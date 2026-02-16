package action_target

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
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

	_, err = client.DeleteTarget(ctx, &actionv2.DeleteTargetRequest{
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

	req := &actionv2.UpdateTargetRequest{
		Id: d.Id(),
	}

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

	if d.HasChange(TargetTypeVar) || d.HasChange(InterruptOnErrorVar) {
		targetType := d.Get(TargetTypeVar).(string)
		interruptOnError := d.Get(InterruptOnErrorVar).(bool)

		switch targetType {
		case targetTypeRestWebhook:
			req.TargetType = &actionv2.UpdateTargetRequest_RestWebhook{
				RestWebhook: &actionv2.RESTWebhook{InterruptOnError: interruptOnError},
			}
		case targetTypeRestCall:
			req.TargetType = &actionv2.UpdateTargetRequest_RestCall{
				RestCall: &actionv2.RESTCall{InterruptOnError: interruptOnError},
			}
		case targetTypeRestAsync:
			req.TargetType = &actionv2.UpdateTargetRequest_RestAsync{
				RestAsync: &actionv2.RESTAsync{},
			}
		default:
			return diag.Errorf("unknown target type %s", targetType)
		}
	}

	if d.HasChange(PayloadTypeVar) {
		payloadType := d.Get(PayloadTypeVar).(string)
		req.PayloadType = stringToPayloadType(payloadType)
	}

	_, err = client.UpdateTarget(ctx, req)
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

	req := &actionv2.CreateTargetRequest{
		Name:        d.Get(NameVar).(string),
		Endpoint:    d.Get(EndpointVar).(string),
		Timeout:     durationpb.New(timeout),
		PayloadType: stringToPayloadType(d.Get(PayloadTypeVar).(string)),
	}

	targetType := d.Get(TargetTypeVar).(string)
	interruptOnError := d.Get(InterruptOnErrorVar).(bool)
	switch targetType {
	case targetTypeRestWebhook:
		req.TargetType = &actionv2.CreateTargetRequest_RestWebhook{
			RestWebhook: &actionv2.RESTWebhook{InterruptOnError: interruptOnError},
		}
	case targetTypeRestCall:
		req.TargetType = &actionv2.CreateTargetRequest_RestCall{
			RestCall: &actionv2.RESTCall{InterruptOnError: interruptOnError},
		}
	case targetTypeRestAsync:
		req.TargetType = &actionv2.CreateTargetRequest_RestAsync{
			RestAsync: &actionv2.RESTAsync{},
		}
	default:
		return diag.Errorf("unknown target type %s", targetType)
	}

	resp, err := client.CreateTarget(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create target: %v", err)
	}
	d.SetId(resp.GetId())

	if err := d.Set(SigningKeyVar, resp.GetSigningKey()); err != nil {
		return diag.Errorf("failed to set signing_key: %v", err)
	}

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

	resp, err := client.GetTarget(ctx, &actionv2.GetTargetRequest{
		Id: helper.GetID(d, TargetIDVar),
	})

	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get target: %v", err)
	}

	target := resp.GetTarget()
	if target != nil {
		set := map[string]interface{}{
			NameVar:        target.GetName(),
			EndpointVar:    target.GetEndpoint(),
			TimeoutVar:     target.GetTimeout().AsDuration().String(),
			PayloadTypeVar: payloadTypeToString(target.GetPayloadType()),
		}

		if target.GetRestWebhook() != nil {
			set[TargetTypeVar] = targetTypeRestWebhook
			set[InterruptOnErrorVar] = target.GetRestWebhook().GetInterruptOnError()
		} else if target.GetRestCall() != nil {
			set[TargetTypeVar] = targetTypeRestCall
			set[InterruptOnErrorVar] = target.GetRestCall().GetInterruptOnError()
		} else if target.GetRestAsync() != nil {
			set[TargetTypeVar] = targetTypeRestAsync
			set[InterruptOnErrorVar] = false
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

func stringToPayloadType(s string) actionv2.PayloadType {
	switch s {
	case payloadTypeJSON:
		return actionv2.PayloadType_PAYLOAD_TYPE_JSON
	case payloadTypeJWT:
		return actionv2.PayloadType_PAYLOAD_TYPE_JWT
	case payloadTypeJWE:
		return actionv2.PayloadType_PAYLOAD_TYPE_JWE
	default:
		return actionv2.PayloadType_PAYLOAD_TYPE_JSON
	}
}

func payloadTypeToString(pt actionv2.PayloadType) string {
	switch pt {
	case actionv2.PayloadType_PAYLOAD_TYPE_JSON:
		return payloadTypeJSON
	case actionv2.PayloadType_PAYLOAD_TYPE_JWT:
		return payloadTypeJWT
	case actionv2.PayloadType_PAYLOAD_TYPE_JWE:
		return payloadTypeJWE
	case actionv2.PayloadType_PAYLOAD_TYPE_UNSPECIFIED:
		// UNSPECIFIED means the target was created before payload_type was added
		// or it's using the API default, which is JSON
		return payloadTypeJSON
	default:
		// Unknown payload type, fall back to JSON as the safe default
		return payloadTypeJSON
	}
}
