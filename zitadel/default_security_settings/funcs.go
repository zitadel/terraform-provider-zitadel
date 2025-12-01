package default_security_settings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default security settings cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSecuritySettingsClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(EnableImpersonationVar, embeddedIframeEnabledVar, embeddedIframeAllowedOriginsVar) {
		req := &settingsv2.SetSecuritySettingsRequest{
			EnableImpersonation: d.Get(EnableImpersonationVar).(bool),
		}

		if d.Get(embeddedIframeEnabledVar) != nil || d.Get(embeddedIframeAllowedOriginsVar) != nil {
			req.EmbeddedIframe = &settingsv2.EmbeddedIframeSettings{
				Enabled:        d.Get(embeddedIframeEnabledVar).(bool),
				AllowedOrigins: helper.GetOkSetToStringSlice(d, embeddedIframeAllowedOriginsVar),
			}
		}

		_, err := client.SetSecuritySettings(ctx, req)
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default security settings: %v", err)
		}
	}

	d.SetId("default_security_settings")
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSecuritySettingsClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetSecuritySettings(ctx, &settingsv2.GetSecuritySettingsRequest{})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get default security settings: %v", err)
	}

	settings := resp.Settings
	set := map[string]interface{}{
		EnableImpersonationVar: settings.GetEnableImpersonation(),
	}

	if settings.GetEmbeddedIframe() != nil {
		set[embeddedIframeEnabledVar] = settings.GetEmbeddedIframe().GetEnabled()
		set[embeddedIframeAllowedOriginsVar] = settings.GetEmbeddedIframe().GetAllowedOrigins()
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default security settings: %v", k, err)
		}
	}

	d.SetId("default_security_settings")
	return nil
}
