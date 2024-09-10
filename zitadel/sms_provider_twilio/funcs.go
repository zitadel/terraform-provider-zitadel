package sms_provider_twilio

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
		return diag.Errorf("failed to delete sms provider twilio: %v", err)
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

	resp, err := client.AddSMSProviderTwilio(ctx, &admin.AddSMSProviderTwilioRequest{
		Sid:          d.Get(sidVar).(string),
		Token:        d.Get(TokenVar).(string),
		SenderNumber: d.Get(SenderNumberVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create sms provider twilio: %v", err)
	}
	d.SetId(resp.Id)

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

	if d.HasChanges(SenderNumberVar, sidVar) {
		_, err = client.UpdateSMSProviderTwilio(ctx, &admin.UpdateSMSProviderTwilioRequest{
			Id:           d.Id(),
			Sid:          d.Get(sidVar).(string),
			SenderNumber: d.Get(SenderNumberVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update sms provider twilio: %v", err)
		}
	}

	if d.HasChange(TokenVar) {
		_, err = client.UpdateSMSProviderTwilioToken(ctx, &admin.UpdateSMSProviderTwilioTokenRequest{
			Id:    d.Id(),
			Token: d.Get(TokenVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update sms provider twilio: %v", err)
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
		return diag.Errorf("failed to get sms provider twilio")
	}

	set := map[string]interface{}{
		sidVar:          resp.GetConfig().GetTwilio().GetSid(),
		SenderNumberVar: resp.GetConfig().GetTwilio().GetSenderNumber(),
	}
	if token, ok := d.GetOk(TokenVar); ok {
		set[TokenVar] = token
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of sms provider twilio: %v", k, err)
		}
	}
	d.SetId(resp.Config.Id)
	return nil
}
