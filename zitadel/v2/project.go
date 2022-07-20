package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"
)

const (
	projectIdVar                     = "id"
	projectNameVar                   = "name"
	projectState                     = "state"
	projectResourceOwner             = "resource_owner"
	projectRoleAssertionVar          = "project_role_assertion"
	projectRoleCheckVar              = "project_role_check"
	projectHasProjectCheckVar        = "has_project_check"
	projectPrivateLabelingSettingVar = "private_labeling_setting"
	projectRoles                     = "roles"
	projectRoleKey                   = "role_key"
	projectRoleDisplayName           = "display_name"
	projectRoleGroup                 = "group"
)

func GetProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			projectIdVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project",
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
				Computed:    true,
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
			projectHasProjectCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			projectPrivateLabelingSettingVar: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Defines from where the private labeling should be triggered",
			},
			projectRoles: {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						projectRoleKey: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key used for project role",
						},
						projectRoleDisplayName: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name used for project role",
						},
						projectRoleGroup: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Group used for project role",
						},
					},
				},
				Required:    true,
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

	client, err := getManagementClient(clientinfo, d.Get(projectResourceOwner).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	plSetting := d.Get(projectPrivateLabelingSettingVar).(int)
	_, err = client.UpdateProject(ctx, &management2.UpdateProjectRequest{
		Id:                     d.Id(),
		Name:                   d.Get(projectNameVar).(string),
		ProjectRoleCheck:       d.Get(projectRoleCheckVar).(bool),
		ProjectRoleAssertion:   d.Get(projectRoleAssertionVar).(bool),
		HasProjectCheck:        d.Get(projectHasProjectCheckVar).(bool),
		PrivateLabelingSetting: project.PrivateLabelingSetting(plSetting),
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

	client, err := getManagementClient(clientinfo, d.Get(projectResourceOwner).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	plSetting := d.Get(projectPrivateLabelingSettingVar).(int)
	resp, err := client.AddProject(ctx, &management2.AddProjectRequest{
		Name:                   d.Get(projectNameVar).(string),
		ProjectRoleAssertion:   d.Get(projectRoleAssertionVar).(bool),
		ProjectRoleCheck:       d.Get(projectRoleCheckVar).(bool),
		HasProjectCheck:        d.Get(projectHasProjectCheckVar).(bool),
		PrivateLabelingSetting: project.PrivateLabelingSetting(plSetting),
	})
	if err != nil {
		return diag.Errorf("failed to create project: %v", err)
	}
	d.SetId(resp.GetId())

	roles := d.Get(projectRoles)
	if roles != nil {
		rolesList := roles.(*schema.Set)
		for _, role := range rolesList.List() {
			roleMap := role.(map[string]interface{})
			_, err := client.AddProjectRole(ctx, &management2.AddProjectRoleRequest{
				ProjectId:   resp.GetId(),
				RoleKey:     roleMap[projectRoleKey].(string),
				DisplayName: roleMap[projectRoleDisplayName].(string),
				Group:       roleMap[projectRoleGroup].(string),
			})
			if err != nil {
				return diag.Errorf("failed to create projectrole: %v", err)
			}
		}
	}

	return nil
}

func readProject(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectResourceOwner).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProjectByID(ctx, &management2.GetProjectByIDRequest{Id: d.Id()})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read project: %v", err)
	}

	project := resp.GetProject()
	set := map[string]interface{}{
		projectIdVar:              project.GetId(),
		projectResourceOwner:      project.GetDetails().GetResourceOwner(),
		projectState:              project.GetState().String(),
		projectNameVar:            project.GetName(),
		projectRoleAssertionVar:   project.GetProjectRoleAssertion(),
		projectRoleCheckVar:       project.GetProjectRoleCheck(),
		projectHasProjectCheckVar: project.GetHasProjectCheck(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	d.SetId(project.GetId())

	respRoles, err := client.ListProjectRoles(ctx, &management2.ListProjectRolesRequest{
		ProjectId: project.GetId(),
	})

	schemaFunc := func(elem interface{}) int {
		elemMap := elem.(map[string]string)
		return schema.HashString(elemMap[projectRoleKey])
	}
	roles := schema.NewSet(schemaFunc, []interface{}{})
	for _, role := range respRoles.Result {
		roleData := map[string]string{
			projectRoleKey:         role.GetKey(),
			projectRoleDisplayName: role.GetDisplayName(),
			projectRoleGroup:       role.GetGroup(),
		}
		roles.Add(roleData)
	}

	return nil
}
