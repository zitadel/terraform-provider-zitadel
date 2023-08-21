package project_role_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_role"
)

func TestAccProjectRole(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_role")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_role.KeyVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedProperty",
		"", "",
		true,
		checkRemoteProperty(*frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, projectID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.ListProjectRoles(frame, &management.ListProjectRolesRequest{
				ProjectId: projectID,
				Queries: []*project.RoleQuery{{
					Query: &project.RoleQuery_KeyQuery{
						KeyQuery: &project.RoleKeyQuery{Key: expect},
					},
				}},
			})
			if err != nil {
				return err
			}
			actualRoles := resp.GetResult()
			if len(actualRoles) == 0 {
				return test_utils.ErrNotFound
			}
			if len(actualRoles) != 1 {
				return fmt.Errorf("expected 1 role, but got %v", actualRoles)
			}
			actualRole := actualRoles[0].GetKey()
			if actualRole != expect {
				return fmt.Errorf("expected role key %s, but got %s", expect, actualRole)
			}
			return nil
		}
	}
}
