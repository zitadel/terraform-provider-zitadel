package default_oidc_settings

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "default oidc settings cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	accessTokenLT, err := time.ParseDuration(d.Get(accessTokenLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	idTokenLT, err := time.ParseDuration(d.Get(idTokenLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	refreshTokenExp, err := time.ParseDuration(d.Get(RefreshTokenExpirationVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	refreshTokenIdleExp, err := time.ParseDuration(d.Get(refreshTokenIdleExpirationVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.UpdateOIDCSettings(ctx, &admin.UpdateOIDCSettingsRequest{
		AccessTokenLifetime:        durationpb.New(accessTokenLT),
		IdTokenLifetime:            durationpb.New(idTokenLT),
		RefreshTokenIdleExpiration: durationpb.New(refreshTokenIdleExp),
		RefreshTokenExpiration:     durationpb.New(refreshTokenExp),
	})
	id := resp.GetDetails().GetResourceOwner()
	if err != nil {
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default oidc settings: %v", err)
		}
	}
	if id == "" {
		getResp, getErr := client.GetOIDCSettings(ctx, &admin.GetOIDCSettingsRequest{})
		if getErr != nil {
			return diag.Errorf("failed to get new default oidc settings id: %v", getErr)
		}
		id = getResp.GetSettings().GetDetails().GetResourceOwner()
	}
	d.SetId(id)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetOIDCSettings(ctx, &admin.GetOIDCSettingsRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default oidc settings: %v", err)
	}

	set := map[string]interface{}{
		accessTokenLifetimeVar:        resp.GetSettings().GetAccessTokenLifetime().AsDuration().String(),
		idTokenLifetimeVar:            resp.GetSettings().GetIdTokenLifetime().AsDuration().String(),
		refreshTokenIdleExpirationVar: resp.GetSettings().GetRefreshTokenIdleExpiration().AsDuration().String(),
		RefreshTokenExpirationVar:     resp.GetSettings().GetRefreshTokenExpiration().AsDuration().String(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default oidc settings: %v", k, err)
		}
	}
	d.SetId(resp.GetSettings().GetDetails().GetResourceOwner())
	return nil
}
