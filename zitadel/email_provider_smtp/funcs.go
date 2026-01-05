package email_provider_smtp

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

	_, err = client.RemoveEmailProvider(ctx, &admin.RemoveEmailProviderRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to delete email smtp provider: %v", err)
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

	resp, err := client.AddEmailProviderSMTP(ctx, &admin.AddEmailProviderSMTPRequest{
		SenderAddress:  d.Get(SenderAddressVar).(string),
		SenderName:     d.Get(SenderNameVar).(string),
		Host:           d.Get(hostVar).(string),
		User:           d.Get(userVar).(string),
		Tls:            d.Get(tlsVar).(bool),
		Password:       d.Get(PasswordVar).(string),
		ReplyToAddress: d.Get(replyToAddressVar).(string),
		Description:    d.Get(DescriptionVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create email provider smtp: %v", err)
	}
	d.SetId(resp.Id)

	if d.Get(setActiveVar).(bool) {
		if _, err := client.ActivateEmailProvider(ctx, &admin.ActivateEmailProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate email smtp provider config: %v", err)
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

	if d.HasChanges(SenderAddressVar, SenderNameVar, tlsVar, hostVar, userVar, replyToAddressVar, PasswordVar, DescriptionVar) {
		_, err = client.UpdateEmailProviderSMTP(ctx, &admin.UpdateEmailProviderSMTPRequest{
			Id:             d.Id(),
			SenderAddress:  d.Get(SenderAddressVar).(string),
			SenderName:     d.Get(SenderNameVar).(string),
			Host:           d.Get(hostVar).(string),
			Tls:            d.Get(tlsVar).(bool),
			User:           d.Get(userVar).(string),
			ReplyToAddress: d.Get(replyToAddressVar).(string),
			Password:       d.Get(PasswordVar).(string),
			Description:    d.Get(DescriptionVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update email provider smtp: %v", err)
		}
	}

	if d.HasChange(setActiveVar) && d.Get(setActiveVar).(bool) {
		if _, err = client.ActivateEmailProvider(ctx, &admin.ActivateEmailProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate email provider smtp: %v", err)
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

	resp, err := client.GetEmailProviderById(ctx, &admin.GetEmailProviderByIdRequest{
		Id: d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get email smtp provider")
	}

	if resp.GetConfig().GetSmtp() == nil {
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		SenderAddressVar:  resp.GetConfig().GetSmtp().GetSenderAddress(),
		SenderNameVar:     resp.GetConfig().GetSmtp().GetSenderName(),
		tlsVar:            resp.GetConfig().GetSmtp().GetTls(),
		hostVar:           resp.GetConfig().GetSmtp().GetHost(),
		userVar:           resp.GetConfig().GetSmtp().GetUser(),
		PasswordVar:       d.Get(PasswordVar).(string),
		replyToAddressVar: resp.GetConfig().GetSmtp().GetReplyToAddress(),
		DescriptionVar:    resp.GetConfig().GetDescription(),
		setActiveVar:      d.Get(setActiveVar).(bool),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of email provider smtp: %v", k, err)
		}
	}
	d.SetId(resp.Config.Id)
	return nil
}
