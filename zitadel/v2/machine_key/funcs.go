package machine_key

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveMachineKey(ctx, &management.RemoveMachineKeyRequest{
		UserId: d.Get(userIDVar).(string),
		KeyId:  d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete machine key: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	keyType := d.Get(keyTypeVar).(string)
	req := &management.AddMachineKeyRequest{
		UserId: d.Get(userIDVar).(string),
		Type:   authn.KeyType(authn.KeyType_value[keyType]),
	}

	if expiration, ok := d.GetOk(expirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, expiration.(string))
		if err != nil {
			return diag.Errorf("failed to parse time: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddMachineKey(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetKeyId())
	if err := d.Set(keyDetailsVar, string(resp.GetKeyDetails())); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(userIDVar).(string)
	resp, err := client.GetMachineKeyByIDs(ctx, &management.GetMachineKeyByIDsRequest{
		UserId: userID,
		KeyId:  d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get machine key")
	}

	d.SetId(resp.GetKey().GetId())
	set := map[string]interface{}{
		expirationDateVar: resp.GetKey().GetExpirationDate().AsTime().Format(time.RFC3339),
		userIDVar:         userID,
		orgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of machine key: %v", k, err)
		}
	}
	return nil
}
