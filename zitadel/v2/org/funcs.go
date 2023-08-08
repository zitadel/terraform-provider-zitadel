package org

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
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
	client, err := helper.GetManagementClient(clientinfo, "")
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddOrg(ctx, &management.AddOrgRequest{
		Name: d.Get(nameVar).(string),
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
	client, err := helper.GetManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrg(ctx, &management.UpdateOrgRequest{
		Name: d.Get(nameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update org: %v", err)
	}
	return nil
}

func getByID(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started getByID")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	orgID := helper.GetID(d, orgIDVar)
	resp, err := client.GetOrgByID(ctx, &admin.GetOrgByIDRequest{
		Id: orgID,
	})
	if err != nil {
		return diag.Errorf("error while getting org by id %s: %v", orgID, err)
	}
	return diag.FromErr(setResourceState(d, resp.GetOrg()))
}

func queryDatasource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started queryDatasource")
	orgID := d.Get(orgIDVar).(string)
	orgName := d.Get(nameVar).(string)
	orgState := d.Get(stateVar).(string)
	orgDomain := d.Get(domainVar).(string)
	if orgID != "" && (orgName != "" || orgState != "" || orgDomain != "") {
		return diag.Errorf("only %s or one or many in %s, %s and %s are supported", orgIDVar, nameVar, stateVar, domainVar)
	}
	if orgID != "" {
		if err := getByID(ctx, d, m); err != nil {
			return err
		}
		return diag.FromErr(d.Set(orgIDVar, orgID))
	}
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
				NameQuery: &org.OrgNameQuery{Name: orgName},
			},
		})
	}
	if orgState != "" {
		req.Queries = append(req.Queries, &org.OrgQuery{
			Query: &org.OrgQuery_StateQuery{
				StateQuery: &org.OrgStateQuery{State: org.OrgState(org.OrgState_value[orgState])},
			},
		})
	}
	if orgDomain != "" {
		req.Queries = append(req.Queries, &org.OrgQuery{
			Query: &org.OrgQuery_DomainQuery{
				DomainQuery: &org.OrgDomainQuery{Domain: orgDomain},
			},
		})
	}
	if len(req.Queries) == 0 {
		return diag.Errorf("specify at least one filter")
	}
	resp, err := client.ListOrgs(ctx, req)
	if err != nil {
		return diag.Errorf("error while getting org by id %s: %v", orgName, err)
	}
	if len(resp.Result) != 1 {
		return diag.Errorf("the filters don't match exactly 1 org, but %d orgs", len(resp.Result))
	}
	if err = setResourceState(d, resp.Result[0]); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set(orgIDVar, resp.Result[0].Id))
}

func setResourceState(d *schema.ResourceData, remoteOrg *org.Org) error {
	d.SetId(remoteOrg.Id)
	if err := d.Set(nameVar, remoteOrg.Name); err != nil {
		return err
	}
	if err := d.Set(primaryDomainVar, remoteOrg.PrimaryDomain); err != nil {
		return err
	}
	if err := d.Set(stateVar, org.OrgState_name[int32(remoteOrg.State)]); err != nil {
		return err
	}
	return nil
}
