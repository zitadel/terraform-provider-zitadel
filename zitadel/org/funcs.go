package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org"
	orgv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteOrganization(ctx, &orgv2.DeleteOrganizationRequest{
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

	client, err := helper.GetOrgV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &orgv2.AddOrganizationRequest{
		Name: d.Get(NameVar).(string),
	}

	if orgID, ok := d.GetOk(OrgIDInputVar); ok {
		orgIDStr := orgID.(string)
		req.OrganizationId = &orgIDStr
	}

	if admins, ok := d.GetOk(adminsVar); ok {
		adminSet := admins.(*schema.Set)
		for _, admin := range adminSet.List() {
			adminMap := admin.(map[string]interface{})
			userId := adminMap["user_id"].(string)

			adminReq := &orgv2.AddOrganizationRequest_Admin{
				UserType: &orgv2.AddOrganizationRequest_Admin_UserId{
					UserId: userId,
				},
			}

			if roles, ok := adminMap["roles"]; ok && roles != nil {
				rolesList := roles.([]interface{})
				for _, role := range rolesList {
					adminReq.Roles = append(adminReq.Roles, role.(string))
				}
			}

			req.Admins = append(req.Admins, adminReq)
		}
	}

	resp, err := client.AddOrganization(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	orgId := resp.GetOrganizationId()
	d.SetId(orgId)

	if val, ok := d.GetOk(IsDefaultVar); ok && val.(bool) {
		adminClient, err := helper.GetAdminClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = adminClient.SetDefaultOrg(ctx, &admin.SetDefaultOrgRequest{
			OrgId: orgId,
		})
		if err != nil {
			return diag.Errorf("error while setting default org id %s: %v", orgId, err)
		}
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	if d.HasChange(NameVar) {
		client, err := helper.GetOrgV2Client(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateOrganization(ctx, &orgv2.UpdateOrganizationRequest{
			OrganizationId: d.Id(),
			Name:           d.Get(NameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update org: %v", err)
		}
	}

	if isDefault, ok := d.GetOk(IsDefaultVar); ok && isDefault.(bool) && d.HasChange(IsDefaultVar) {
		adminClient, err := helper.GetAdminClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = adminClient.SetDefaultOrg(ctx, &admin.SetDefaultOrgRequest{
			OrgId: d.Id(),
		})
		if err != nil {
			return diag.Errorf("error while setting default org id %s: %v", d.Id(), err)
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	orgID := helper.GetID(d, OrgIDVar)
	tflog.Info(ctx, fmt.Sprintf("Reading org ID: %s", orgID))
	resp, err := client.GetOrgByID(ctx, &admin.GetOrgByIDRequest{
		Id: orgID,
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
	tflog.Info(ctx, "Org found, updating state")
	remoteOrg := resp.GetOrg()
	d.SetId(remoteOrg.Id)
	if err := d.Set(NameVar, remoteOrg.Name); err != nil {
		return diag.Errorf("error while setting org name %s: %v", remoteOrg.Name, err)
	}
	if err := d.Set(primaryDomainVar, remoteOrg.PrimaryDomain); err != nil {
		return diag.Errorf("error while setting org primary domain %s: %v", remoteOrg.PrimaryDomain, err)
	}
	state := org.OrgState_name[int32(remoteOrg.State)]
	if err := d.Set(stateVar, state); err != nil {
		return diag.Errorf("error while setting org state %s: %v", state, err)
	}
	adminClient, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	defaultOrg, err := adminClient.GetDefaultOrg(ctx, &admin.GetDefaultOrgRequest{})
	if err != nil {
		return diag.Errorf("error while getting default instance org: %v", err)
	}
	if defaultOrg.Org.Id == remoteOrg.Id {
		if err := d.Set(IsDefaultVar, true); err != nil {
			return diag.Errorf("error while setting org is_default: %v", err)
		}
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
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &admin.ListOrgsRequest{}
	if orgName != "" {
		req.Queries = append(req.Queries, &org.OrgQuery{
			Query: &org.OrgQuery_NameQuery{
				NameQuery: &org.OrgNameQuery{
					Name:   orgName,
					Method: object.TextQueryMethod(object.TextQueryMethod_value[orgNameMethod]),
				},
			},
		})
	}
	if orgState != "" {
		req.Queries = append(req.Queries, &org.OrgQuery{
			Query: &org.OrgQuery_StateQuery{
				StateQuery: &org.OrgStateQuery{
					State: org.OrgState(org.OrgState_value[orgState]),
				},
			},
		})
	}
	if orgDomain != "" {
		req.Queries = append(req.Queries, &org.OrgQuery{
			Query: &org.OrgQuery_DomainQuery{
				DomainQuery: &org.OrgDomainQuery{
					Domain: orgDomain,
					Method: object.TextQueryMethod(object.TextQueryMethod_value[orgDomainMethod]),
				},
			},
		})
	}
	resp, err := client.ListOrgs(ctx, req)
	if err != nil {
		return diag.Errorf("error while getting org by id %s: %v", orgName, err)
	}
	orgIDs := make([]string, len(resp.Result))
	for i, org := range resp.Result {
		orgIDs[i] = org.Id
	}
	d.SetId("-")
	return diag.FromErr(d.Set(orgIDsVar, orgIDs))
}
