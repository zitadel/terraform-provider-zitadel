package idp_utils

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteProvider(ctx, &admin.DeleteProviderRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to delete idp: %v", err)
	}
	return nil
}

func StringValue(d *schema.ResourceData, attributeVar string) string {
	return d.Get(attributeVar).(string)
}

// WriteOnlyStringValue reads a write-only attribute from the raw config. See
// helper.WriteOnlyStringValue.
func WriteOnlyStringValue(d *schema.ResourceData, attributeVar string) string {
	return helper.WriteOnlyStringValue(d, attributeVar)
}

// ClientSecretHashDiff keeps the computed client_secret_hash attribute in sync
// with the write-only client_secret, so that rotating the secret is detected as
// a normal diff without the practitioner having to manage a version field.
func ClientSecretHashDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	return helper.WriteOnlyHashDiff(d, ClientSecretVar, ClientSecretHashVar)
}

func BoolValue(d *schema.ResourceData, attributeVar string) bool {
	return d.Get(attributeVar).(bool)
}

func ScopesValue(d *schema.ResourceData) []string {
	return helper.GetOkSetToStringSlice(d, ScopesVar)
}

func AutoLinkingValue(d *schema.ResourceData, attributeVar string) idp.AutoLinkingOption {
	return idp.AutoLinkingOption(idp.AutoLinkingOption_value[StringValue(d, attributeVar)])
}

func AutoLinkingString(value idp.AutoLinkingOption) string {
	return idp.AutoLinkingOption_name[int32(value)]
}

func ProviderOptionsValue(d *schema.ResourceData) *idp.Options {
	return &idp.Options{
		IsLinkingAllowed:  BoolValue(d, IsLinkingAllowedVar),
		IsCreationAllowed: BoolValue(d, IsCreationAllowedVar),
		IsAutoUpdate:      BoolValue(d, IsAutoUpdateVar),
		IsAutoCreation:    BoolValue(d, IsAutoCreationVar),
		AutoLinking:       AutoLinkingValue(d, AutoLinkingVar),
	}
}

func InterfaceToStringSlice(in interface{}) []string {
	slice := in.([]interface{})
	ret := make([]string, 0)
	for _, item := range slice {
		ret = append(ret, item.(string))
	}
	return ret
}

// ListPageSize is the number of identity providers requested per ListProviders call.
const ListPageSize uint32 = 100

// ListQuery returns the pagination query for the given offset.
func ListQuery(offset uint64) *object.ListQuery {
	return &object.ListQuery{
		Offset: offset,
		Limit:  ListPageSize,
		Asc:    true,
	}
}

// NameQuery builds the server-side name filter from the datasource configuration, or nil if no name is set.
func NameQuery(d *schema.ResourceData) *idp.IDPNameQuery {
	name, ok := d.GetOk(NameVar)
	if !ok {
		return nil
	}
	nameMethod := d.Get(NameMethodVar).(string)
	return &idp.IDPNameQuery{
		Name:   name.(string),
		Method: object.TextQueryMethod(object.TextQueryMethod_value[nameMethod]),
	}
}

// OwnerTypeQuery builds the server-side owner type filter from the datasource configuration,
// or nil if no owner type is set. IDP_OWNER_TYPE_UNSPECIFIED applies no filter.
func OwnerTypeQuery(d *schema.ResourceData) *idp.IDPOwnerTypeQuery {
	ownerType := d.Get(OwnerTypeVar).(string)
	if ownerType == "" || ownerType == idp.IDPOwnerType_IDP_OWNER_TYPE_UNSPECIFIED.String() {
		return nil
	}
	return &idp.IDPOwnerTypeQuery{
		OwnerType: idp.IDPOwnerType(idp.IDPOwnerType_value[ownerType]),
	}
}

// FlattenProviders applies the client-side type filter and maps the providers to the idps
// datasource schema, sorted by ID so list indexes stay stable across refreshes.
// PROVIDER_TYPE_UNSPECIFIED applies no filter.
func FlattenProviders(d *schema.ResourceData, providers []*idp.Provider) []interface{} {
	typeFilter := d.Get(TypeVar).(string)
	filterByType := typeFilter != "" && typeFilter != idp.ProviderType_PROVIDER_TYPE_UNSPECIFIED.String()
	idps := make([]interface{}, 0, len(providers))
	for _, provider := range providers {
		if filterByType && provider.GetType().String() != typeFilter {
			continue
		}
		idps = append(idps, map[string]interface{}{
			IdpIDVar:     provider.GetId(),
			NameVar:      provider.GetName(),
			TypeVar:      provider.GetType().String(),
			StateVar:     provider.GetState().String(),
			OwnerTypeVar: provider.GetOwner().String(),
		})
	}
	sort.Slice(idps, func(i, j int) bool {
		return idps[i].(map[string]interface{})[IdpIDVar].(string) < idps[j].(map[string]interface{})[IdpIDVar].(string)
	})
	return idps
}
