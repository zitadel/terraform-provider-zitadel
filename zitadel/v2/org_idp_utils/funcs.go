package org_idp_utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteProvider(ctx, &management.DeleteProviderRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to delete idp: %v", err)
	}
	return nil
}

func ImportIDPWithOrgAndClientSecret(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	id := data.Id()
	if id == "" {
		return nil, fmt.Errorf("%s is not set", IdpIDVar)
	}
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected orgid:idpid:clientsecret", id)
	}
	if err := data.Set(OrgIDVar, parts[0]); err != nil {
		return nil, err
	}
	data.SetId(parts[1])
	if err := data.Set(ClientSecretVar, parts[2]); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil

}
