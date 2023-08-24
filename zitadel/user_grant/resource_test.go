package user_grant_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project_role/project_role_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/user_grant"
)

func TestAccUserGrant(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grant")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, user_grant.RoleKeysVar, exampleAttributes).AsValueSlice()[0].AsString()
	updatedProperty := "updatedProperty"
	projectDep, projectID := project_test_dep.Create(t, frame)
	project_role_test_dep.Create(t, frame, projectID, exampleProperty, updatedProperty)
	userDep, userID := human_user_test_dep.Create(t, frame)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(*frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, user_grant.UserIDVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetUserGrantByID(frame, &management.GetUserGrantByIDRequest{
				UserId:  userID,
				GrantId: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actualRoleKeys := resp.GetUserGrant().GetRoleKeys()
			if len(actualRoleKeys) != 1 {
				return fmt.Errorf("expected 1 role, but got %d", len(actualRoleKeys))
			}
			if expect != actualRoleKeys[0] {
				return fmt.Errorf("expected role key %s, but got %s", expect, actualRoleKeys[0])
			}
			return nil
		}
	}
}
