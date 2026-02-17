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

func TestAccActionTarget(t *testing.T) {
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
			rs, ok := state.RootModule().Resources[frame.TerraformName]
			if !ok {
				return fmt.Errorf("not found: %s", frame.TerraformName)
			}

			client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			remoteResource, err := client.GetTarget(
				context.Background(),
				&actionv2.GetTargetRequest{Id: rs.Primary.ID},
			)
			if err != nil {
				return err
			}
			actualEndpoint := remoteResource.GetTarget().GetEndpoint()
			if expectedEndpoint != "" && actualEndpoint != expectedEndpoint {
				return fmt.Errorf("expected endpoint %q, but got %q", expectedEndpoint, actualEndpoint)
			}
			actualPayloadType := remoteResource.GetTarget().GetPayloadType().String()
			if expectedPayloadType != "" && actualPayloadType != expectedPayloadType {
				return fmt.Errorf("expected payload_type %q, but got %q", expectedPayloadType, actualPayloadType)
			}
			return nil
		}
	}
}

func TestAccActionTargetPayloadTypeJSON(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "payload_type", "PAYLOAD_TYPE_JSON"),
					checkRemoteProperty(frame, "PAYLOAD_TYPE_JSON")(""),
				),
			},
		},
	})
}

func TestAccActionTargetPayloadTypeJWT(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWT"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "payload_type", "PAYLOAD_TYPE_JWT"),
					checkRemoteProperty(frame, "PAYLOAD_TYPE_JWT")(""),
				),
			},
		},
	})
}

func TestAccActionTargetPayloadTypeJWE(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "payload_type", "PAYLOAD_TYPE_JWE"),
					checkRemoteProperty(frame, "PAYLOAD_TYPE_JWE")(""),
				),
			},
		},
	})
}

func TestAccActionTargetPayloadTypeUpdate(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWT"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "payload_type", "PAYLOAD_TYPE_JSON"),
					checkRemoteProperty(frame, "PAYLOAD_TYPE_JSON")(""),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "payload_type", "PAYLOAD_TYPE_JWT"),
					checkRemoteProperty(frame, "PAYLOAD_TYPE_JWT")(""),
				),
			},
		},
	})
}

func TestAccActionTargetTypeRestWebhook(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_WEBHOOK"
  timeout            = "10s"
  interrupt_on_error = true
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "target_type", "REST_WEBHOOK"),
					checkTargetType(frame, "REST_WEBHOOK"),
				),
			},
		},
	})
}

func TestAccActionTargetTypeRestCall(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_CALL"
  timeout            = "10s"
  interrupt_on_error = true
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "target_type", "REST_CALL"),
					checkTargetType(frame, "REST_CALL"),
				),
			},
		},
	})
}

func TestAccActionTargetTypeRestAsync(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "target_type", "REST_ASYNC"),
					checkTargetType(frame, "REST_ASYNC"),
				),
			},
		},
	})
}

func checkTargetType(frame *test_utils.InstanceTestFrame, expectedTargetType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}

		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
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
