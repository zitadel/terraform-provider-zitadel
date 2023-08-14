package smtp_config_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/smtp_config"
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
	_, err = frame.RemoveSMTPConfig(frame, &admin.RemoveSMTPConfigRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Fatalf("failed to remove smtp config: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, secretProperty string) string {
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
		smtp_config.PasswordVar, initialSecret, updatedSecret,
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
