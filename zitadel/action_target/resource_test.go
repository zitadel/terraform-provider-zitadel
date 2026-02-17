package action_target_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccTarget(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)

	nameAttribute := test_utils.AttributeValue(t, action_target.NameVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, nameAttribute, frame.UniqueResourcesID, 1)

	exampleProperty := test_utils.AttributeValue(t, action_target.EndpointVar, exampleAttributes).AsString()
	examplePayloadType := test_utils.AttributeValue(t, action_target.PayloadTypeVar, exampleAttributes).AsString()
	updatedProperty := exampleProperty + "-updated"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty,
		updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame, examplePayloadType),
		test_utils.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, examplePayloadType), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
		),
		action_target.SigningKeyVar,
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, expectedPayloadType string) func(string) resource.TestCheckFunc {
	return func(expectedEndpoint string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			remoteResource, err := client.GetTarget(
				context.Background(),
				&actionv2.GetTargetRequest{Id: frame.State(state).ID},
			)
			if err != nil {
				return err
			}
			actualEndpoint := remoteResource.GetTarget().GetEndpoint()
			if actualEndpoint != expectedEndpoint {
				return fmt.Errorf("expected endpoint %q, but got %q", expectedEndpoint, actualEndpoint)
			}
			actualPayloadType := remoteResource.GetTarget().GetPayloadType().String()
			if actualPayloadType != expectedPayloadType {
				return fmt.Errorf("expected payload_type %q, but got %q", expectedPayloadType, actualPayloadType)
			}
			return nil
		}
	}
}

func TestAccTargetPayloadTypes(t *testing.T) {
	tests := []struct {
		name        string
		payloadType string
	}{
		{"JSON", "PAYLOAD_TYPE_JSON"},
		{"JWT", "PAYLOAD_TYPE_JWT"},
		{"JWE", "PAYLOAD_TYPE_JWE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
			resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "test" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "%s"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, tt.payloadType)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: frame.V6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: resourceConfig,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("zitadel_action_target.test", "payload_type", tt.payloadType),
							checkPayloadType(frame, tt.payloadType),
						),
					},
				},
			})
		})
	}
}

func checkPayloadType(frame *test_utils.InstanceTestFrame, expectedPayloadType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}
		rs, ok := state.RootModule().Resources["zitadel_action_target.test"]
		if !ok {
			return fmt.Errorf("resource not found")
		}
		remoteResource, err := client.GetTarget(
			context.Background(),
			&actionv2.GetTargetRequest{Id: rs.Primary.ID},
		)
		if err != nil {
			return err
		}
		actualPayloadType := remoteResource.GetTarget().GetPayloadType().String()
		if actualPayloadType != expectedPayloadType {
			return fmt.Errorf("expected payload_type %q, but got %q", expectedPayloadType, actualPayloadType)
		}
		return nil
	}
}

func TestAccTargetTypes(t *testing.T) {
	tests := []struct {
		name             string
		targetType       string
		interruptOnError bool
	}{
		{"REST_WEBHOOK", "REST_WEBHOOK", true},
		{"REST_CALL", "REST_CALL", true},
		{"REST_ASYNC", "REST_ASYNC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
			resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "test" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "%s"
  timeout            = "10s"
  interrupt_on_error = %t
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, tt.targetType, tt.interruptOnError)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: frame.V6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: resourceConfig,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("zitadel_action_target.test", "target_type", tt.targetType),
							checkTargetType(frame, tt.targetType),
						),
					},
				},
			})
		})
	}
}

func checkTargetType(frame *test_utils.InstanceTestFrame, expectedTargetType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}
		rs, ok := state.RootModule().Resources["zitadel_action_target.test"]
		if !ok {
			return fmt.Errorf("resource not found")
		}
		remoteResource, err := client.GetTarget(
			context.Background(),
			&actionv2.GetTargetRequest{Id: rs.Primary.ID},
		)
		if err != nil {
			return err
		}
		target := remoteResource.GetTarget()
		var actualTargetType string
		if target.GetRestWebhook() != nil {
			actualTargetType = "REST_WEBHOOK"
		} else if target.GetRestCall() != nil {
			actualTargetType = "REST_CALL"
		} else if target.GetRestAsync() != nil {
			actualTargetType = "REST_ASYNC"
		}
		if actualTargetType != expectedTargetType {
			return fmt.Errorf("expected target_type %q, but got %q", expectedTargetType, actualTargetType)
		}
		return nil
	}
}
