package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	projectIdVar              = "id"
	projectNameVar            = "name"
	projectState              = "state"
	projectResourceOwner      = "resource_owner"
	projectRoleAssertionVar   = "project_role_assertion"
	projectRoleCheckVar       = "project_role_check"
	hasProjectCheckVar        = "has_project_check"
	privateLabelingSettingVar = "private_labeling_setting"
)

func GetProjectDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			projectIdVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project",
			},
			projectNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the project",
			},
			projectResourceOwner: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization in which the project is located",
			},
			projectState: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "State of the project",
			},
			projectRoleAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "describes if roles of user should be added in token",
			},
			projectRoleCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the user has at least one on this project",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			privateLabelingSettingVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Defines from where the private labeling should be triggered",
			},
		},
	}
}

func readProjectsOfOrg(ctx context.Context, projects *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListProjects(ctx, &management2.ListProjectsRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of projects: %v", err)
	}

	projectResource := GetProjectDatasource()
	for i := range resp.Result {
		project := resp.Result[i]
		projectdata := projectResource.Data(&terraform.InstanceState{})
		projectdata.SetId(project.Id)
		if errDiag := readProject(ctx, projectdata, m, clientinfo, org); errDiag != nil {
			return errDiag
		}

		data := resourceToValueMap(projectResource, projectdata)
		projects.Add(data)
	}
	return nil
}

func readProject(ctx context.Context, d *schema.ResourceData, m interface{}, info *ClientInfo, org string) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	client, err := getManagementClient(info, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProjectByID(ctx, &management2.GetProjectByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to get list of users: %v", err)
	}

	project := resp.GetProject()
	set := map[string]interface{}{
		projectIdVar:            project.GetId(),
		projectResourceOwner:    project.GetDetails().GetResourceOwner(),
		projectState:            project.GetState().Number(),
		projectNameVar:          project.GetName(),
		projectRoleAssertionVar: project.GetProjectRoleAssertion(),
		projectRoleCheckVar:     project.GetProjectRoleCheck(),
		hasProjectCheckVar:      project.GetHasProjectCheck(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	d.SetId(project.GetId())
	return nil
}

func getProjectValueMap(d *schema.ResourceData) map[string]interface{} {
	res := GetProjectDatasource()

	values := make(map[string]interface{}, 0)
	for key := range res.Schema {
		values[key] = d.Get(key)
	}
	return values
}
