package v2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/action"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"
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
		Description: "Resource representing an action belonging to an organization.",
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
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
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

	_, err = client.UpdateAction(ctx, &management.UpdateActionRequest{
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

	_, err = client.DeleteAction(ctx, &management.DeleteActionRequest{
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

	resp, err := client.CreateAction(ctx, &management.CreateActionRequest{
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

	client, err := getManagementClient(clientinfo, d.Get(actionOrgId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActions(ctx, &management.ListActionsRequest{
		Queries: []*management.ActionQuery{
			{Query: &management.ActionQuery_ActionIdQuery{
				ActionIdQuery: &action.ActionIDQuery{
					Id: d.Id(),
				},
			}},
		},
	})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read action: %v", err)
	}

	if len(resp.Result) == 1 {
		action := resp.Result[0]
		set := map[string]interface{}{
			actionOrgId:         action.GetDetails().GetResourceOwner(),
			actionName:          action.GetName(),
			actionState:         action.GetState(),
			actionScript:        action.GetScript(),
			actionTimeout:       action.GetTimeout().AsDuration().String(),
			actionAllowedToFail: action.GetAllowedToFail(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of action: %v", k, err)
			}
		}
		d.SetId(action.GetId())
		return nil
	}

	d.SetId("")
	return nil
}
