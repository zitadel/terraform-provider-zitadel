package user_metadata

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

	userID := d.Get(UserIDVar).(string)
	key := d.Get(KeyVar).(string)
	value := []byte(d.Get(ValueVar).(string))
	_, err = client.SetUserMetadata(helper.CtxWithOrgID(ctx, d), &management.SetUserMetadataRequest{
		Id:    userID,
		Key:   key,
		Value: value,
	})
	if err != nil {
		return diag.Errorf("failed to set metadata entry: %v", err)
	}
	d.SetId(getUserMetadataID(userID, key))
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
	userID := d.Get(UserIDVar).(string)
	key := d.Get(KeyVar).(string)
	resp, err := client.GetUserMetadata(helper.CtxWithOrgID(ctx, d), &management.GetUserMetadataRequest{Id: userID, Key: key})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get metadata object")
	}
	set := map[string]interface{}{
		UserIDVar: userID,
		KeyVar:    key,
		ValueVar:  string(resp.GetMetadata().GetValue()),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set metadata with key %s: %v", k, err)
		}
	}
	d.SetId(getUserMetadataID(userID, key))
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
	userID := d.Get(UserIDVar).(string)
	key := d.Get(KeyVar).(string)
	_, err = client.RemoveUserMetadata(helper.CtxWithOrgID(ctx, d), &management.RemoveUserMetadataRequest{Id: userID, Key: key})
	if err != nil {
		return diag.Errorf("failed to remove metadata entry: %v", err)
	}
	return nil
}

func getUserMetadataID(userID string, key string) string {
	return userID + "_" + key
}
