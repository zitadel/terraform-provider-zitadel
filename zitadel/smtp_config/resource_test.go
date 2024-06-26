package smtp_config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/smtp_config"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
)

func TestAccSMTPConfig(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_smtp_config")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	senderAddressProperty := test_utils.AttributeValue(t, smtp_config.SenderAddressVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, senderAddressProperty, fmt.Sprintf("zitadel@%s", frame.InstanceDomain), 1)
	exampleProperty := test_utils.AttributeValue(t, smtp_config.SenderNameVar, exampleAttributes).AsString()
	exampleSecret := test_utils.AttributeValue(t, smtp_config.PasswordVar, exampleAttributes).AsString()
	updatedProperty := "updatedProperty"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, updatedProperty,
		smtp_config.PasswordVar, exampleSecret, "updatedSecret",
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame), updatedProperty),
		test_utils.ImportResourceId(frame.BaseTestFrame),
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
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
