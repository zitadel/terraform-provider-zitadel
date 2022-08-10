package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	projectGrantOrgIDVar        = "org_id"
	projectGrantProjectIDVar    = "project_id"
	projectGrantGrantedOrgIDVar = "granted_org_id"
	projectGrantRoleKeysVar     = "role_keys"
)

func GetProjectGrant() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the grant of a project to a different organization, also containing the available roles which can be given to the members of the projectgrant.",
		Schema: map[string]*schema.Schema{
			projectGrantProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			projectGrantGrantedOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization granted the project",
				ForceNew:    true,
			},
			projectGrantRoleKeysVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of roles granted",
			},
			projectGrantOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
			},
		},
		DeleteContext: deleteProjectGrant,
		CreateContext: createProjectGrant,
		UpdateContext: updateProjectGrant,
		ReadContext:   readProjectGrant,
	}
}

func deleteProjectGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProjectGrant(ctx, &management2.RemoveProjectGrantRequest{
		GrantId:   d.Id(),
		ProjectId: d.Get(projectGrantProjectIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete projectgrant: %v", err)
	}
	return nil
}

func updateProjectGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectGrant(ctx, &management2.UpdateProjectGrantRequest{
		GrantId:   d.Id(),
		ProjectId: d.Get(projectGrantProjectIDVar).(string),
		RoleKeys:  d.Get(projectGrantRoleKeysVar).([]string),
	})
	if err != nil {
		return diag.Errorf("failed to update projectgrant: %v", err)
	}
	return nil
}

func createProjectGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	roles := make([]string, 0)
	for _, role := range d.Get(projectGrantRoleKeysVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}

	resp, err := client.AddProjectGrant(ctx, &management2.AddProjectGrantRequest{
		GrantedOrgId: d.Get(projectGrantGrantedOrgIDVar).(string),
		ProjectId:    d.Get(projectGrantProjectIDVar).(string),
		RoleKeys:     roles,
	})
	if err != nil {
		return diag.Errorf("failed to create projectgrant: %v", err)
	}
	d.SetId(resp.GetGrantId())
	return nil
}

func readProjectGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProjectGrantByID(ctx, &management2.GetProjectGrantByIDRequest{ProjectId: d.Get(projectGrantProjectIDVar).(string), GrantId: d.Id()})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read projectgrant: %v", err)
	}

	projectGrant := resp.GetProjectGrant()
	set := map[string]interface{}{
		projectGrantProjectIDVar:    projectGrant.GetProjectId(),
		projectGrantGrantedOrgIDVar: projectGrant.GetGrantedOrgId(),
		projectGrantRoleKeysVar:     projectGrant.GetGrantedRoleKeys(),
		projectGrantOrgIDVar:        projectGrant.GetDetails().GetResourceOwner(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of projectgrant: %v", k, err)
		}
	}
	d.SetId(projectGrant.GetGrantId())
	return nil
}
