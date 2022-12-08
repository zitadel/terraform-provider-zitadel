package smtp_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveSMTPConfig(ctx, &admin.RemoveSMTPConfigRequest{})
	if err != nil {
		return diag.Errorf("failed to delete smtp config: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &admin.AddSMTPConfigRequest{
		SenderAddress: d.Get(senderAddressVar).(string),
		SenderName:    d.Get(senderNameVar).(string),
		Host:          d.Get(hostVar).(string),
		User:          d.Get(userVar).(string),
		Tls:           d.Get(tlsVar).(bool),
		Password:      d.Get(passwordVar).(string),
	}

	resp, err := client.AddSMTPConfig(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create smtp config: %v", err)
	}
	d.SetId(resp.Details.ResourceOwner)

	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(senderAddressVar, senderNameVar, tlsVar, hostVar, userVar) {
		_, err = client.UpdateSMTPConfig(ctx, &admin.UpdateSMTPConfigRequest{
			SenderAddress: d.Get(senderAddressVar).(string),
			SenderName:    d.Get(senderNameVar).(string),
			Host:          d.Get(hostVar).(string),
			Tls:           d.Get(tlsVar).(bool),
			User:          d.Get(userVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config: %v", err)
		}
	}

	if d.HasChange(passwordVar) {
		_, err = client.UpdateSMTPConfigPassword(ctx, &admin.UpdateSMTPConfigPasswordRequest{
			Password: d.Get(passwordVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config password: %v", err)
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

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetSMTPConfig(ctx, &admin.GetSMTPConfigRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get smtp config")
	}

	set := map[string]interface{}{
		senderAddressVar: resp.GetSmtpConfig().GetSenderAddress(),
		senderNameVar:    resp.GetSmtpConfig().GetSenderName(),
		tlsVar:           resp.GetSmtpConfig().GetTls(),
		hostVar:          resp.GetSmtpConfig().GetHost(),
		userVar:          resp.GetSmtpConfig().GetUser(),
		passwordVar:      d.Get(passwordVar).(string),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of smtp config: %v", k, err)
		}
	}
	d.SetId(resp.SmtpConfig.Details.ResourceOwner)
	return nil
}
