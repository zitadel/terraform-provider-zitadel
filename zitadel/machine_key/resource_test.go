package machine_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/machine_key"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/machine_user/machine_user_test_dep"
)

func TestAccMachineKey(t *testing.T) {
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
		checkRemoteProperty(*frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, machine_key.UserIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, machine_key.KeyDetailsVar),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID string) func(string) resource.TestCheckFunc {
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
