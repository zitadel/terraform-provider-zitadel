package user_grant_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccUserGrant(t *testing.T) {
	resourceName := "zitadel_user_grant"
	initialProperty := "initialProperty"
	updatedProperty := "updatedProperty"
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
	for _, role := range []string{initialProperty, updatedProperty} {
		_, err = frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     role,
			DisplayName: role,
		})
		if err != nil {
			t.Fatalf("failed to create project role %s: %v", role, err)
		}
	}
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
  org_id         = "%s"
  project_id     = "%s"
  user_id		 = "%s"
  role_keys      = ["%s"]
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, projectID, userID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		true,
		checkRemoteProperty(*frame, userID),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		nil, nil, "", "",
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
