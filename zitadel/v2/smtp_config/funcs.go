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

	tls, tlsOk := d.GetOk(tlsVar)
	user, userOk := d.GetOk(userVar)
	password, pwOk := d.GetOk(passwordVar)
	req := &admin.AddSMTPConfigRequest{
		SenderAddress: d.Get(senderAddressVar).(string),
		SenderName:    d.Get(senderNameVar).(string),
		Host:          d.Get(hostVar).(string),
	}
	if tlsOk {
		req.Tls = tls.(bool)
	}
	if userOk {
		req.User = user.(string)
	}
	if pwOk {
		req.Password = password.(string)
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
		smtp, err := client.GetSMTPConfig(ctx, &admin.GetSMTPConfigRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		senderAddress := d.Get(senderAddressVar).(string)
		senderName := d.Get(senderNameVar).(string)
		tls, tlsOk := d.GetOk(tlsVar)
		host := d.Get(hostVar).(string)
		user, userOk := d.GetOk(userVar)

		if smtp.SmtpConfig.SenderName != senderName ||
			smtp.SmtpConfig.SenderAddress != senderAddress ||
			smtp.SmtpConfig.Tls != tls ||
			smtp.SmtpConfig.Host != host ||
			smtp.SmtpConfig.User != user {

			req := &admin.UpdateSMTPConfigRequest{
				SenderAddress: senderAddress,
				SenderName:    senderName,
				Host:          host,
			}
			if tlsOk {
				req.Tls = tls.(bool)
			}
			if userOk {
				req.User = user.(string)
			}

			_, err = client.UpdateSMTPConfig(ctx, req)
			if err != nil {
				return diag.Errorf("failed to update smtp config: %v", err)
			}
		}
	}

	if d.HasChange(passwordVar) {
		password, pwOk := d.GetOk(passwordVar)
		req := &admin.UpdateSMTPConfigPasswordRequest{}
		if pwOk {
			req.Password = password.(string)
		}
		_, err = client.UpdateSMTPConfigPassword(ctx, req)
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
	if err != nil {
		d.SetId("")
		return nil
	}
	password, pwOk := d.GetOk(passwordVar)
	set := map[string]interface{}{
		senderAddressVar: resp.GetSmtpConfig().GetSenderAddress(),
		senderNameVar:    resp.GetSmtpConfig().GetSenderName(),
		tlsVar:           resp.GetSmtpConfig().GetTls(),
		hostVar:          resp.GetSmtpConfig().GetHost(),
		userVar:          resp.GetSmtpConfig().GetUser(),
	}
	if pwOk {
		set[passwordVar] = password.(string)
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of smtp config: %v", k, err)
		}
	}
	d.SetId(resp.SmtpConfig.Details.ResourceOwner)
	return nil
}
