package org_idp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	idpName := d.Get(NameVar).(string)
	idpNameMethod := d.Get(nameMethodVar).(string)
	idpType := d.Get(typeVar).(string)
	idpOwnerType := d.Get(ownerTypeVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListProvidersRequest{}
	if idpName != "" {
		req.Queries = append(req.Queries, &management.ProviderQuery{
			Query: &management.ProviderQuery_IdpNameQuery{
				IdpNameQuery: &idp.IDPNameQuery{
					Name:   idpName,
					Method: object.TextQueryMethod(object.TextQueryMethod_value[idpNameMethod]),
				},
			},
		})
	}
	if idpOwnerType != "" {
		req.Queries = append(req.Queries, &management.ProviderQuery{
			Query: &management.ProviderQuery_OwnerTypeQuery{
				OwnerTypeQuery: &idp.IDPOwnerTypeQuery{
					OwnerType: idp.IDPOwnerType(idp.IDPOwnerType_value[idpOwnerType]),
				},
			},
		})
	}
	resp, err := client.ListProviders(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while getting idp list: %v", err)
	}
	idpIDs := make([]string, 0, len(resp.Result))
	for _, idp := range resp.Result {
		if idpType != "" && idp.Type.String() != idpType {
			continue
		}
		idpIDs = append(idpIDs, idp.Id)
	}
	d.SetId("-")
	return diag.FromErr(d.Set(idpIDsVar, idpIDs))
}
