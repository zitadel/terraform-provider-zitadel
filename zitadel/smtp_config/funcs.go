package smtp_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
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

	_, err = client.RemoveSMTPConfig(ctx, &admin.RemoveSMTPConfigRequest{
		Id: d.Id(),
	})
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

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &admin.AddSMTPConfigRequest{
		SenderAddress:  d.Get(SenderAddressVar).(string),
		SenderName:     d.Get(SenderNameVar).(string),
		Host:           d.Get(hostVar).(string),
		User:           d.Get(userVar).(string),
		Tls:            d.Get(tlsVar).(bool),
		Password:       d.Get(PasswordVar).(string),
		ReplyToAddress: d.Get(replyToAddressVar).(string),
	}

	resp, err := client.AddSMTPConfig(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create smtp config: %v", err)
	}
	d.SetId(resp.GetId())

	if d.Get(SetActiveVar).(bool) {
		if _, err := client.ActivateSMTPConfig(ctx, &admin.ActivateSMTPConfigRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate smtp config: %v", err)
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

	if d.HasChanges(SenderAddressVar, SenderNameVar, tlsVar, hostVar, userVar, replyToAddressVar) {
		_, err = client.UpdateSMTPConfig(ctx, &admin.UpdateSMTPConfigRequest{
			Id:             d.Id(),
			SenderAddress:  d.Get(SenderAddressVar).(string),
			SenderName:     d.Get(SenderNameVar).(string),
			Host:           d.Get(hostVar).(string),
			Tls:            d.Get(tlsVar).(bool),
			User:           d.Get(userVar).(string),
			ReplyToAddress: d.Get(replyToAddressVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config: %v", err)
		}
	}

	if d.HasChange(PasswordVar) {
		_, err = client.UpdateSMTPConfigPassword(ctx, &admin.UpdateSMTPConfigPasswordRequest{
			Id:       d.Id(),
			Password: d.Get(PasswordVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config password: %v", err)
		}
	}

	if d.HasChange(PasswordVar) && d.Get(SetActiveVar).(bool) {
		if _, err := client.ActivateSMTPConfig(ctx, &admin.ActivateSMTPConfigRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate smtp config: %v", err)
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

	resp, err := client.GetSMTPConfigById(ctx, &admin.GetSMTPConfigByIdRequest{
		Id: d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get smtp config")
	}

	set := map[string]interface{}{
		SenderAddressVar:  resp.GetSmtpConfig().GetSenderAddress(),
		SenderNameVar:     resp.GetSmtpConfig().GetSenderName(),
		tlsVar:            resp.GetSmtpConfig().GetTls(),
		hostVar:           resp.GetSmtpConfig().GetHost(),
		userVar:           resp.GetSmtpConfig().GetUser(),
		PasswordVar:       d.Get(PasswordVar).(string),
		replyToAddressVar: resp.GetSmtpConfig().GetReplyToAddress(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of smtp config: %v", k, err)
		}
	}
	d.SetId(resp.GetSmtpConfig().GetId())
	return nil
}
