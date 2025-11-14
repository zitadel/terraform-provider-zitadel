package action_execution

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

type BuildConditionFunc func(d *schema.ResourceData) (*action.Condition, error)
type IdFromConditionFunc func(condition *action.Condition) (*string, error)

// NewSetExecution returns a Terraform handler for create and update that
// builds an action.Condition from the resource data, calls the Zitadel
// SetExecution API with the configured targets, and sets the resource ID
// based on the provided idFromCondition function.
func NewSetExecution(buildCondition BuildConditionFunc, idFromCondition IdFromConditionFunc) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		condition, err := buildCondition(d)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "started set (create/update)")

		clientinfo, ok := m.(*helper.ClientInfo)
		if !ok {
			return diag.Errorf("failed to get client")
		}

		client, err := helper.GetActionClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		var targetIDs []string
		if ids, ok := d.GetOk(TargetIDsVar); ok {
			for _, id := range ids.([]interface{}) {
				targetIDs = append(targetIDs, id.(string))
			}
		}

		_, err = client.SetExecution(ctx, &action.SetExecutionRequest{
			Condition: condition,
			Targets:   targetIDs,
		})
		if err != nil {
			return diag.Errorf("failed to set execution: %v", err)
		}

		idPtr, err := idFromCondition(condition)
		if err != nil {
			return diag.FromErr(err)
		}
		if idPtr == nil {
			return diag.Errorf("failed to derive execution ID from condition")
		}

		d.SetId(*idPtr)
		return nil
	}
}

func NewDeleteExecution(buildCondition BuildConditionFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		condition, err := buildCondition(d)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "started delete")

		clientinfo, ok := m.(*helper.ClientInfo)
		if !ok {
			return diag.Errorf("failed to get client")
		}

		client, err := helper.GetActionClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.SetExecution(ctx, &action.SetExecutionRequest{
			Condition: condition,
			Targets:   []string{},
		})
		if err != nil {
			return diag.Errorf("failed to delete execution: %v", err)
		}
		return nil
	}
}
