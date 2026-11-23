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

	if d.IsNewResource() || d.HasChanges(
		enableImpersonationVar,
		embeddedIframeEnabledVar,
		embeddedIframeAllowedOriginsVar,
		dynamicClientRegistrationEnabledVar,
		dynamicClientRegistrationAllowUnauthenticatedVar,
	) {
		dcrEnabled := configuredBool(d, dynamicClientRegistrationEnabledVar)
		dcrAllowUnauthenticated := configuredBool(d, dynamicClientRegistrationAllowUnauthenticatedVar)

		// SetSecuritySettings replaces the whole setting, so a field left out of the
		// request is reset. Keep the current dynamic client registration values for the
		// attributes the configuration does not set.
		var current *settingsv2.DynamicClientRegistrationSettings
		if dcrEnabled == nil || dcrAllowUnauthenticated == nil {
			resp, err := client.GetSecuritySettings(ctx, &settingsv2.GetSecuritySettingsRequest{})
			if err != nil {
				return diag.Errorf("failed to get default security settings: %v", err)
			}
			current = resp.GetSettings().GetDynamicClientRegistration()
		}

		req := &settingsv2.SetSecuritySettingsRequest{
			EnableImpersonation: d.Get(enableImpersonationVar).(bool),
			EmbeddedIframe: &settingsv2.EmbeddedIframeSettings{
				Enabled:        d.Get(embeddedIframeEnabledVar).(bool),
				AllowedOrigins: helper.GetOkSetToStringSlice(d, embeddedIframeAllowedOriginsVar),
			},
			DynamicClientRegistration: dynamicClientRegistrationSettings(dcrEnabled, dcrAllowUnauthenticated, current),
		}

		_, err := client.SetSecuritySettings(ctx, req)
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default security settings: %v", err)
		}

		// Store the values sent, including those kept from ZITADEL, so computed
		// attributes are known after create.
		set := map[string]interface{}{
			dynamicClientRegistrationEnabledVar:              req.GetDynamicClientRegistration().GetEnabled(),
			dynamicClientRegistrationAllowUnauthenticatedVar: req.GetDynamicClientRegistration().GetAllowUnauthenticated(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of default security settings: %v", k, err)
			}
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
		enableImpersonationVar: settings.GetEnableImpersonation(),
	}

	if settings.GetEmbeddedIframe() != nil {
		set[embeddedIframeEnabledVar] = settings.GetEmbeddedIframe().GetEnabled()
		set[embeddedIframeAllowedOriginsVar] = settings.GetEmbeddedIframe().GetAllowedOrigins()
	}
	set[dynamicClientRegistrationEnabledVar] = settings.GetDynamicClientRegistration().GetEnabled()
	set[dynamicClientRegistrationAllowUnauthenticatedVar] = settings.GetDynamicClientRegistration().GetAllowUnauthenticated()

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default security settings: %v", k, err)
		}
	}

	d.SetId("default_security_settings")
	return nil
}

// configuredBool returns the value of a bool attribute in the practitioner's
// configuration, or nil if the configuration does not set it. Unlike d.Get it
// does not fall back to state, which matters for Optional+Computed attributes.
func configuredBool(d *schema.ResourceData, attr string) *bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	v := raw.GetAttr(attr)
	if v.IsNull() || !v.IsKnown() {
		return nil
	}
	b := v.True()
	return &b
}

// dynamicClientRegistrationSettings returns the dynamic client registration part of
// a SetSecuritySettings request: the configured values, and for those not
// configured (nil) the current values from ZITADEL.
func dynamicClientRegistrationSettings(enabled, allowUnauthenticated *bool, current *settingsv2.DynamicClientRegistrationSettings) *settingsv2.DynamicClientRegistrationSettings {
	settings := &settingsv2.DynamicClientRegistrationSettings{
		Enabled:              current.GetEnabled(),
		AllowUnauthenticated: current.GetAllowUnauthenticated(),
	}
	if enabled != nil {
		settings.Enabled = *enabled
	}
	if allowUnauthenticated != nil {
		settings.AllowUnauthenticated = *allowUnauthenticated
	}
	return settings
}
