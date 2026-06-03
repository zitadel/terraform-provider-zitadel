package project

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	filter "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"

	projectpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetProjectV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteProject(ctx, &projectpb.DeleteProjectRequest{
		ProjectId: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete project: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetProjectV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectRoleAssertion := d.Get(roleAssertionVar).(bool)
	authRequired := d.Get(roleCheckVar).(bool)
	projectAccessRequired := d.Get(hasProjectCheckVar).(bool)
	plSetting := projectpb.PrivateLabelingSetting(projectpb.PrivateLabelingSetting_value[d.Get(privateLabelingSettingVar).(string)])

	name := d.Get(NameVar).(string)

	_, err = client.UpdateProject(ctx, &projectpb.UpdateProjectRequest{
		ProjectId:                d.Id(),
		Name:                     &name,
		ProjectRoleAssertion:     &projectRoleAssertion,
		AuthorizationRequired:    &authRequired,
		ProjectAccessRequired:    &projectAccessRequired,
		PrivateLabelingSetting:   &plSetting,
	})
	if err != nil {
		return diag.Errorf("failed to update project: %v", err)
	}

	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetProjectV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	plSetting := projectpb.PrivateLabelingSetting(projectpb.PrivateLabelingSetting_value[d.Get(privateLabelingSettingVar).(string)])

	resp, err := client.CreateProject(ctx, &projectpb.CreateProjectRequest{
		OrganizationId:         d.Get(helper.OrgIDVar).(string),
		Name:                   d.Get(NameVar).(string),
		ProjectRoleAssertion:   d.Get(roleAssertionVar).(bool),
		AuthorizationRequired:  d.Get(roleCheckVar).(bool),
		ProjectAccessRequired:  d.Get(hasProjectCheckVar).(bool),
		PrivateLabelingSetting: plSetting,
	})
	if err != nil {
		return diag.Errorf("failed to create project: %v", err)
	}
	d.SetId(resp.GetProjectId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetProjectV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProject(ctx, &projectpb.GetProjectRequest{ProjectId: helper.GetID(d, ProjectIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get project: %v", err)
	}

	p := resp.GetProject()
	set := map[string]interface{}{
		helper.OrgIDVar:           p.GetOrganizationId(),
		stateVar:                  p.GetState().String(),
		NameVar:                   p.GetName(),
		roleAssertionVar:          p.GetProjectRoleAssertion(),
		roleCheckVar:              p.GetAuthorizationRequired(),
		hasProjectCheckVar:        p.GetProjectAccessRequired(),
		privateLabelingSettingVar: p.GetPrivateLabelingSetting().String(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	d.SetId(p.GetProjectId())

	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	name := d.Get(NameVar).(string)
	nameMethod := d.Get(nameMethodVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetProjectV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &projectpb.ListProjectsRequest{
		Filters: make([]*projectpb.ProjectSearchFilter, 0),
	}

	orgID := d.Get(helper.OrgIDVar).(string)
	if orgID != "" {
		req.Filters = append(req.Filters, &projectpb.ProjectSearchFilter{
			Filter: &projectpb.ProjectSearchFilter_OrganizationIdFilter{
				OrganizationIdFilter: &projectpb.ProjectOrganizationIDFilter{
					OrganizationId: orgID,
					Type:           projectpb.ProjectOrganizationIDFilter_OWNED,
				},
			},
		})
	}

	if name != "" {
		// Convert V1 text query method names to V2 filter method names
		v2MethodStr := strings.Replace(nameMethod, "TEXT_QUERY_METHOD", "TEXT_FILTER_METHOD", 1)
		v2Method := filter.TextFilterMethod(filter.TextFilterMethod_value[v2MethodStr])

		req.Filters = append(req.Filters,
			&projectpb.ProjectSearchFilter{
				Filter: &projectpb.ProjectSearchFilter_ProjectNameFilter{
					ProjectNameFilter: &projectpb.ProjectNameFilter{
						ProjectName: name,
						Method:      v2Method,
					},
				},
			})
	}

	resp, err := client.ListProjects(ctx, req)
	if err != nil {
		return diag.Errorf("error while getting project by name %s: %v", name, err)
	}
	ids := make([]string, len(resp.GetProjects()))
	for i, res := range resp.GetProjects() {
		ids[i] = res.GetProjectId()
	}
	d.SetId("-")
	return diag.FromErr(d.Set(projectIDsVar, ids))
}
