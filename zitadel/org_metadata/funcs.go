package org_metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func set(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	key := d.Get(KeyVar).(string)
	value, err := helper.Base64Decode(d.Get(ValueVar).(string))
	if err != nil {
		return diag.Errorf("failed to decode base64 value: %v", err)
	}
	_, err = client.SetOrgMetadata(helper.CtxWithOrgID(ctx, d), &management.SetOrgMetadataRequest{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return diag.Errorf("failed to set metadata entry: %v", err)
	}
	d.SetId(key)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	key := helper.GetID(d, KeyVar)
	resp, err := client.GetOrgMetadata(helper.CtxWithOrgID(ctx, d), &management.GetOrgMetadataRequest{Key: key})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get metadata object")
	}
	value := helper.Base64Encode(resp.GetMetadata().GetValue())
	set := map[string]interface{}{
		KeyVar:   key,
		ValueVar: value,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set metadata with key %s: %v", k, err)
		}
	}
	d.SetId(key)
	return nil
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.RemoveOrgMetadata(helper.CtxWithOrgID(ctx, d), &management.RemoveOrgMetadataRequest{Key: d.Id()})
	if err != nil {
		return diag.Errorf("failed to remove metadata entry: %v", err)
	}
	return nil
}
