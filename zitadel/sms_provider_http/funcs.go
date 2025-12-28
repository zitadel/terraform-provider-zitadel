package sms_provider_http

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveSMSProvider(ctx, &admin.RemoveSMSProviderRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to delete sms http provider: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddSMSProviderHTTP(ctx, &admin.AddSMSProviderHTTPRequest{
		Endpoint:    d.Get(EndPointVar).(string),
		Description: d.Get(DescriptionVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create sms provider http: %v", err)
	}
	d.SetId(resp.Id)

	if d.Get(setActiveVar).(bool) {
		if _, err := client.ActivateSMSProvider(ctx, &admin.ActivateSMSProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate sms http provider config: %v", err)
		}
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(EndPointVar, DescriptionVar) {
		_, err = client.UpdateSMSProviderHTTP(ctx, &admin.UpdateSMSProviderHTTPRequest{
			Id:          d.Id(),
			Endpoint:    d.Get(EndPointVar).(string),
			Description: d.Get(DescriptionVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update sms provider http: %v", err)
		}
	}

	if d.HasChange(setActiveVar) && d.Get(setActiveVar).(bool) {
		if _, err = client.ActivateSMSProvider(ctx, &admin.ActivateSMSProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate sms provider http: %v", err)
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetSMSProvider(ctx, &admin.GetSMSProviderRequest{
		Id: d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get sms http provider: %v", err)
	}

	set := map[string]interface{}{
		EndPointVar:    resp.GetConfig().GetHttp().GetEndpoint(),
		DescriptionVar: resp.GetConfig().GetDescription(),
		setActiveVar:   d.Get(setActiveVar).(bool),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of sms provider: %v", k, err)
		}
	}
	d.SetId(resp.Config.Id)
	return nil
}
