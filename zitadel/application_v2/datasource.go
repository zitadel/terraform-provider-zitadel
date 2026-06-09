package application_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// GetDatasource returns the singular datasource for one application looked up by app_id.
func GetDatasource() *schema.Resource {
	r := GetResource()
	ds := &schema.Resource{
		Description: "Datasource for a single application via the unified Application v2 API.",
		Schema:      cloneSchemaAsComputed(r.Schema),
		ReadContext: read,
	}
	ds.Schema[AppIDVar] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The ID of the application.",
	}
	// project_id is exposed for state-shape parity with the resource but is
	// not required to look up an application: the v2 GetApplication RPC
	// keys off application_id alone.
	ds.Schema[ProjectIDVar] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The ID of the project the application belongs to. Optional on the datasource; the application is looked up by `app_id`.",
	}
	ds.Schema[helper.OrgIDVar] = helper.OrgIDDatasourceField

	// client_secret is only ever returned by Zitadel on CreateApplication;
	// GetApplication does not return it. Exposing it on the datasource
	// would imply the secret is readable, which it is not. Prune it from
	// the cloned OIDC and API nested schemas so the datasource surface
	// only advertises fields that can actually be populated from a Get.
	pruneClientSecret(ds.Schema, oidcBlockVar)
	pruneClientSecret(ds.Schema, apiBlockVar)
	return ds
}

func pruneClientSecret(s map[string]*schema.Schema, block string) {
	field, ok := s[block]
	if !ok {
		return
	}
	res, ok := field.Elem.(*schema.Resource)
	if !ok {
		return
	}
	// The package has a top-level `delete` CRUD function which shadows
	// the built-in `delete` identifier, so rebuild the schema map
	// without the unwanted key instead of calling the builtin.
	pruned := make(map[string]*schema.Schema, len(res.Schema))
	for k, v := range res.Schema {
		if k == clientSecretVar {
			continue
		}
		pruned[k] = v
	}
	res.Schema = pruned
}

// ListDatasources returns app IDs filtered by project_id and (optionally) name.
func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource listing application IDs in a project, regardless of application type.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project to list applications from.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by exact application name.",
			},
			appIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Matching application IDs.",
			},
		},
		ReadContext: list,
	}
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	// Scope the call to the org_id attribute so middleware metadata is set
	// consistently with the rest of the provider.
	ctx = helper.CtxWithOrgID(ctx, d)

	filters := make([]*apppb.ApplicationSearchFilter, 0, 2)
	if pid := d.Get(ProjectIDVar).(string); pid != "" {
		filters = append(filters, &apppb.ApplicationSearchFilter{
			Filter: &apppb.ApplicationSearchFilter_ProjectIdFilter{
				ProjectIdFilter: &apppb.ProjectIDFilter{ProjectId: pid},
			},
		})
	}
	if name := d.Get(NameVar).(string); name != "" {
		filters = append(filters, &apppb.ApplicationSearchFilter{
			Filter: &apppb.ApplicationSearchFilter_NameFilter{
				NameFilter: &apppb.ApplicationNameFilter{Name: name},
			},
		})
	}

	resp, err := client.ListApplications(ctx, &apppb.ListApplicationsRequest{Filters: filters})
	if err != nil {
		return diag.Errorf("failed to list applications: %v", err)
	}
	ids := make([]string, 0, len(resp.GetApplications()))
	for _, a := range resp.GetApplications() {
		ids = append(ids, a.GetApplicationId())
	}
	d.SetId("-")
	return diag.FromErr(d.Set(appIDsVar, ids))
}

// cloneSchemaAsComputed returns a deep copy of the resource schema flipped
// to a read-only (Computed=true) shape for use as a datasource. The
// recursion is needed because nested *schema.Resource Elems carry their own
// ExactlyOneOf / Required / MaxItems constraints that aren't valid on a
// fully-computed schema.
func cloneSchemaAsComputed(in map[string]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema, len(in))
	for k, v := range in {
		out[k] = cloneFieldAsComputed(v)
	}
	return out
}

func cloneFieldAsComputed(s *schema.Schema) *schema.Schema {
	dup := *s
	dup.Required = false
	dup.Optional = false
	dup.Default = nil
	dup.ValidateDiagFunc = nil
	dup.ValidateFunc = nil
	dup.ForceNew = false
	dup.ExactlyOneOf = nil
	dup.ConflictsWith = nil
	dup.AtLeastOneOf = nil
	dup.RequiredWith = nil
	dup.MaxItems = 0
	dup.MinItems = 0
	dup.DiffSuppressFunc = nil
	dup.Computed = true

	switch elem := dup.Elem.(type) {
	case *schema.Resource:
		dup.Elem = &schema.Resource{Schema: cloneSchemaAsComputed(elem.Schema)}
	case *schema.Schema:
		// Primitive list/set element: SDK requires only Type to be set, so
		// pass through unchanged.
		dup.Elem = elem
	}
	return &dup
}
