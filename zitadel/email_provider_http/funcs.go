package email_provider_http

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
		return diag.Errorf("failed to delete email http provider: %v", err)
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

	resp, err := client.AddEmailProviderHTTP(ctx, &admin.AddEmailProviderHTTPRequest{
		Endpoint:    d.Get(EndpointVar).(string),
		Description: d.Get(DescriptionVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create email provider http: %v", err)
	}
	d.SetId(resp.Id)

	if err := d.Set(SigningKeyVar, resp.SigningKey); err != nil {
		return diag.Errorf("failed to set signing_key: %v", err)
	}

	if d.Get(setActiveVar).(bool) {
		if _, err := client.ActivateEmailProvider(ctx, &admin.ActivateEmailProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate email http provider config: %v", err)
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

	if d.HasChanges(EndpointVar, DescriptionVar) {
		_, err = client.UpdateEmailProviderHTTP(ctx, &admin.UpdateEmailProviderHTTPRequest{
			Id:          d.Id(),
			Endpoint:    d.Get(EndpointVar).(string),
			Description: d.Get(DescriptionVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update email provider http: %v", err)
		}
	}

	if d.HasChange(setActiveVar) && d.Get(setActiveVar).(bool) {
		if _, err = client.ActivateEmailProvider(ctx, &admin.ActivateEmailProviderRequest{Id: d.Id()}); err != nil {
			return diag.Errorf("failed to activate email provider http: %v", err)
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

	resp, err := client.GetEmailProvider(ctx, &admin.GetEmailProviderRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get email http provider")
	}

	set := map[string]interface{}{
		EndpointVar:    resp.GetConfig().GetHttp().GetEndpoint(),
		DescriptionVar: resp.GetConfig().GetDescription(),
		SigningKeyVar:  d.Get(SigningKeyVar).(string),
		setActiveVar:   d.Get(setActiveVar).(bool),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of email provider: %v", k, err)
		}
	}
	d.SetId(resp.Config.Id)
	return nil
}
