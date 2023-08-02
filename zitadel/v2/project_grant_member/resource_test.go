package project_grant_member_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccProjectGrantMember(t *testing.T) {
	resourceName := "zitadel_project_grant_member"
	initialProperty := "PROJECT_GRANT_OWNER"
	updatedProperty := "PROJECT_GRANT_OWNER_VIEWER"
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
	otherOrgFrame, err := frame.AnotherOrg(frame.UniqueResourcesID)
	if err != nil {
		t.Fatalf("failed to switch to another org: %v", err)
	}
	grant, err := frame.AddProjectGrant(frame, &management.AddProjectGrantRequest{
		ProjectId:    projectID,
		GrantedOrgId: otherOrgFrame.OrgID,
	})
	if err != nil {
		t.Fatalf("failed create project grant: %v", err)
	}
	grantID := grant.GetGrantId()
	otherOrgUser, err := otherOrgFrame.ImportHumanUser(otherOrgFrame, &management.ImportHumanUserRequest{
		UserName: otherOrgFrame.UniqueResourcesID,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName: "Don't",
			LastName:  "Care",
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           "dont@care.com",
			IsEmailVerified: true,
		},
	})
	otherOrgUserID := otherOrgUser.GetUserId()
	if err != nil {
		t.Fatalf("failed to create otherOrgUser: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		otherOrgFrame.BaseTestFrame,
		func(cfg, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id     = "%s"
  project_id = "%s"
  grant_id   = "%s"
  user_id    = "%s"
  roles      = ["%s"]
}`, resourceName, otherOrgFrame.UniqueResourcesID, otherOrgFrame.OrgID, projectID, grantID, otherOrgUserID, cfg)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*otherOrgFrame, projectID, grantID, otherOrgUserID),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*otherOrgFrame, projectID, grantID, otherOrgUserID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, projectID, grantID, userID string) func(interface{}) resource.TestCheckFunc {
	return func(expected interface{}) resource.TestCheckFunc {
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
			if expected != actualRoleKeys[0] {
				return fmt.Errorf("expected role key %s, but got %s", expected, actualRoleKeys[0])
			}
			return nil
		}
	}
}
