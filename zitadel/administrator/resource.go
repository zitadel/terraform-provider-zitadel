package administrator

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	filter "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"
	permission "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/internal_permission/v2"
)

const (
	UserIDVar       = "user_id"
	ResourceTypeVar = "resource_type"
	ProjectIDVar    = "project_id"
	RolesVar        = "roles"
)

func GetResource() *schema.Resource {
	return newResource(func(ctx context.Context, m interface{}) (permission.InternalPermissionServiceClient, error) {
		info, ok := m.(*helper.ClientInfo)
		if !ok || info == nil {
			return nil, fmt.Errorf("failed to get client")
		}
		return helper.GetInternalPermissionV2Client(ctx, info)
	})
}

func newResource(client helper.BindingClientFactory[permission.InternalPermissionServiceClient]) *schema.Resource {
	r := &schema.Resource{
		Description: "One user's complete administrator role set on one explicit target, using stable InternalPermission v2. Import format: user_id:instance, user_id:organization:org_id, user_id:project:project_id, or user_id:project_grant:project_id:org_id. Percent-encode colons in IDs. Only one resource may manage a given user/target tuple.",
		Schema: map[string]*schema.Schema{
			UserIDVar:       {Type: schema.TypeString, Required: true, ForceNew: true, ValidateDiagFunc: helper.NonEmptyString(UserIDVar), Description: "ID of the human or machine user."},
			ResourceTypeVar: {Type: schema.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringInSlice([]string{"instance", "organization", "project", "project_grant"}, false), Description: "Explicit administrator scope: instance, organization, project or project_grant. There is no default scope."},
			helper.OrgIDVar: {Type: schema.TypeString, Optional: true, ForceNew: true, ValidateDiagFunc: helper.NonEmptyString(helper.OrgIDVar), Description: "Organization target, or grantee organization for project_grant. Forbidden for instance/project scope."},
			ProjectIDVar:    {Type: schema.TypeString, Optional: true, ForceNew: true, ValidateDiagFunc: helper.NonEmptyString(ProjectIDVar), Description: "Project target or project owning the grant. Forbidden for instance/organization scope."},
			RolesVar:        {Type: schema.TypeSet, Required: true, MinItems: 1, Elem: &schema.Schema{Type: schema.TypeString, ValidateDiagFunc: helper.NonEmptyString("role")}, Description: "Complete administrator role set. Roles are validated by Zitadel for the selected scope; updates revoke omitted roles."},
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			// References can be unknown during planning. Validate once the target
			// is known; create also checks before any request reaches the API.
			for _, key := range []string{UserIDVar, ResourceTypeVar, helper.OrgIDVar, ProjectIDVar} {
				if !d.NewValueKnown(key) {
					return nil
				}
			}
			return targetFrom(d).validate()
		},
		Importer: &schema.ResourceImporter{StateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
			t, err := parseIdentity(d.Id())
			if err != nil {
				return nil, err
			}
			if err := t.set(d); err != nil {
				return nil, err
			}
			d.SetId(t.identity())
			return []*schema.ResourceData{d}, nil
		}},
	}
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		t := targetFrom(d)
		if err := t.validate(); err != nil {
			return diag.FromErr(err)
		}
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.CreateAdministrator(ctx, &permission.CreateAdministratorRequest{UserId: t.user, Resource: t.resource(), Roles: helper.GetOkSetToStringSlice(d, RolesVar)})
		// Do not adopt AlreadyExists implicitly: another declaration may own it.
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(t.identity())
		return nil
	}
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		t, err := parseIdentity(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.UpdateAdministrator(ctx, &permission.UpdateAdministratorRequest{UserId: t.user, Resource: t.resource(), Roles: helper.GetOkSetToStringSlice(d, RolesVar)})
		return diag.FromErr(err)
	}
	r.DeleteContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		t, err := parseIdentity(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DeleteAdministrator(ctx, &permission.DeleteAdministratorRequest{UserId: t.user, Resource: t.resource()})
		if err = helper.IgnoreIfNotFoundError(err); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("")
		return nil
	}
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		t, err := parseIdentity(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		c, err := client(ctx, m)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := c.ListAdministrators(ctx, &permission.ListAdministratorsRequest{
			Pagination: &filter.PaginationRequest{Limit: 2},
			Filters: []*permission.AdministratorSearchFilter{
				{Filter: &permission.AdministratorSearchFilter_InUserIdsFilter{InUserIdsFilter: &filter.InIDsFilter{Ids: []string{t.user}}}},
				{Filter: &permission.AdministratorSearchFilter_Resource{Resource: t.filter()}},
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(resp.GetAdministrators()) == 0 {
			d.SetId("")
			return nil
		}
		if len(resp.GetAdministrators()) != 1 || !t.matches(resp.GetAdministrators()[0]) {
			return diag.Errorf("administrator query returned a different or ambiguous identity")
		}
		if err := t.set(d); err != nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(d.Set(RolesVar, resp.GetAdministrators()[0].GetRoles()))
	}
	return r
}

// Administrator has no server-generated ID. The durable identity is the user
// plus the full target tuple. Operations use that identity, including after an
// import, rather than accidentally applying a new spec target to an old ID.
type target struct{ user, scope, project, org string }

func targetFrom(d interface{ Get(string) interface{} }) target {
	return target{user: d.Get(UserIDVar).(string), scope: d.Get(ResourceTypeVar).(string), project: d.Get(ProjectIDVar).(string), org: d.Get(helper.OrgIDVar).(string)}
}

func (t target) validate() error {
	if strings.TrimSpace(t.user) == "" {
		return fmt.Errorf("user_id must not be empty")
	}
	valid := false
	switch t.scope {
	case "instance":
		valid = t.org == "" && t.project == ""
	case "organization":
		valid = strings.TrimSpace(t.org) != "" && t.project == ""
	case "project":
		valid = strings.TrimSpace(t.project) != "" && t.org == ""
	case "project_grant":
		valid = strings.TrimSpace(t.project) != "" && strings.TrimSpace(t.org) != ""
	}
	if !valid {
		return fmt.Errorf("invalid administrator target: instance requires no IDs; organization requires only org_id; project requires only project_id; project_grant requires project_id and grantee org_id")
	}
	return nil
}

func (t target) identity() string {
	parts := []string{t.user, t.scope}
	if t.project != "" {
		parts = append(parts, t.project)
	}
	if t.org != "" {
		parts = append(parts, t.org)
	}
	for i := range parts {
		parts[i] = url.QueryEscape(parts[i])
	}
	return strings.Join(parts, ":")
}

func parseIdentity(id string) (target, error) {
	parts := strings.Split(id, ":")
	if len(parts) < 2 || len(parts) > 4 {
		return target{}, fmt.Errorf("invalid administrator import ID: expected user_id:resource_type[:project_id][:org_id]")
	}
	for i := range parts {
		v, err := url.QueryUnescape(parts[i])
		if err != nil {
			return target{}, fmt.Errorf("invalid administrator import ID: %w", err)
		}
		parts[i] = v
	}
	t := target{user: parts[0], scope: parts[1]}
	switch {
	case t.scope == "instance" && len(parts) == 2:
	case t.scope == "organization" && len(parts) == 3:
		t.org = parts[2]
	case t.scope == "project" && len(parts) == 3:
		t.project = parts[2]
	case t.scope == "project_grant" && len(parts) == 4:
		t.project, t.org = parts[2], parts[3]
	default:
		return target{}, fmt.Errorf("invalid administrator import ID: target scope and number of IDs do not match")
	}
	return t, t.validate()
}

func (t target) set(d *schema.ResourceData) error {
	for k, v := range map[string]string{UserIDVar: t.user, ResourceTypeVar: t.scope, ProjectIDVar: t.project, helper.OrgIDVar: t.org} {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (t target) resource() *permission.ResourceType {
	r := &permission.ResourceType{}
	switch t.scope {
	case "instance":
		r.Resource = &permission.ResourceType_Instance{Instance: true}
	case "organization":
		r.Resource = &permission.ResourceType_OrganizationId{OrganizationId: t.org}
	case "project":
		r.Resource = &permission.ResourceType_ProjectId{ProjectId: t.project}
	case "project_grant":
		r.Resource = &permission.ResourceType_ProjectGrant_{ProjectGrant: &permission.ResourceType_ProjectGrant{ProjectId: t.project, OrganizationId: t.org}}
	}
	return r
}

func (t target) filter() *permission.ResourceFilter {
	r := &permission.ResourceFilter{}
	switch t.scope {
	case "instance":
		r.Resource = &permission.ResourceFilter_Instance{Instance: true}
	case "organization":
		r.Resource = &permission.ResourceFilter_OrganizationId{OrganizationId: t.org}
	case "project":
		r.Resource = &permission.ResourceFilter_ProjectId{ProjectId: t.project}
	case "project_grant":
		r.Resource = &permission.ResourceFilter_ProjectGrant_{ProjectGrant: &permission.ResourceFilter_ProjectGrant{ProjectId: t.project, OrganizationId: t.org}}
	}
	return r
}

func (t target) matches(a *permission.Administrator) bool {
	if a.GetUser().GetId() != t.user {
		return false
	}
	switch t.scope {
	case "instance":
		return a.GetInstance()
	case "organization":
		return a.GetOrganization().GetId() == t.org
	case "project":
		return a.GetProject().GetId() == t.project
	case "project_grant":
		return a.GetProjectGrant().GetProjectId() == t.project && a.GetProjectGrant().GetOrganizationId() == t.org
	}
	return false
}
