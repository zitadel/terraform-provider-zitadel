package sms_provider_twilio_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/sms_provider_twilio"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
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
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(configProperty, secretProperty interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  sid           = "sid"
  sender_number = "%s"
  token         = "%s"
}`, resourceName, frame.UniqueResourcesID, configProperty, secretProperty)
		},
		initialProperty, updatedProperty,
		initialSecret, updatedSecret,
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckNothing,
		nil, nil, "", sms_provider_twilio.TokenVar,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
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
