package authorization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	authz "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authorization/v2"
	filter "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"
)

const (
	UserIDVar    = "user_id"
	ProjectIDVar = "project_id"
	RoleKeysVar  = "role_keys"
)

func GetResource() *schema.Resource {
	return newResource(func(ctx context.Context, m interface{}) (authz.AuthorizationServiceClient, error) {
		info, ok := m.(*helper.ClientInfo)
		if !ok || info == nil {
			return nil, fmt.Errorf("failed to get client")
		}
		return helper.GetAuthorizationV2Client(ctx, info)
	})
}

func newResource(client helper.BindingClientFactory[authz.AuthorizationServiceClient]) *schema.Resource {
	r := &schema.Resource{
		Description: "One user's complete application role assignment, using the stable Authorization v2 API. Import an existing authorization by ID. Only one resource may manage a given user/project/organization assignment.",
		Schema: map[string]*schema.Schema{
			UserIDVar:       requiredID("ID of the human or machine user."),
			ProjectIDVar:    requiredID("ID of the project."),
			helper.OrgIDVar: requiredID("ID of the organization the authorization belongs to. It must own the project or already have a project grant; it need not be the user's organization."),
			RoleKeysVar:     {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, ValidateDiagFunc: helper.NonEmptyString("role key")}, Description: "Complete set of role keys. Updating this set revokes roles omitted from it."},
			"state":         {Type: schema.TypeString, Computed: true, Description: "Observed authorization state. This resource does not activate or deactivate authorizations."},
		},
		Importer: &schema.ResourceImporter{StateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
			if err := validateID(d.Id()); err != nil {
				return nil, err
			}
			return []*schema.ResourceData{d}, nil
		}},
	}
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := c.CreateAuthorization(ctx, &authz.CreateAuthorizationRequest{
			UserId: d.Get(UserIDVar).(string), ProjectId: d.Get(ProjectIDVar).(string),
			OrganizationId: d.Get(helper.OrgIDVar).(string), RoleKeys: helper.GetOkSetToStringSlice(d, RoleKeysVar),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if resp.GetId() == "" {
			return diag.Errorf("create authorization returned no ID")
		}
		d.SetId(resp.GetId())
		return nil
	}
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.UpdateAuthorization(ctx, &authz.UpdateAuthorizationRequest{Id: d.Id(), RoleKeys: helper.GetOkSetToStringSlice(d, RoleKeysVar)})
		return diag.FromErr(err)
	}
	r.DeleteContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DeleteAuthorization(ctx, &authz.DeleteAuthorizationRequest{Id: d.Id()})
		if err = helper.IgnoreIfNotFoundError(err); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("")
		return nil
	}
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		// There is no GetAuthorization RPC. The ID filter uniquely addresses one
		// assignment, so this query cannot silently truncate an unfiltered list.
		resp, err := c.ListAuthorizations(ctx, &authz.ListAuthorizationsRequest{
			Pagination: &filter.PaginationRequest{Limit: 2},
			Filters:    []*authz.AuthorizationsSearchFilter{{Filter: &authz.AuthorizationsSearchFilter_AuthorizationIds{AuthorizationIds: &filter.InIDsFilter{Ids: []string{d.Id()}}}}},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(resp.GetAuthorizations()) == 0 {
			d.SetId("")
			return nil
		}
		if len(resp.GetAuthorizations()) != 1 || resp.GetAuthorizations()[0].GetId() != d.Id() {
			return diag.Errorf("authorization query returned a different or ambiguous identity")
		}
		a := resp.GetAuthorizations()[0]
		roles := make([]string, 0, len(a.GetRoles()))
		for _, role := range a.GetRoles() {
			roles = append(roles, role.GetKey())
		}
		for key, value := range map[string]interface{}{
			UserIDVar: a.GetUser().GetId(), ProjectIDVar: a.GetProject().GetId(), helper.OrgIDVar: a.GetOrganization().GetId(),
			RoleKeysVar: roles, "state": a.GetState().String(),
		} {
			if err := d.Set(key, value); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}
	return r
}

func requiredID(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true, ValidateDiagFunc: helper.NonEmptyString("ID"), Description: description}
}

func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("authorization ID must not be empty")
	}
	return nil
}
