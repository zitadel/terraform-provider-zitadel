package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

const (
	actionOrgId         = "org_id"
	actionState         = "state"
	actionName          = "name"
	actionScript        = "script"
	actionTimeout       = "timeout"
	actionAllowedToFail = "allowed_to_fail"
)

func GetAction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			actionOrgId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			actionState: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the state of the action",
			},
			actionName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			actionScript: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			actionTimeout: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "after which time the action will be terminated if not finished",
			},
			actionAllowedToFail: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "when true, the next action will be called even if this action fails",
			},
		},
		CreateContext: createAction,
		DeleteContext: deleteAction,
		ReadContext:   readAction,
		UpdateContext: updateAction,
	}
}

func updateAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(actionOrgId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(actionTimeout).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAction(ctx, &management2.UpdateActionRequest{
		Id:            d.Id(),
		Name:          d.Get(actionName).(string),
		Script:        d.Get(actionScript).(string),
		Timeout:       durationpb.New(timeout),
		AllowedToFail: d.Get(actionAllowedToFail).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update action: %v", err)
	}
	return nil
}

func deleteAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(actionOrgId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteAction(ctx, &management2.DeleteActionRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete action: %v", err)
	}
	return nil
}

func createAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(actionOrgId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(actionTimeout).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateAction(ctx, &management2.CreateActionRequest{
		Name:          d.Get(actionName).(string),
		Script:        d.Get(actionScript).(string),
		Timeout:       durationpb.New(timeout),
		AllowedToFail: d.Get(actionAllowedToFail).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create action: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func readAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(actionOrgId).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActions(ctx, &management2.ListActionsRequest{})
	if err != nil {
		return diag.Errorf("failed to read action: %v", err)
	}

	set := map[string]interface{}{}
	actionIDStr := ""
	for i := range resp.Result {
		action := resp.Result[i]
		if action.GetId() == d.Id() {
			actionIDStr = d.Id()
			set[actionOrgId] = action.GetDetails().GetResourceOwner()
			set[actionName] = action.GetName()
			set[actionState] = action.GetState()
			set[actionScript] = action.GetScript()
			set[actionTimeout] = action.GetTimeout().String()
			set[actionAllowedToFail] = action.GetAllowedToFail()
		}
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of action: %v", k, err)
		}
	}
	d.SetId(actionIDStr)
	return nil
}
