package target

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2beta"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteTarget(helper.CtxWithOrgID(ctx, d), &action.DeleteTargetRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete target: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	// Start with the request struct containing only the mandatory ID
	req := &action.UpdateTargetRequest{
		Id: d.Id(),
	}

	// Conditionally set fields only if they have changed in the configuration
	if d.HasChange(NameVar) {
		name := d.Get(NameVar).(string)
		req.Name = &name
	}

	if d.HasChange(EndpointVar) {
		endpoint := d.Get(EndpointVar).(string)
		req.Endpoint = &endpoint
	}

	if d.HasChange(TimeoutVar) {
		timeout, err := time.ParseDuration(d.Get(TimeoutVar).(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.Timeout = durationpb.New(timeout)
	}

	// The target type and interrupt on error must be checked together
	if d.HasChange(TargetTypeVar) || d.HasChange(InterruptOnErrorVar) {
		targetType := d.Get(TargetTypeVar).(string)
		interruptOnError := d.Get(InterruptOnErrorVar).(bool)

		switch targetType {
		case "REST_WEBHOOK":
			req.TargetType = &action.UpdateTargetRequest_RestWebhook{
				RestWebhook: &action.RESTWebhook{InterruptOnError: interruptOnError},
			}
		case "REST_CALL":
			req.TargetType = &action.UpdateTargetRequest_RestCall{
				RestCall: &action.RESTCall{InterruptOnError: interruptOnError},
			}
		case "REST_ASYNC":
			req.TargetType = &action.UpdateTargetRequest_RestAsync{
				RestAsync: &action.RESTAsync{},
			}
		default:
			return diag.Errorf("unknown target type %s", targetType)
		}
	}

	_, err = client.UpdateTarget(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to update target: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	timeout, err := time.ParseDuration(d.Get(TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &action.CreateTargetRequest{
		Name:     d.Get(NameVar).(string),
		Endpoint: d.Get(EndpointVar).(string),
		Timeout:  durationpb.New(timeout),
	}

	targetType := d.Get(TargetTypeVar).(string)
	interruptOnError := d.Get(InterruptOnErrorVar).(bool)
	switch targetType {
	case "REST_WEBHOOK":
		req.TargetType = &action.CreateTargetRequest_RestWebhook{
			RestWebhook: &action.RESTWebhook{InterruptOnError: interruptOnError},
		}
	case "REST_CALL":
		req.TargetType = &action.CreateTargetRequest_RestCall{
			RestCall: &action.RESTCall{InterruptOnError: interruptOnError},
		}
	case "REST_ASYNC":
		req.TargetType = &action.CreateTargetRequest_RestAsync{
			RestAsync: &action.RESTAsync{},
		}
	default:
		return diag.Errorf("unknown target type %s", targetType)
	}

	resp, err := client.CreateTarget(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create target: %v", err)
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

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use GetTarget with the correct field name 'Id'
	resp, err := client.GetTarget(helper.CtxWithOrgID(ctx, d), &action.GetTargetRequest{
		Id: helper.GetID(d, TargetIDVar),
	})

	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get target")
	}

	target := resp.GetTarget()
	if target != nil {
		// Do not set OrgIDVar, as the v2beta API does not return it in the response body.
		// The value from the configuration will be preserved in the state.
		set := map[string]interface{}{
			NameVar:     target.GetName(),
			EndpointVar: target.GetEndpoint(),
			TimeoutVar:  target.GetTimeout().AsDuration().String(),
		}

		if target.GetRestWebhook() != nil {
			set[TargetTypeVar] = "REST_WEBHOOK"
			set[InterruptOnErrorVar] = target.GetRestWebhook().GetInterruptOnError()
		} else if target.GetRestCall() != nil {
			set[TargetTypeVar] = "REST_CALL"
			set[InterruptOnErrorVar] = target.GetRestCall().GetInterruptOnError()
		} else if target.GetRestAsync() != nil {
			set[TargetTypeVar] = "REST_ASYNC"
			set[InterruptOnErrorVar] = false // Not applicable, so set to default
		}

		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of target: %v", k, err)
			}
		}
		d.SetId(target.GetId())
		return nil
	}

	d.SetId("")
	return nil
}
