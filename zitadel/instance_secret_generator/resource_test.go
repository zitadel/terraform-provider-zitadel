package instance_secret_generator_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceSecretGenerator(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_secret_generator")
	resourceExample := `
resource "zitadel_instance_secret_generator" "default" {
  generator_type        = "init_code"
  length                = 8
  include_lower_letters = true
  include_upper_letters = true
  include_digits        = true
  include_symbols       = false
}
`
	updatedExample := `
resource "zitadel_instance_secret_generator" "default" {
  generator_type        = "init_code"
  length                = 10
  include_lower_letters = true
  include_upper_letters = true
  include_digits        = true
  include_symbols       = false
}
`
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		func(property bool, secret string) string {
			if property {
				return resourceExample
			}
			return updatedExample
		},
		true, false,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		regexp.MustCompile(`^init_code$`),
		test_utils.CheckNothing,
		func(_ *terraform.State) (string, error) { return "init_code", nil },
	)
}

func TestAccInstanceSecretGeneratorPartialConfig(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_secret_generator")

	// Apply with only generator_type set; all other fields should be
	// adopted from the server and not clobbered to zero values.
	partialConfig := fmt.Sprintf(`%s
resource "zitadel_instance_secret_generator" "partial" {
  generator_type = "otp_sms"
}
`, frame.ProviderSnippet)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:             partialConfig,
				ExpectNonEmptyPlan: true,
				Check: func(state *terraform.State) error {
					client, err := helper.GetAdminClient(context.Background(), frame.ClientInfo)
					if err != nil {
						return fmt.Errorf("failed to get client: %w", err)
					}
					resp, err := client.GetSecretGenerator(context.Background(), &admin.GetSecretGeneratorRequest{
						GeneratorType: settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_OTP_SMS,
					})
					if err != nil {
						return fmt.Errorf("getting secret generator failed: %w", err)
					}
					sg := resp.GetSecretGenerator()
					if sg.GetLength() == 0 {
						return fmt.Errorf("expected length to be preserved from server, but got 0")
					}
					return nil
				},
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetAdminClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			resp, err := client.GetSecretGenerator(context.Background(), &admin.GetSecretGeneratorRequest{
				GeneratorType: settings.SecretGeneratorType_SECRET_GENERATOR_TYPE_INIT_CODE,
			})
			if err != nil {
				return fmt.Errorf("getting secret generator failed: %w", err)
			}
			length := resp.GetSecretGenerator().GetLength()
			if expect && length != 8 {
				return fmt.Errorf("expected length 8, but got %d", length)
			}
			if !expect && length != 10 {
				return fmt.Errorf("expected length 10, but got %d", length)
			}
			return nil
		}
	}
}
