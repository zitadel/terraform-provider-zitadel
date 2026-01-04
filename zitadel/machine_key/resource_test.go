package machine_key_test

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_key"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user/machine_user_test_dep"
)

func TestAccMachineKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_machine_key")
	userDep, userID := machine_user_test_dep.Create(t, frame, frame.UniqueResourcesID)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	resourceExample = removePublicKeyFromResource(resourceExample)
	exampleProperty := test_utils.AttributeValue(t, machine_key.ExpirationDateVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "2051-01-01T00:00:00Z",
		"", "", "",
		false,
		checkRemoteProperty(frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, machine_key.UserIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, machine_key.KeyDetailsVar),
		),
	)
}

func TestAccMachineKeyWithPublicKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_machine_key")
	userDep, userID := machine_user_test_dep.Create(t, frame, frame.UniqueResourcesID)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, machine_key.ExpirationDateVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "2051-01-01T00:00:00Z",
		"", "", "",
		false,
		checkRemoteProperty(frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, machine_key.UserIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportNothing,
			importStateAttributeBase64(frame.BaseTestFrame, machine_key.PublicKeyVar),
		),
	)
}

func removePublicKeyFromResource(resource string) string {
	// Pattern matches: public_key = <<EOT through to EOT including newlines
	re := regexp.MustCompile(`(?s)\s*public_key\s*=\s*<<-EOT.*?EOT`)
	return re.ReplaceAllString(resource, "")
}

func importStateAttributeBase64(frame test_utils.BaseTestFrame, attr string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		primary := frame.State(state)
		val, ok := primary.Attributes[attr]
		if !ok {
			return "", fmt.Errorf("attribute %s not found in attributes %+v", attr, primary.Attributes)
		}
		// Remove the quotes wrapping - just use raw value
		val = strings.ReplaceAll(val, ":", helper.SemicolonPlaceholder)

		valBase64 := base64.StdEncoding.EncodeToString([]byte(val))
		return valBase64, nil
	}
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetMachineKeyByIDs(frame, &management.GetMachineKeyByIDsRequest{
				UserId: userID,
				KeyId:  frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := resp.GetKey().GetExpirationDate().AsTime().Format("2006-01-02T15:04:05Z")
			if expect != actual {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
