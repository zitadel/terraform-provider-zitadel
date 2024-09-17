package trigger_actions

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	flowType := d.Get(FlowTypeVar).(string)
	flowTypeValues := helper.EnumValueMap(FlowTypes())
	triggerType := d.Get(TriggerTypeVar).(string)
	triggerTypeValues := helper.EnumValueMap(TriggerTypes())
	_, err = client.SetTriggerActions(helper.CtxWithOrgID(ctx, d), &management.SetTriggerActionsRequest{
		FlowType:    strconv.Itoa(int(flowTypeValues[flowType])),
		TriggerType: strconv.Itoa(int(triggerTypeValues[triggerType])),
		ActionIds:   []string{},
	})
	if helper.IgnoreIfNotFoundError(err) != nil {
		return diag.Errorf("failed to delete trigger actions: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	flowType := d.Get(FlowTypeVar).(string)
	flowTypeValues := helper.EnumValueMap(FlowTypes())
	triggerType := d.Get(TriggerTypeVar).(string)
	triggerTypeValues := helper.EnumValueMap(TriggerTypes())
	_, err = client.SetTriggerActions(helper.CtxWithOrgID(ctx, d), &management.SetTriggerActionsRequest{
		FlowType:    strconv.Itoa(int(flowTypeValues[flowType])),
		TriggerType: strconv.Itoa(int(triggerTypeValues[triggerType])),
		ActionIds:   helper.GetOkSetToStringSlice(d, actionsVar),
	})
	if err != nil {
		return diag.Errorf("failed to update trigger actions: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	orgID := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	flowType := d.Get(FlowTypeVar).(string)
	flowTypeValues := helper.EnumValueMap(FlowTypes())
	triggerType := d.Get(TriggerTypeVar).(string)
	triggerTypeValues := helper.EnumValueMap(TriggerTypes())
	actionIDs := helper.GetOkSetToStringSlice(d, actionsVar)
	_, err = client.SetTriggerActions(helper.CtxWithOrgID(ctx, d), &management.SetTriggerActionsRequest{
		FlowType:    strconv.Itoa(int(flowTypeValues[flowType])),
		TriggerType: strconv.Itoa(int(triggerTypeValues[triggerType])),
		ActionIds:   actionIDs,
	})
	if err != nil {
		return diag.Errorf("failed to create trigger actions: %v", err)
	}
	d.SetId(getTriggerActionsID(orgID, flowType, triggerType))
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	orgID := d.Get(helper.OrgIDVar).(string)
	flowType := d.Get(FlowTypeVar).(string)
	triggerType := d.Get(TriggerTypeVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	flowTypeValues := helper.EnumValueMap(FlowTypes())
	triggerTypeNames := TriggerTypes()
	resp, err := client.GetFlow(helper.CtxWithOrgID(ctx, d), &management.GetFlowRequest{Type: strconv.Itoa(int(flowTypeValues[flowType]))})
	if err != nil {
		return diag.FromErr(err)
	}
	var actionIDs []string
	for _, triggerAction := range resp.GetFlow().GetTriggerActions() {
		triggerTypeID, err := strconv.Atoi(triggerAction.GetTriggerType().GetId())
		if err != nil {
			return diag.FromErr(err)
		}
		if triggerTypeNames[int32(triggerTypeID)] != triggerType {
			continue
		}
		for _, action := range triggerAction.GetActions() {
			actionIDs = append(actionIDs, action.GetId())
		}
	}
	if err = d.Set(actionsVar, actionIDs); err != nil {
		return diag.Errorf("setting action ids %s to property %s failed: %v", actionIDs, actionsVar, err)
	}
	d.SetId(getTriggerActionsID(orgID, flowType, triggerType))
	return nil
}

func getTriggerActionsID(orgID, flowType string, triggerType string) string {
	return orgID + "_" + flowType + "_" + triggerType
}
