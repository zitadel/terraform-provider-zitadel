package sms_provider_http_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/sms_provider_http"
)

func TestAccSMSHttpProvider(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_http")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, sms_provider_http.EndPointVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "https://relay.example.com/test",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
		),
	)
}

func TestAccSMSHttpProviderActivation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_http")

	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_http" "default" {
  endpoint   = "https://relay.example.com/sms"
  set_active = false
}
`, frame.ProviderSnippet)

	activatedConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_http" "default" {
  endpoint   = "https://relay.example.com/sms"
  set_active = true
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "false"),
				),
			},
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetSMSProvider(frame, &admin.GetSMSProviderRequest{Id: frame.State(state).ID})
			if err != nil {
				return fmt.Errorf("getting sms provider failed: %w", err)
			}
			actual := resp.GetConfig().GetHttp().GetEndpoint()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
