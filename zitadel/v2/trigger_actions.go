package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/action"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	triggerActionsOrgIDVar       = "org_id"
	triggerActionsFlowTypeVar    = "flow_type"
	triggerActionsTriggerTypeVar = "trigger_type"
	triggerActionsActionsVar     = "action_ids"
)

func GetTriggerActions() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing triggers, when actions get started",
		Schema: map[string]*schema.Schema{
			triggerActionsOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			triggerActionsFlowTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the flow to which the action triggers belong",
				ForceNew:    true,
			},
			triggerActionsTriggerTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trigger type on when the actions get triggered",
				ForceNew:    true,
			},
			triggerActionsActionsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "IDs of the triggered actions",
			},
		},
		DeleteContext: deleteTriggerActions,
		CreateContext: createTriggerActions,
		UpdateContext: updateTriggerActions,
		ReadContext:   readTriggerActions,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func deleteTriggerActions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(triggerActionsOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.SetTriggerActions(ctx, &management2.SetTriggerActionsRequest{
		FlowType:    action.FlowType(action.FlowType_value[d.Get(triggerActionsFlowTypeVar).(string)]),
		TriggerType: action.TriggerType(action.TriggerType_value[d.Get(triggerActionsTriggerTypeVar).(string)]),
		ActionIds:   []string{},
	})
	if err != nil {
		return diag.Errorf("failed to delete trigger actions: %v", err)
	}
	return nil
}

func updateTriggerActions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(triggerActionsOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	actionsSet := d.Get(triggerActionsActionsVar).(*schema.Set)
	actions := make([]string, 0)
	for _, action := range actionsSet.List() {
		actions = append(actions, action.(string))
	}
	_, err = client.SetTriggerActions(ctx, &management2.SetTriggerActionsRequest{
		FlowType:    action.FlowType(action.FlowType_value[d.Get(triggerActionsFlowTypeVar).(string)]),
		TriggerType: action.TriggerType(action.TriggerType_value[d.Get(triggerActionsTriggerTypeVar).(string)]),
		ActionIds:   actions,
	})
	if err != nil {
		return diag.Errorf("failed to update trigger actions: %v", err)
	}

	return nil
}

func createTriggerActions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(triggerActionsOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	actionsSet := d.Get(triggerActionsActionsVar).(*schema.Set)
	actions := make([]string, 0)
	for _, action := range actionsSet.List() {
		actions = append(actions, action.(string))
	}
	flowType := d.Get(triggerActionsFlowTypeVar).(string)
	triggerType := d.Get(triggerActionsTriggerTypeVar).(string)
	_, err = client.SetTriggerActions(ctx, &management2.SetTriggerActionsRequest{
		FlowType:    action.FlowType(action.FlowType_value[flowType]),
		TriggerType: action.TriggerType(action.TriggerType_value[triggerType]),
		ActionIds:   actions,
	})
	d.SetId(getTriggerActionsID(orgID, flowType, triggerType))

	return nil
}

func readTriggerActions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	orgID := d.Get(triggerActionsOrgIDVar).(string)
	flowType := d.Get(triggerActionsFlowTypeVar).(string)
	triggerType := d.Get(triggerActionsTriggerTypeVar).(string)
	d.SetId(getTriggerActionsID(orgID, flowType, triggerType))
	return nil
}

func getTriggerActionsID(orgID, flowType string, triggerType string) string {
	return orgID + "_" + flowType + "_" + triggerType
}
