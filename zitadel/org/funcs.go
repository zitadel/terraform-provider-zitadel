package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	objectv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object/v2"
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.RemoveOrg(ctx, &admin.RemoveOrgRequest{
		OrgId: d.Id(),
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
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddOrg(ctx, &management.AddOrgRequest{
		Name: d.Get(NameVar).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	orgId := resp.GetId()
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
	// If try updating the name to the same value API will return an error.
	if d.HasChange(NameVar) {
		client, err := helper.GetManagementClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateOrg(helper.CtxSetOrgID(ctx, d.Id()), &management.UpdateOrgRequest{
			Name: d.Get(NameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update org: %v", err)
		}
	}
	// To unset the default org, we need to set another org as default org.
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

	// Map old v1 state values to v2 for backwards compatibility
	if orgState != "" {
		stateMapping := map[string]string{
			"ORG_STATE_UNSPECIFIED": "ORGANIZATION_STATE_UNSPECIFIED",
			"ORG_STATE_ACTIVE":      "ORGANIZATION_STATE_ACTIVE",
			"ORG_STATE_INACTIVE":    "ORGANIZATION_STATE_INACTIVE",
			"ORG_STATE_REMOVED":     "ORGANIZATION_STATE_REMOVED",
		}
		if mappedState, ok := stateMapping[orgState]; ok {
			orgState = mappedState
		}
	}

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &orgv2.ListOrganizationsRequest{}
	if orgName != "" {
		req.Queries = append(req.Queries, &orgv2.SearchQuery{
			Query: &orgv2.SearchQuery_NameQuery{
				NameQuery: &orgv2.OrganizationNameQuery{
					Name:   orgName,
					Method: objectv2.TextQueryMethod(objectv2.TextQueryMethod_value[orgNameMethod]),
				},
			},
		})
	}
	if orgState != "" {
		req.Queries = append(req.Queries, &orgv2.SearchQuery{
			Query: &orgv2.SearchQuery_StateQuery{
				StateQuery: &orgv2.OrganizationStateQuery{
					State: orgv2.OrganizationState(orgv2.OrganizationState_value[orgState]),
				},
			},
		})
	}
	if orgDomain != "" {
		req.Queries = append(req.Queries, &orgv2.SearchQuery{
			Query: &orgv2.SearchQuery_DomainQuery{
				DomainQuery: &orgv2.OrganizationDomainQuery{
					Domain: orgDomain,
					Method: objectv2.TextQueryMethod(objectv2.TextQueryMethod_value[orgDomainMethod]),
				},
			},
		})
	}
	resp, err := client.ListOrganizations(ctx, req)
	if err != nil {
		return diag.Errorf("error while listing orgs (name=%q, domain=%q, state=%q): %v", orgName, orgDomain, orgState, err)
	}
	orgIDs := make([]string, len(resp.Result))
	for i, org := range resp.Result {
		orgIDs[i] = org.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(orgIDsVar, orgIDs))
}
