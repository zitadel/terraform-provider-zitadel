package active_webkey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2beta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func createActiveWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create active_webkey")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	keyID := d.Get(KeyIDVar).(string)
	_, err = client.ActivateWebKey(helper.CtxWithOrgID(ctx, d), &webkey.ActivateWebKeyRequest{Id: keyID})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get(helper.OrgIDVar).(string), keyID))
	return nil
}

func readActiveWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read active_webkey")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.ListWebKeys(helper.CtxWithOrgID(ctx, d), &webkey.ListWebKeysRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	for _, key := range resp.GetWebKeys() {
		if key.State == webkey.State_STATE_ACTIVE {
			if err := d.Set(KeyIDVar, key.Id); err != nil {
				return diag.FromErr(err)
			}
			d.SetId(fmt.Sprintf("%s:%s", d.Get(helper.OrgIDVar).(string), key.Id))
			return nil
		}
	}
	d.SetId("")
	return nil
}

func updateActiveWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update active_webkey")
	return createActiveWebKey(ctx, d, m)
}

func deleteActiveWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete active_webkey: reverting to default key")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.ListWebKeys(helper.CtxWithOrgID(ctx, d), &webkey.ListWebKeysRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var oldestKey *webkey.WebKey
	for _, key := range resp.GetWebKeys() {
		if oldestKey == nil {
			oldestKey = key
			continue
		}
		if key.GetCreationDate().AsTime().Before(oldestKey.GetCreationDate().AsTime()) {
			oldestKey = key
			continue
		}
		if key.GetCreationDate().AsTime().Equal(oldestKey.GetCreationDate().AsTime()) && key.GetId() < oldestKey.GetId() {
			oldestKey = key
		}
	}

	if oldestKey == nil {
		tflog.Info(ctx, "no keys found, nothing to revert to")
		return nil
	}

	tflog.Info(ctx, fmt.Sprintf("re-activating default key with id %s", oldestKey.GetId()))
	_, err = client.ActivateWebKey(helper.CtxWithOrgID(ctx, d), &webkey.ActivateWebKeyRequest{Id: oldestKey.GetId()})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.FailedPrecondition {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}
