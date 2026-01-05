package organization_metadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metadata "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/metadata/v2"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func set(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started set")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	key := d.Get(KeyVar).(string)
	value := []byte(d.Get(ValueVar).(string))
	orgID := d.Get(OrganizationIDVar).(string)

	_, err = client.SetOrganizationMetadata(ctx, &org.SetOrganizationMetadataRequest{
		OrganizationId: orgID,
		Metadata: []*org.Metadata{
			{
				Key:   key,
				Value: value,
			},
		},
	})
	if err != nil {
		return diag.Errorf("failed to set metadata: %v", err)
	}
	d.SetId(key)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	key := d.Id()
	orgID := d.Get(OrganizationIDVar).(string)

	resp, err := client.ListOrganizationMetadata(ctx, &org.ListOrganizationMetadataRequest{
		OrganizationId: orgID,
		Filters: []*metadata.MetadataSearchFilter{
			{
				Filter: &metadata.MetadataSearchFilter_KeyFilter{
					KeyFilter: &metadata.MetadataKeyFilter{
						Key: key,
					},
				},
			},
		},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get metadata: %v", err)
	}

	if len(resp.Metadata) == 0 {
		d.SetId("")
		return nil
	}

	remoteMetadata := resp.Metadata[0]
	value := string(remoteMetadata.Value)

	if err := d.Set(KeyVar, key); err != nil {
		return diag.Errorf("failed to set key: %v", err)
	}
	if err := d.Set(ValueVar, value); err != nil {
		return diag.Errorf("failed to set value: %v", err)
	}

	return nil
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteOrganizationMetadata(ctx, &org.DeleteOrganizationMetadataRequest{
		OrganizationId: d.Get(OrganizationIDVar).(string),
		Keys:           []string{d.Id()},
	})
	if err != nil {
		return diag.Errorf("failed to delete metadata: %v", err)
	}
	return nil
}

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	key := helper.GetID(d, KeyVar)
	orgID := helper.GetID(d, OrganizationIDVar)

	resp, err := client.ListOrganizationMetadata(ctx, &org.ListOrganizationMetadataRequest{
		OrganizationId: orgID,
		Filters: []*metadata.MetadataSearchFilter{
			{
				Filter: &metadata.MetadataSearchFilter_KeyFilter{
					KeyFilter: &metadata.MetadataKeyFilter{
						Key: key,
					},
				},
			},
		},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get metadata: %v", err)
	}

	if len(resp.Metadata) == 0 {
		d.SetId("")
		return nil
	}

	remoteMetadata := resp.Metadata[0]

	d.SetId(key)
	if err := d.Set(KeyVar, key); err != nil {
		return diag.Errorf("failed to set key: %v", err)
	}
	value := string(remoteMetadata.Value)
	if err := d.Set(ValueVar, value); err != nil {
		return diag.Errorf("failed to set value: %v", err)
	}

	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get(OrganizationIDVar).(string)
	keyFilter := d.Get(KeyVar).(string)

	req := &org.ListOrganizationMetadataRequest{
		OrganizationId: orgID,
	}

	if keyFilter != "" {
		req.Filters = []*metadata.MetadataSearchFilter{
			{
				Filter: &metadata.MetadataSearchFilter_KeyFilter{
					KeyFilter: &metadata.MetadataKeyFilter{
						Key: keyFilter,
					},
				},
			},
		}
	}

	resp, err := client.ListOrganizationMetadata(ctx, req)
	if err != nil {
		return diag.Errorf("failed to list metadata: %v", err)
	}

	metadataList := make([]interface{}, len(resp.Metadata))
	for i, meta := range resp.Metadata {
		metadataMap := map[string]interface{}{
			KeyVar:   meta.Key,
			ValueVar: string(meta.Value),
		}
		metadataList[i] = metadataMap
	}

	d.SetId(fmt.Sprintf("%s", orgID))
	return diag.FromErr(d.Set(metadataVar, metadataList))
}
