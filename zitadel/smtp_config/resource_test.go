package smtp_config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/smtp_config"
)

func TestAccSMTPConfig(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_smtp_config")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	senderAddressProperty := test_utils.AttributeValue(t, smtp_config.SenderAddressVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, senderAddressProperty, fmt.Sprintf("zitadel@%s", frame.InstanceDomain), 1)
	exampleProperty := test_utils.AttributeValue(t, smtp_config.SenderNameVar, exampleAttributes).AsString()
	exampleSecret := test_utils.AttributeValue(t, smtp_config.PasswordVar, exampleAttributes).AsString()
	importParts := []resource.ImportStateIdFunc{
		test_utils.ImportResourceId(frame.BaseTestFrame),
		test_utils.ImportStateAttribute(frame.BaseTestFrame, smtp_config.PasswordVar),
	}

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, "updatedProperty",
		smtp_config.PasswordVar, exampleSecret, "updatedSecret",
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		CheckDestroy(*frame),
		test_utils.ChainImportStateIdFuncs(importParts...),
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

func CheckDestroy(frame test_utils.InstanceTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := checkRemoteProperty(frame)("something")(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}
