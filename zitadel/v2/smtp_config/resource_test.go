package smtp_config_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/smtp_config"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccSMTPConfig(t *testing.T) {
	resourceName := "zitadel_smtp_config"
	initialProperty := "initialProperty"
	updatedProperty := "updatedProperty"
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
  sender_address = "address"
  sender_name    = "%s"
  tls            = true
  host           = "localhost:25"
  user           = "user"
  password       = "%s"
}`, resourceName, frame.UniqueResourcesID, configProperty, secretProperty)
		},
		initialProperty, updatedProperty,
		initialSecret, updatedSecret,
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckNothing,
		nil, nil, "", smtp_config.PasswordVar,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetSMTPConfig(frame, &admin.GetSMTPConfigRequest{})
			if err != nil {
				return fmt.Errorf("getting smtp config failed: %w", err)
			}
			actual := resp.GetSmtpConfig().GetSenderName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
