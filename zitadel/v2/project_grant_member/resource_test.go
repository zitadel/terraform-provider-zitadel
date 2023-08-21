package project_grant_member_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org/org_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_grant/project_grant_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_grant_member"
)

func TestAccProjectGrantMember(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_grant_member")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_grant_member.RolesVar, exampleAttributes).AsValueSlice()[0].AsString()
	grantIDProperty := test_utils.AttributeValue(t, project_grant_member.GrantIDVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame)
	userDep, userID := human_user_test_dep.Create(t, frame)
	_, grantedOrgID, _ := org_test_dep.Create(t, frame, "granting_org")
	grantID := project_grant_test_dep.Create(t, frame, projectID, grantedOrgID)
	resourceExample = strings.Replace(resourceExample, grantIDProperty, grantID, 1)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "PROJECT_GRANT_OWNER_VIEWER",
		"", "",
		true,
		checkRemoteProperty(*frame, projectID, grantID, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID, grantID, userID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, projectID, grantID, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.ListProjectGrantMembers(frame, &management.ListProjectGrantMembersRequest{
				ProjectId: projectID,
				GrantId:   grantID,
				Queries: []*member.SearchQuery{{
					Query: &member.SearchQuery_UserIdQuery{
						UserIdQuery: &member.UserIDQuery{
							UserId: userID,
						},
					},
				}},
			})
			if err != nil {
				return err
			}
			if len(resp.Result) != 1 {
				return fmt.Errorf("expected 1 result, but got %d: %w", len(resp.Result), test_utils.ErrNotFound)
			}
			actualRoleKeys := resp.Result[0].GetRoles()
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
