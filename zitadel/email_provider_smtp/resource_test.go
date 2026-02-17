package email_provider_smtp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_smtp"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccEmailSMTPProvider(t *testing.T) {
	t.Skip("Skipping: flaky due to email provider state conflicts with email_provider_http on same instance")
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	senderAddressProperty := test_utils.AttributeValue(t, email_provider_smtp.SenderAddressVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, senderAddressProperty, fmt.Sprintf("zitadel@%s", frame.InstanceDomain), 1)

	exampleProperty := test_utils.AttributeValue(t, email_provider_smtp.SenderNameVar, exampleAttributes).AsString()
	updatedProperty := "updatedProperty"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
		),
		"password", "set_active",
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetEmailProviderById(frame, &admin.GetEmailProviderByIdRequest{Id: frame.State(state).ID})
			if err != nil {
				return fmt.Errorf("getting email provider failed: %w", err)
			}
			actual := resp.GetConfig().GetSmtp().GetSenderName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
