package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/pkg/client/zitadel/project"
)

const (
	projectIdVar              = "id"
	projectOldIdVar           = "old_id"
	projectNameVar            = "name"
	projectState              = "state"
	projectResourceOwner      = "resource_owner"
	projectRoleAssertionVar   = "project_role_assertion"
	projectRoleCheckVar       = "project_role_check"
	hasProjectCheckVar        = "has_project_check"
	privateLabelingSettingVar = "private_labeling_setting"
)

func GetProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			projectIdVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project",
			},
			projectOldIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Old ID of the project",
				ForceNew:    true,
			},
			projectNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project",
			},
			projectResourceOwner: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization in which the project is located",
			},
			projectState: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State of the project",
			},
			projectRoleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "describes if roles of user should be added in token",
			},
			projectRoleCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ZITADEL checks if the user has at least one on this project",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			privateLabelingSettingVar: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Defines from where the private labeling should be triggered",
			},
		},
		DeleteContext: deleteProject,
		CreateContext: createProject,
		UpdateContext: updateProject,
		ReadContext:   readProject,
	}
}

func deleteProject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectResourceOwner).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProject(ctx, &management2.RemoveProjectRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete project: %v", err)
	}
	return nil
}

func updateProject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProject(ctx, &management2.UpdateProjectRequest{
		Id:                     d.Id(),
		Name:                   d.Get(projectNameVar).(string),
		ProjectRoleCheck:       d.Get(projectRoleCheckVar).(bool),
		ProjectRoleAssertion:   d.Get(projectRoleAssertionVar).(bool),
		HasProjectCheck:        d.Get(hasProjectCheckVar).(bool),
		PrivateLabelingSetting: d.Get(privateLabelingSettingVar).(project.PrivateLabelingSetting),
	})
	if err != nil {
		return diag.Errorf("failed to update project: %v", err)
	}
	return nil
}

func createProject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddProject(ctx, &management2.AddProjectRequest{
		Name:                   d.Get(projectNameVar).(string),
		ProjectRoleAssertion:   d.Get(projectRoleAssertionVar).(bool),
		ProjectRoleCheck:       d.Get(projectRoleCheckVar).(bool),
		HasProjectCheck:        d.Get(hasProjectCheckVar).(bool),
		PrivateLabelingSetting: d.Get(privateLabelingSettingVar).(project.PrivateLabelingSetting),
	})
	if err != nil {
		return diag.Errorf("failed to create project: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func readProject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProjectByID(ctx, &management2.GetProjectByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to read project: %v", err)
	}

	project := resp.GetProject()
	set := map[string]interface{}{
		projectIdVar:            project.GetId(),
		projectResourceOwner:    project.GetDetails().GetResourceOwner(),
		projectState:            project.GetState(),
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
