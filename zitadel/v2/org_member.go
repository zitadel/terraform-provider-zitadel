package v2

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	orgMemberOrgIDVar  = "org_id"
	orgMemberUserIDVar = "user_id"
	orgMemberRolesVar  = "roles"
)

func GetOrgMember() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an organization, defined with the given role.",
		Schema: map[string]*schema.Schema{
			orgMemberOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the organization",
			},
			orgMemberUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			orgMemberRolesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of roles granted",
			},
		},
		DeleteContext: deleteOrgMember,
		CreateContext: createOrgMember,
		UpdateContext: updateOrgMember,
		ReadContext:   readOrgMember,
	}
}

func deleteOrgMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(orgMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveOrgMember(ctx, &management2.RemoveOrgMemberRequest{
		UserId: d.Get(orgMemberUserIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete orgmember: %v", err)
	}
	return nil
}

func updateOrgMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(orgMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrgMember(ctx, &management2.UpdateOrgMemberRequest{
		UserId: d.Get(orgMemberUserIDVar).(string),
		Roles:  d.Get(orgMemberRolesVar).([]string),
	})
	if err != nil {
		return diag.Errorf("failed to update orgmember: %v", err)
	}
	return nil
}

func createOrgMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(orgMemberUserIDVar).(string)
	roles := make([]string, 0)
	for _, role := range d.Get(orgMemberRolesVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}

	_, err = client.AddOrgMember(ctx, &management2.AddOrgMemberRequest{
		UserId: userID,
		Roles:  roles,
	})
	if err != nil {
		return diag.Errorf("failed to create orgmember: %v", err)
	}
	d.SetId(getOrgMemberID(org, userID))
	return nil
}

func readOrgMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	org := d.Get(orgMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgMembers(ctx, &management2.ListOrgMembersRequest{})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read orgmember: %v", err)
	}

	userID := d.Get(orgMemberUserIDVar).(string)
	for _, orgMember := range resp.Result {
		if orgMember.UserId == userID {
			set := map[string]interface{}{
				orgMemberUserIDVar: orgMember.GetUserId(),
				orgMemberOrgIDVar:  orgMember.GetDetails().GetResourceOwner(),
				orgMemberRolesVar:  orgMember.GetRoles(),
			}
			for k, v := range set {
				if err := d.Set(k, v); err != nil {
					return diag.Errorf("failed to set %s of orgmember: %v", k, err)
				}
			}
			d.SetId(getOrgMemberID(org, userID))
			return nil
		}
	}
	d.SetId("")
	return nil
}

func getOrgMemberID(org string, userID string) string {
	return org + "_" + userID
}

func splitOrgMemberID(orgMemberID string) (string, string) {
	parts := strings.Split(orgMemberID, "_")
	return parts[0], parts[1]
}
