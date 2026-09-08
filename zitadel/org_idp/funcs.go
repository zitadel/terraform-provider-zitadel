package org_idp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	idppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp"
)

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	var queries []*management.ProviderQuery
	if nameQuery := idp.NameQuery(d); nameQuery != nil {
		queries = append(queries, &management.ProviderQuery{
			Query: &management.ProviderQuery_IdpNameQuery{IdpNameQuery: nameQuery},
		})
	}
	if ownerType, ok := d.GetOk(ownerTypeVar); ok {
		queries = append(queries, &management.ProviderQuery{
			Query: &management.ProviderQuery_OwnerTypeQuery{
				OwnerTypeQuery: &idppb.IDPOwnerTypeQuery{
					OwnerType: idppb.IDPOwnerType(idppb.IDPOwnerType_value[ownerType.(string)]),
				},
			},
		})
	}

	providers := make([]*idppb.Provider, 0)
	for offset := uint64(0); ; offset += uint64(idp.PageSize()) {
		resp, err := client.ListProviders(helper.CtxWithOrgID(ctx, d), &management.ListProvidersRequest{
			Query:   idp.ListQuery(offset),
			Queries: queries,
		})
		if err != nil {
			return diag.Errorf("failed to list idps: %v", err)
		}
		providers = append(providers, resp.GetResult()...)
		if len(resp.GetResult()) < idp.PageSize() {
			break
		}
	}

	d.SetId("-")
	return diag.FromErr(d.Set(idpsVar, idp.Flatten(d, providers)))
}
