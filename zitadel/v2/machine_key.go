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
	machineKeyOrgIDVar          = "org_id"
	machineKeyUserIDVar         = "user_id"
	machineKeyKeyTypeVar        = "key_type"
	machineKeyKeyDetailsVar     = "key_details"
	machineKeyExpirationDateVar = "expiration_date"
)

func GetMachineKey() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a machine key",
		Schema: map[string]*schema.Schema{
			machineKeyOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			machineKeyUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			machineKeyKeyTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the machine key",
				ForceNew:    true,
			},
			machineKeyExpirationDateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration date of the machine key",
				ForceNew:    true,
			},
			machineKeyKeyDetailsVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the machine key",
				Sensitive:   true,
			},
		},
		DeleteContext: deleteMachineKey,
		CreateContext: createMachineKey,
		ReadContext:   readMachineKey,
	}
}

func deleteMachineKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(machineKeyOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveMachineKey(ctx, &management2.RemoveMachineKeyRequest{
		UserId: d.Get(machineKeyUserIDVar).(string),
		KeyId:  d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete machine key: %v", err)
	}
	return nil
}

func createMachineKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(machineKeyOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	t, err := time.Parse(time.RFC3339, d.Get(machineKeyExpirationDateVar).(string))
	if err != nil {
		return diag.Errorf("failed to parse time: %v", err)
	}

	keyType := d.Get(machineKeyKeyTypeVar).(string)
	resp, err := client.AddMachineKey(ctx, &management2.AddMachineKeyRequest{
		UserId:         d.Get(machineKeyUserIDVar).(string),
		Type:           authn.KeyType(authn.KeyType_value[keyType]),
		ExpirationDate: timestamppb.New(t),
	})
	d.SetId(resp.GetKeyId())

	if err := d.Set(machineKeyKeyDetailsVar, string(resp.GetKeyDetails())); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readMachineKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(machineKeyOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(machineKeyUserIDVar).(string)
	resp, err := client.GetMachineKeyByIDs(ctx, &management2.GetMachineKeyByIDsRequest{
		UserId: userID,
		KeyId:  d.Id(),
	})
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(resp.GetKey().GetId())

	set := map[string]interface{}{
		machineKeyExpirationDateVar: resp.GetKey().GetExpirationDate().AsTime().Format(time.RFC3339),
		machineKeyUserIDVar:         userID,
		machineKeyOrgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of machine key: %v", k, err)
		}
	}
	return nil
}
