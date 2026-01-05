package project_grant_member_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org/org_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_grant/project_grant_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_grant_member"
)

func TestAccProjectGrantMember(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_grant_member")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_grant_member.RolesVar, exampleAttributes).AsValueSlice()[0].AsString()
	grantIDProperty := test_utils.AttributeValue(t, project_grant_member.GrantIDVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
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
		"", "", "",
		true,
		checkRemoteProperty(frame, projectID, grantID, userID),
		regexp.MustCompile(fmt.Sprintf(
			"^%s_%s_%s_%s$",
			helper.ZitadelGeneratedIdPattern,
			helper.ZitadelGeneratedIdPattern,
			helper.ZitadelGeneratedIdPattern,
			helper.ZitadelGeneratedIdPattern,
		)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID, grantID, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_grant_member.ProjectIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_grant_member.GrantIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_grant_member.UserIDVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectID, grantID, userID string) func(string) resource.TestCheckFunc {
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
