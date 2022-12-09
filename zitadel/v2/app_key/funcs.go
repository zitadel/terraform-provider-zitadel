package app_key

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

	_, err = client.RemoveAppKey(ctx, &management.RemoveAppKeyRequest{
		ProjectId: d.Get(projectIDVar).(string),
		AppId:     d.Get(appIDVar).(string),
		KeyId:     d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete app key: %v", err)
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
	req := &management.AddAppKeyRequest{
		ProjectId: d.Get(projectIDVar).(string),
		AppId:     d.Get(appIDVar).(string),
		Type:      authn.KeyType(authn.KeyType_value[keyType]),
	}

	if expiration, ok := d.GetOk(expirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, expiration.(string))
		if err != nil {
			return diag.Errorf("failed to parse time: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddAppKey(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetId())
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

	projectID := d.Get(projectIDVar).(string)
	appID := d.Get(appIDVar).(string)
	resp, err := client.GetAppKey(ctx, &management.GetAppKeyRequest{
		ProjectId: projectID,
		AppId:     appID,
		KeyId:     d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get app key")
	}
	d.SetId(resp.GetKey().GetId())

	set := map[string]interface{}{
		expirationDateVar: resp.GetKey().GetExpirationDate().AsTime().Format(time.RFC3339),
		projectIDVar:      projectID,
		appIDVar:          appID,
		orgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of app key: %v", k, err)
		}
	}
	return nil
}
