package idp

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
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
	if nameQuery := NameQuery(d); nameQuery != nil {
		queries = append(queries, &admin.ProviderQuery{
			Query: &admin.ProviderQuery_IdpNameQuery{IdpNameQuery: nameQuery},
		})
	}
	providers := make([]*idp.Provider, 0)
	for offset := uint64(0); ; offset += uint64(listPageSize) {
		resp, err := client.ListProviders(ctx, &admin.ListProvidersRequest{
			Query:   ListQuery(offset),
			Queries: queries,
		})
		if err != nil {
			return diag.Errorf("failed to list idps: %v", err)
		}
		providers = append(providers, resp.GetResult()...)
		if len(resp.GetResult()) < int(listPageSize) {
			break
		}
	}

	d.SetId("-")
	return diag.FromErr(d.Set(idpsVar, Flatten(d, providers)))
}

// NameQuery builds the server-side name filter from the datasource configuration, or nil if no name is set.
func NameQuery(d *schema.ResourceData) *idp.IDPNameQuery {
	name, ok := d.GetOk(NameVar)
	if !ok {
		return nil
	}
	nameMethod := d.Get(nameMethodVar).(string)
	return &idp.IDPNameQuery{
		Name:   name.(string),
		Method: object.TextQueryMethod(object.TextQueryMethod_value[nameMethod]),
	}
}

// ListQuery returns the pagination query for the given offset.
func ListQuery(offset uint64) *object.ListQuery {
	return &object.ListQuery{
		Offset: offset,
		Limit:  listPageSize,
		Asc:    true,
	}
}

// PageSize returns the number of identity providers requested per page.
func PageSize() int {
	return int(listPageSize)
}

// Flatten applies the client-side type filter and maps the providers to the nested schema,
// sorted by ID so list indexes stay stable across refreshes.
func Flatten(d *schema.ResourceData, providers []*idp.Provider) []interface{} {
	typeFilter, filterByType := d.GetOk(TypeVar)
	idps := make([]interface{}, 0, len(providers))
	for _, provider := range providers {
		if filterByType && provider.GetType().String() != typeFilter.(string) {
			continue
		}
		idps = append(idps, map[string]interface{}{
			idVar:        provider.GetId(),
			NameVar:      provider.GetName(),
			TypeVar:      provider.GetType().String(),
			stateVar:     provider.GetState().String(),
			ownerTypeVar: provider.GetOwner().String(),
		})
	}
	sort.Slice(idps, func(i, j int) bool {
		return idps[i].(map[string]interface{})[idVar].(string) < idps[j].(map[string]interface{})[idVar].(string)
	})
	return idps
}
