package v2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	appKeyOrgIDVar          = "org_id"
	appKeyProjectIDVar      = "project_id"
	appKeyAppIDVar          = "app_id"
	appKeyKeyTypeVar        = "key_type"
	appKeyKeyDetailsVar     = "key_details"
	appKeyExpirationDateVar = "expiration_date"
)

func GetAppKey() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a app key",
		Schema: map[string]*schema.Schema{
			appKeyOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			appKeyProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			appKeyAppIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the application",
				ForceNew:    true,
			},
			appKeyKeyTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the app key",
				ForceNew:    true,
			},
			appKeyExpirationDateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration date of the app key",
				ForceNew:    true,
			},
			appKeyKeyDetailsVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the app key",
				Sensitive:   true,
			},
		},
		DeleteContext: deleteAppKey,
		CreateContext: createAppKey,
		ReadContext:   readAppKey,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func deleteAppKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(appKeyOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveAppKey(ctx, &management2.RemoveAppKeyRequest{
		ProjectId: d.Get(appKeyProjectIDVar).(string),
		AppId:     d.Get(appKeyAppIDVar).(string),
		KeyId:     d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete app key: %v", err)
	}
	return nil
}

func createAppKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(appKeyOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	t, err := time.Parse(time.RFC3339, d.Get(appKeyExpirationDateVar).(string))
	if err != nil {
		return diag.Errorf("failed to parse time: %v", err)
	}

	keyType := d.Get(appKeyKeyTypeVar).(string)
	resp, err := client.AddAppKey(ctx, &management2.AddAppKeyRequest{
		ProjectId:      d.Get(appKeyProjectIDVar).(string),
		AppId:          d.Get(appKeyAppIDVar).(string),
		Type:           authn.KeyType(authn.KeyType_value[keyType]),
		ExpirationDate: timestamppb.New(t),
	})

	d.SetId(resp.GetId())
	if err := d.Set(appKeyKeyDetailsVar, string(resp.GetKeyDetails())); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readAppKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(appKeyOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(appKeyProjectIDVar).(string)
	appID := d.Get(appKeyAppIDVar).(string)
	resp, err := client.GetAppKey(ctx, &management2.GetAppKeyRequest{
		ProjectId: projectID,
		AppId:     appID,
		KeyId:     d.Id(),
	})
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(resp.GetKey().GetId())

	set := map[string]interface{}{
		appKeyExpirationDateVar: resp.GetKey().GetExpirationDate().AsTime().Format(time.RFC3339),
		appKeyProjectIDVar:      projectID,
		appKeyAppIDVar:          appID,
		appKeyOrgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of app key: %v", k, err)
		}
	}
	return nil
}
