package user_metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/metadata"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func set(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
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
	client, err := helper.GetManagementClient(ctx, clientinfo)
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
		return diag.Errorf("failed to get metadata object: %v", err)
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
	client, err := helper.GetManagementClient(ctx, clientinfo)
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

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := helper.GetID(d, UserIDVar)
	key := helper.GetID(d, KeyVar)

	resp, err := client.GetUserMetadata(helper.CtxWithOrgID(ctx, d), &management.GetUserMetadataRequest{Id: userID, Key: key})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get metadata: %v", err)
	}

	d.SetId(getUserMetadataID(userID, key))
	if err := d.Set(UserIDVar, userID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(KeyVar, key); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(ValueVar, string(resp.GetMetadata().GetValue())); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

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

	userID := helper.GetID(d, UserIDVar)
	keyFilter := d.Get(KeyVar).(string)

	req := &management.ListUserMetadataRequest{
		Id: userID,
	}

	if keyFilter != "" {
		req.Queries = []*metadata.MetadataQuery{
			{
				Query: &metadata.MetadataQuery_KeyQuery{
					KeyQuery: &metadata.MetadataKeyQuery{
						Key: keyFilter,
					},
				},
			},
		}
	}

	resp, err := client.ListUserMetadata(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to list metadata: %v", err)
	}

	metadataList := make([]interface{}, len(resp.Result))
	for i, meta := range resp.Result {
		metadataMap := map[string]interface{}{
			KeyVar:   meta.Key,
			ValueVar: string(meta.Value),
		}
		metadataList[i] = metadataMap
	}

	d.SetId(userID)
	return diag.FromErr(d.Set(metadataVar, metadataList))
}
