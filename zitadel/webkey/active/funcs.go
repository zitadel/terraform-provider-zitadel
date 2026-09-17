package active_webkey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

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

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

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

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "the active web key cannot be deactivated, leaving it unchanged")
	return nil
}
