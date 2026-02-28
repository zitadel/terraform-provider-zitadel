package project_role_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_role"
)

func TestAccProjectRole(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_role")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_role.KeyVar, exampleAttributes).AsString()
	updatedProperty := "updatedProperty"
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame, projectID),
		regexp.MustCompile(fmt.Sprintf("^%s_%s_(%s|%s)$", helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern, exampleProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_role.ProjectIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, project_role.KeyVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func TestAccProjectRoleDisplayNameUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_role")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_project_role" "default" {
  org_id       = data.zitadel_org.default.id
  project_id   = "%s"
  role_key     = "%s"
  display_name = "Initial Role Name"
  group        = "initial_group"
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_project_role" "default" {
  org_id       = data.zitadel_org.default.id
  project_id   = "%s"
  role_key     = "%s"
  display_name = "Updated Role Name"
  group        = "updated_group"
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "display_name", "Initial Role Name"),
					resource.TestCheckResourceAttr(frame.TerraformName, "group", "initial_group"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "display_name", "Updated Role Name"),
					resource.TestCheckResourceAttr(frame.TerraformName, "group", "updated_group"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectID string) func(string) resource.TestCheckFunc {
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
