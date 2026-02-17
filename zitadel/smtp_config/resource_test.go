package smtp_config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/smtp_config"
)

func TestAccSMTPConfig(t *testing.T) {
	t.Skip("Skipping test: uses deprecated SMTP API that conflicts with email_provider_smtp tests on same instance")
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_smtp_config")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	senderAddressProperty := test_utils.AttributeValue(t, smtp_config.SenderAddressVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, senderAddressProperty, fmt.Sprintf("zitadel@%s", frame.InstanceDomain), 1)

	exampleProperty := test_utils.AttributeValue(t, smtp_config.SenderNameVar, exampleAttributes).AsString()
	updatedProperty := "updatedProperty"

	exampleSecret := test_utils.AttributeValue(t, smtp_config.PasswordVar, exampleAttributes).AsString()
	updatedSecret := "updatedSecret"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, updatedProperty,
		smtp_config.PasswordVar, exampleSecret, updatedSecret,
		true,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, smtp_config.PasswordVar),
		),
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetSMTPConfigById(frame, &admin.GetSMTPConfigByIdRequest{Id: frame.State(state).ID})
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
