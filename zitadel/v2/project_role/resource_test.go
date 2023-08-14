package project_role_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccProjectRole(t *testing.T) {
	resourceName := "zitadel_project_role"
	initialProperty := "initialProperty"
	updatedProperty := "updatedProperty"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	proj, err := frame.AddProject(frame, &management.AddProjectRequest{
		Name: frame.UniqueResourcesID,
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	projectID := proj.GetId()
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id         = "%s"
  project_id     = "%s"
  role_key     = "%s"
  display_name = "display_name2"
  group        = "role_group"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, projectID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(*frame, projectID),
		regexp.MustCompile(fmt.Sprintf("^%s_%s_(%s|%s)$", helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern, initialProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, projectID), ""),
		nil,
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
