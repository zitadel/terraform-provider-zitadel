package org

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
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
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddOrg(ctx, &management.AddOrgRequest{
		Name: d.Get(NameVar).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrg(helper.CtxSetOrgID(ctx, d.Id()), &management.UpdateOrgRequest{
		Name: d.Get(NameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update org: %v", err)
	}
	return nil
}

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	orgID := helper.GetID(d, OrgIDVar)
	resp, err := client.GetOrgByID(ctx, &admin.GetOrgByIDRequest{
		Id: orgID,
	})
	if err != nil {
		return diag.Errorf("error while getting org by id %s: %v", orgID, err)
	}
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
	client, err := helper.GetAdminClient(clientinfo)
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
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(orgIDsVar, orgIDs))
}
