package project_member_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_member"
)

func TestAccProjectMember(t *testing.T) {
	resourceName := "zitadel_project_member"
	initialProperty := "PROJECT_OWNER"
	updatedProperty := "PROJECT_OWNER_VIEWER"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	project, err := frame.AddProject(frame, &management.AddProjectRequest{
		Name: frame.UniqueResourcesID,
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	projectID := project.GetId()
	user, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
		UserName: frame.UniqueResourcesID,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName: "Don't",
			LastName:  "Care",
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           "dont@care.com",
			IsEmailVerified: true,
		},
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	userID := user.GetUserId()
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	org_id              = "%s"
	project_id          = "%s"
	user_id 			= "%s"
  	roles  				= ["%s"]
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, projectID, userID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(*frame, projectID, userID),
		regexp.MustCompile(fmt.Sprintf("^%s_%s_%s$", helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID, userID), ""),
		test_utils.ConcatImportStateIdFuncs(
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
