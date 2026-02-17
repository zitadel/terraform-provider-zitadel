package action_target_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionTargetDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	config := fmt.Sprintf(`
%s
resource "zitadel_action_target" "test" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_WEBHOOK"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}

data "zitadel_action_target" "test" {
  target_id = zitadel_action_target.test.id
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zitadel_action_target.test", "id", "zitadel_action_target.test", "id"),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "name", frame.UniqueResourcesID),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "endpoint", "https://example.com/datasource-test"),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "target_type", "REST_WEBHOOK"),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "timeout", "10s"),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "interrupt_on_error", "false"),
					resource.TestCheckResourceAttr("data.zitadel_action_target.test", "payload_type", "PAYLOAD_TYPE_JSON"),
				),
			},
		},
	})
}

func TestAccActionTargetDatasourcePayloadTypes(t *testing.T) {
	tests := []struct {
		name        string
		payloadType string
	}{
		{"JWT", "PAYLOAD_TYPE_JWT"},
		{"JWE", "PAYLOAD_TYPE_JWE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
			config := fmt.Sprintf(`
%s
resource "zitadel_action_target" "test" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_ASYNC"
  timeout            = "5s"
  interrupt_on_error = false
  payload_type       = "%s"
}

data "zitadel_action_target" "test" {
  target_id = zitadel_action_target.test.id
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, tt.payloadType)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: frame.V6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrPair("data.zitadel_action_target.test", "id", "zitadel_action_target.test", "id"),
							resource.TestCheckResourceAttr("data.zitadel_action_target.test", "payload_type", tt.payloadType),
						),
					},
				},
			})
		})
	}
}

func TestAccActionTargetDatasourceTargetTypes(t *testing.T) {
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
			config := fmt.Sprintf(`
%s
resource "zitadel_action_target" "test" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "%s"
  timeout            = "10s"
  interrupt_on_error = %t
  payload_type       = "PAYLOAD_TYPE_JSON"
}

data "zitadel_action_target" "test" {
  target_id = zitadel_action_target.test.id
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, tt.targetType, tt.interruptOnError)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: frame.V6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrPair("data.zitadel_action_target.test", "id", "zitadel_action_target.test", "id"),
							resource.TestCheckResourceAttr("data.zitadel_action_target.test", "target_type", tt.targetType),
						),
					},
				},
			})
		})
	}
}
