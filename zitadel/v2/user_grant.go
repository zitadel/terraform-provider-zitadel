package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	userGrantProjectIDVar      = "project_id"
	userGrantProjectGrantIDVar = "project_grant_id"
	userGrantUserIDVar         = "user_id"
	userGrantRoleKeysVar       = "role_keys"
	userGrantOrgIDVar          = "org_id"
)

func GetUserGrant() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the authorization given to a user directly, including the given roles.",
		Schema: map[string]*schema.Schema{
			userGrantProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			userGrantProjectGrantIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the granted project",
				ForceNew:    true,
			},
			userGrantUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			userGrantRoleKeysVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of roles granted",
			},
			userGrantOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
			},
		},
		DeleteContext: deleteUserGrant,
		CreateContext: createUserGrant,
		UpdateContext: updateUserGrant,
		ReadContext:   readUserGrant,
	}
}

func deleteUserGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUserGrant(ctx, &management2.RemoveUserGrantRequest{
		GrantId: d.Id(),
		UserId:  d.Get(userGrantUserIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete usergrant: %v", err)
	}
	return nil
}

func updateUserGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserGrant(ctx, &management2.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(userGrantUserIDVar).(string),
		RoleKeys: d.Get(userGrantRoleKeysVar).([]string),
	})
	if err != nil {
		return diag.Errorf("failed to update usergrant: %v", err)
	}
	return nil
}

func createUserGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddUserGrant(ctx, &management2.AddUserGrantRequest{
		UserId:         d.Get(userGrantUserIDVar).(string),
		ProjectGrantId: d.Get(userGrantProjectGrantIDVar).(string),
		ProjectId:      d.Get(userGrantProjectIDVar).(string),
		RoleKeys:       d.Get(userGrantRoleKeysVar).([]string),
	})
	if err != nil {
		return diag.Errorf("failed to create usergrant: %v", err)
	}
	d.SetId(resp.GetUserGrantId())
	return nil
}

func readUserGrant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetUserGrantByID(ctx, &management2.GetUserGrantByIDRequest{UserId: d.Get(userGrantUserIDVar).(string), GrantId: d.Id()})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read usergrant: %v", err)
	}

	grant := resp.GetUserGrant()
	set := map[string]interface{}{
		userGrantUserIDVar:         grant.GetProjectId(),
		userGrantProjectIDVar:      grant.GetProjectId(),
		userGrantProjectGrantIDVar: grant.GetProjectGrantId(),
		userGrantRoleKeysVar:       grant.GetRoleKeys(),
		userGrantOrgIDVar:          grant.GetDetails().GetResourceOwner(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of projectgrant: %v", k, err)
		}
	}
	d.SetId(grant.GetId())
	return nil
}
