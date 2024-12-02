package machine_user_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user"
)

func TestAccMachineUser(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_machine_user")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleUsername := test_utils.AttributeValue(t, machine_user.UserNameVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, exampleUsername, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, machine_user.DescriptionVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedproperty",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			func(state *terraform.State) (string, error) {
				return strconv.FormatBool(test_utils.AttributeValue(t, machine_user.WithSecretVar, exampleAttributes).True()), nil
			},
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetUserByID(frame, &management.GetUserByIDRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetUser().GetMachine().GetDescription()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
