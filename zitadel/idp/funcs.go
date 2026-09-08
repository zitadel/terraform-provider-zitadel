package idp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	var queries []*admin.ProviderQuery
	if nameQuery := idp_utils.NameQuery(d); nameQuery != nil {
		queries = append(queries, &admin.ProviderQuery{
			Query: &admin.ProviderQuery_IdpNameQuery{IdpNameQuery: nameQuery},
		})
	}

	providers := make([]*idp.Provider, 0)
	for offset := uint64(0); ; offset += uint64(idp_utils.ListPageSize) {
		resp, err := client.ListProviders(ctx, &admin.ListProvidersRequest{
			Query:   idp_utils.ListQuery(offset),
			Queries: queries,
		})
		if err != nil {
			return diag.Errorf("failed to list idps: %v", err)
		}
		providers = append(providers, resp.GetResult()...)
		if len(resp.GetResult()) < int(idp_utils.ListPageSize) {
			break
		}
	}

	d.SetId("-")
	return diag.FromErr(d.Set(idp_utils.IdpsVar, idp_utils.FlattenProviders(d, providers)))
}
