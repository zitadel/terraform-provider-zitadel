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

	resp, err := client.AddSMTPConfig(ctx, &admin.AddSMTPConfigRequest{
		SenderAddress: d.Get(senderAddressVar).(string),
		SenderName:    d.Get(senderNameVar).(string),
		Tls:           d.Get(tlsVar).(bool),
		Host:          d.Get(hostVar).(string),
		User:          d.Get(userVar).(string),
		Password:      d.Get(passwordVar).(string),
	})
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

	smtp, err := client.GetSMTPConfig(ctx, &admin.GetSMTPConfigRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	senderAddress := d.Get(senderAddressVar).(string)
	senderName := d.Get(senderNameVar).(string)
	tls := d.Get(tlsVar).(bool)
	host := d.Get(hostVar).(string)
	user := d.Get(userVar).(string)
	if smtp.SmtpConfig.SenderName != senderName ||
		smtp.SmtpConfig.SenderAddress != senderAddress ||
		smtp.SmtpConfig.Tls != tls ||
		smtp.SmtpConfig.Host != host ||
		smtp.SmtpConfig.User != user {

		_, err = client.UpdateSMTPConfig(ctx, &admin.UpdateSMTPConfigRequest{
			SenderAddress: senderAddress,
			SenderName:    senderName,
			Tls:           tls,
			Host:          host,
			User:          user,
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config: %v", err)
		}
	} else {
		_, err = client.UpdateSMTPConfigPassword(ctx, &admin.UpdateSMTPConfigPasswordRequest{
			Password: d.Get(passwordVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update smtp config: %v", err)
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
		//return diag.Errorf("error while reading smtp config: %v", err)
	}
	d.SetId(resp.SmtpConfig.Details.ResourceOwner)
	return nil
}
