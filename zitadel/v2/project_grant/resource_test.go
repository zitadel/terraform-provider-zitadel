package project_grant_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org/org_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_grant"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_role/project_role_test_dep"
)

func TestAccProjectGrant(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_grant")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_grant.RoleKeysVar, exampleAttributes).AsValueSlice()[0].AsString()
	updatedProperty := "updatedproperty"
	projectDep, projectID := project_test_dep.Create(t, frame)
	project_role_test_dep.Create(t, frame, projectID, exampleProperty, updatedProperty)
	grantedOrgDep, _, _ := org_test_dep.Create(t, frame, "granted_org")
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep, grantedOrgDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "",
		false,
		checkRemoteProperty(*frame, projectID),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, projectID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetProjectGrantByID(frame, &management.GetProjectGrantByIDRequest{
				ProjectId: projectID,
				GrantId:   frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actualRoleKeys := resp.GetProjectGrant().GetGrantedRoleKeys()
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
