package organization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	object "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object/v2"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
		OrganizationId: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &org.AddOrganizationRequest{
		Name: d.Get(NameVar).(string),
	}
	if v, ok := d.GetOk(OrganizationIDVar); ok {
		id := v.(string)
		req.OrganizationId = &id
	}

	if v, ok := d.GetOk(adminsVar); ok {
		adminsList := v.([]interface{})
		admins := make([]*org.AddOrganizationRequest_Admin, len(adminsList))
		for i, adminRaw := range adminsList {
			adminMap := adminRaw.(map[string]interface{})
			admin := &org.AddOrganizationRequest_Admin{
				UserType: &org.AddOrganizationRequest_Admin_UserId{
					UserId: adminMap[adminUserIDVar].(string),
				},
			}
			if rolesRaw, ok := adminMap[adminRolesVar]; ok && rolesRaw != nil {
				rolesList := rolesRaw.([]interface{})
				roles := make([]string, len(rolesList))
				for j, role := range rolesList {
					roles[j] = role.(string)
				}
				admin.Roles = roles
			}
			admins[i] = admin
		}
		req.Admins = admins
	}

	resp, err := client.AddOrganization(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetOrganizationId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	if d.HasChange(NameVar) {
		client, err := helper.GetOrgClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateOrganization(ctx, &org.UpdateOrganizationRequest{
			OrganizationId: d.Id(),
			Name:           d.Get(NameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update org: %v", err)
		}
	}
	return nil
}

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	orgID := helper.GetID(d, OrgIDVar)
	tflog.Info(ctx, fmt.Sprintf("Reading org ID: %s", orgID))
	resp, err := client.ListOrganizations(ctx, &org.ListOrganizationsRequest{
		Queries: []*org.SearchQuery{
			{
				Query: &org.SearchQuery_IdQuery{
					IdQuery: &org.OrganizationIDQuery{
						Id: orgID,
					},
				},
			},
		},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		tflog.Info(ctx, "Org not found, clearing from state")
		d.SetId("")
		return nil
	}
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error getting org: %v", err))
		return diag.Errorf("error while getting org by id %s: %v", orgID, err)
	}

	if len(resp.Result) == 0 {
		tflog.Info(ctx, "Org not found in list, clearing from state")
		d.SetId("")
		return nil
	}
	remoteOrg := resp.Result[0]

	tflog.Info(ctx, "Org found, updating state")
	d.SetId(remoteOrg.Id)
	if err := d.Set(NameVar, remoteOrg.Name); err != nil {
		return diag.Errorf("error while setting org name %s: %v", remoteOrg.Name, err)
	}
	if err := d.Set(primaryDomainVar, remoteOrg.PrimaryDomain); err != nil {
		return diag.Errorf("error while setting org primary domain %s: %v", remoteOrg.PrimaryDomain, err)
	}
	state := org.OrganizationState_name[int32(remoteOrg.State)]
	if err := d.Set(stateVar, state); err != nil {
		return diag.Errorf("error while setting org state %s: %v", state, err)
	}
	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	orgName := d.Get(NameVar).(string)
	orgNameMethod := d.Get(nameMethodVar).(string)
	orgDomain := d.Get(DomainVar).(string)
	orgDomainMethod := d.Get(domainMethodVar).(string)
	orgState := d.Get(stateVar).(string)
	isDefault := d.Get(isDefaultVar).(bool)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &org.ListOrganizationsRequest{}
	if orgName != "" {
		req.Queries = append(req.Queries, &org.SearchQuery{
			Query: &org.SearchQuery_NameQuery{
				NameQuery: &org.OrganizationNameQuery{
					Name:   orgName,
					Method: object.TextQueryMethod(object.TextQueryMethod_value[orgNameMethod]),
				},
			},
		})
	}
	if orgState != "" {
		req.Queries = append(req.Queries, &org.SearchQuery{
			Query: &org.SearchQuery_StateQuery{
				StateQuery: &org.OrganizationStateQuery{
					State: org.OrganizationState(org.OrganizationState_value[orgState]),
				},
			},
		})
	}
	if orgDomain != "" {
		req.Queries = append(req.Queries, &org.SearchQuery{
			Query: &org.SearchQuery_DomainQuery{
				DomainQuery: &org.OrganizationDomainQuery{
					Domain: orgDomain,
					Method: object.TextQueryMethod(object.TextQueryMethod_value[orgDomainMethod]),
				},
			},
		})
	}
	if isDefault {
		req.Queries = append(req.Queries, &org.SearchQuery{
			Query: &org.SearchQuery_DefaultQuery{
				DefaultQuery: &org.DefaultOrganizationQuery{},
			},
		})
	}
	resp, err := client.ListOrganizations(ctx, req)
	if err != nil {
		return diag.Errorf("error while getting org list: %v", err)
	}
	orgIDs := make([]string, len(resp.Result))
	for i, org := range resp.Result {
		orgIDs[i] = org.Id
	}
	d.SetId("-")
	return diag.FromErr(d.Set(orgIDsVar, orgIDs))
}
