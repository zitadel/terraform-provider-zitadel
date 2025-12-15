package action

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

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

	timeout, err := time.ParseDuration(d.Get(timeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAction(helper.CtxWithOrgID(ctx, d), &management.UpdateActionRequest{
		Id:            d.Id(),
		Name:          d.Get(NameVar).(string),
		Script:        d.Get(ScriptVar).(string),
		Timeout:       durationpb.New(timeout),
		AllowedToFail: d.Get(allowedToFailVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update action: %v", err)
	}
	return nil
}

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

	_, err = client.DeleteAction(helper.CtxWithOrgID(ctx, d), &management.DeleteActionRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete action: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(timeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateAction(helper.CtxWithOrgID(ctx, d), &management.CreateActionRequest{
		Name:          d.Get(NameVar).(string),
		Script:        d.Get(ScriptVar).(string),
		Timeout:       durationpb.New(timeout),
		AllowedToFail: d.Get(allowedToFailVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create action: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActions(helper.CtxWithOrgID(ctx, d), &management.ListActionsRequest{
		Queries: []*management.ActionQuery{
			{Query: &management.ActionQuery_ActionIdQuery{
				ActionIdQuery: &action.ActionIDQuery{
					Id: helper.GetID(d, ActionIDVar),
				},
			}},
		},
	})
	if err != nil {
		return diag.Errorf("failed to list actions")
	}

	if len(resp.Result) == 1 {
		action := resp.Result[0]
		set := map[string]interface{}{
			helper.OrgIDVar:  action.GetDetails().GetResourceOwner(),
			NameVar:          action.GetName(),
			stateVar:         action.GetState(),
			ScriptVar:        action.GetScript(),
			timeoutVar:       action.GetTimeout().AsDuration().String(),
			allowedToFailVar: action.GetAllowedToFail(),
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
