package action

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/action"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(timeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAction(ctx, &management.UpdateActionRequest{
		Id:            d.Id(),
		Name:          d.Get(nameVar).(string),
		Script:        d.Get(scriptVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
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

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(timeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateAction(ctx, &management.CreateActionRequest{
		Name:          d.Get(nameVar).(string),
		Script:        d.Get(scriptVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActions(ctx, &management.ListActionsRequest{
		Queries: []*management.ActionQuery{
			{Query: &management.ActionQuery_ActionIdQuery{
				ActionIdQuery: &action.ActionIDQuery{
					Id: helper.GetID(d, actionIDVar),
				},
			}},
		},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to list actions")
	}

	if len(resp.Result) == 1 {
		action := resp.Result[0]
		set := map[string]interface{}{
			orgIDVar:         action.GetDetails().GetResourceOwner(),
			nameVar:          action.GetName(),
			stateVar:         action.GetState(),
			scriptVar:        action.GetScript(),
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
