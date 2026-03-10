package instance_restrictions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "instance restrictions cannot be deleted")
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	d.SetId("instance_restrictions")

	return update(ctx, d, m)
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

	if d.HasChanges(disallowPublicOrgRegistrationVar, allowedLanguagesVar) {
		req := &admin.SetRestrictionsRequest{}

		if d.HasChange(disallowPublicOrgRegistrationVar) {
			v := d.Get(disallowPublicOrgRegistrationVar).(bool)
			req.DisallowPublicOrgRegistration = &v
		}

		if d.HasChange(allowedLanguagesVar) {
			req.AllowedLanguages = &admin.SelectLanguages{
				List: helper.GetOkSetToStringSlice(d, allowedLanguagesVar),
			}
		}

		_, err := client.SetRestrictions(ctx, req)
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update instance restrictions: %v", err)
		}
	}

	d.SetId("instance_restrictions")
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

	resp, err := client.GetRestrictions(ctx, &admin.GetRestrictionsRequest{})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get instance restrictions: %v", err)
	}

	set := map[string]interface{}{
		disallowPublicOrgRegistrationVar: resp.GetDisallowPublicOrgRegistration(),
		allowedLanguagesVar:              resp.GetAllowedLanguages(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of instance restrictions: %v", k, err)
		}
	}

	d.SetId("instance_restrictions")
	return nil
}
