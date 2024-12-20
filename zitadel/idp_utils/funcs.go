package idp_utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

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
