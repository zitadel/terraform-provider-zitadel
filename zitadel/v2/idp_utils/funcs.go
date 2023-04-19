package idp_utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
)

func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteProvider(ctx, &admin.DeleteProviderRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to delete idp: %v", err)
	}
	return nil
}

func ImportIDPWithSecret(secretVar string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		id := data.Id()
		if id == "" {
			return nil, fmt.Errorf("%s is not set", IdpIDVar)
		}
		parts := strings.SplitN(id, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("unexpected format of ID (%s), expected %s:%s", id, IdpIDVar, secretVar)
		}
		data.SetId(parts[0])
		if err := data.Set(secretVar, parts[1]); err != nil {
			return nil, err
		}
		return []*schema.ResourceData{data}, nil
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
