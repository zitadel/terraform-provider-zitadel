package project_member_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_grant_member"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_member"
)

func TestAccProjectMember(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_member")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_grant_member.RolesVar, exampleAttributes).AsValueSlice()[0].AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	userDep, userID := human_user_test_dep.Create(t, frame)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "PROJECT_OWNER_VIEWER",
		"", "", "",
		true,
		checkRemoteProperty(*frame, projectID, userID),
		regexp.MustCompile(fmt.Sprintf("^%s_%s_%s$", helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID, userID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_member.ProjectIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_member.UserIDVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, projectID, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.ListProjectMembers(frame, &management.ListProjectMembersRequest{
				ProjectId: projectID,
				Queries: []*member.SearchQuery{{
					Query: &member.SearchQuery_UserIdQuery{UserIdQuery: &member.UserIDQuery{UserId: userID}},
				}},
			})
			if err != nil {
				return err
			}
			if len(resp.Result) == 0 || len(resp.Result[0].Roles) == 0 {
				return fmt.Errorf("expected 1 user with 1 role, but got %d: %w", len(resp.Result), test_utils.ErrNotFound)
			}
			actual := resp.Result[0].Roles[0]
			if expect != actual {
				return fmt.Errorf("expected role %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
