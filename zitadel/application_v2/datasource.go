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
	ds.Schema[ProjectIDVar] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The ID of the project the application belongs to.",
	}
	ds.Schema[helper.OrgIDVar] = helper.OrgIDDatasourceField
	return ds
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
	dup.ForceNew = false
	dup.ExactlyOneOf = nil
	dup.ConflictsWith = nil
	dup.AtLeastOneOf = nil
	dup.RequiredWith = nil
	dup.MaxItems = 0
	dup.MinItems = 0
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
