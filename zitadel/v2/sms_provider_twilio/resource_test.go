package sms_provider_twilio_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/sms_provider_twilio"
)

func TestAccSMSProviderTwilio(t *testing.T) {
	resourceName := "zitadel_sms_provider_twilio"
	initialProperty := "123456789"
	updatedProperty := "987654321"
	initialSecret := "initialSecret"
	updatedSecret := "updatedSecret"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, secretProperty string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  sid           = "sid"
  sender_number = "%s"
  token         = "%s"
}`, resourceName, frame.UniqueResourcesID, configProperty, secretProperty)
		},
		initialProperty, updatedProperty,
		sms_provider_twilio.TokenVar, initialSecret, updatedSecret,
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		nil,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetSMSProvider(frame, &admin.GetSMSProviderRequest{Id: frame.State(state).ID})
			if err != nil {
				return fmt.Errorf("getting sms provider failed: %w", err)
			}
			actual := resp.GetConfig().GetTwilio().GetSenderNumber()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
