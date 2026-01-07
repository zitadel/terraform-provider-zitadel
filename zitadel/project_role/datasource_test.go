package project_role_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

func TestAccProjectRolesDatasource_All(t *testing.T) {
	datasourceName := "zitadel_project_roles"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	projectDep, projectID := project_test_dep.Create(t, frame, "project_roles_datasource_"+frame.UniqueResourcesID)

	roleKeys := []string{"role1_" + frame.UniqueResourcesID, "role2_" + frame.UniqueResourcesID, "role3_" + frame.UniqueResourcesID}
	for _, key := range roleKeys {
		_, err := frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     key,
			DisplayName: key,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_project_roles" "default" {
  org_id     = "%s"
  project_id = "%s"
}
`, frame.OrgID, projectID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		nil,
		map[string]string{
			"role_keys.#": "3",
		},
	)
}

func TestAccProjectRolesDatasource_FilterByKey(t *testing.T) {
	datasourceName := "zitadel_project_roles"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	projectDep, projectID := project_test_dep.Create(t, frame, "project_roles_datasource_"+frame.UniqueResourcesID)
	matchingKey := "admin_" + frame.UniqueResourcesID
	roleKeys := []string{matchingKey, "user_" + frame.UniqueResourcesID, "viewer_" + frame.UniqueResourcesID}

	for _, key := range roleKeys {
		_, err := frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     key,
			DisplayName: key,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_project_roles" "default" {
  org_id     = "%s"
  project_id = "%s"
  role_key   = "%s"
}
`, frame.OrgID, projectID, matchingKey)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		checkRoleKeyMatch(frame, projectID, matchingKey),
		map[string]string{
			"role_keys.#": "1",
			"role_keys.0": matchingKey,
		},
	)
}

func TestAccProjectRolesDatasource_FilterByDisplayName(t *testing.T) {
	datasourceName := "zitadel_project_roles"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	projectDep, projectID := project_test_dep.Create(t, frame, "project_roles_datasource_"+frame.UniqueResourcesID)
	roleKey := "admin_" + frame.UniqueResourcesID
	displayName := "Administrator_" + frame.UniqueResourcesID

	_, err := frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     roleKey,
		DisplayName: displayName,
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_project_roles" "default" {
  org_id       = "%s"
  project_id   = "%s"
  display_name = "%s"
}
`, frame.OrgID, projectID, displayName)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		nil,
		map[string]string{
			"role_keys.#": "1",
			"role_keys.0": roleKey,
		},
	)
}

func TestAccProjectRolesDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_project_roles"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	projectDep, projectID := project_test_dep.Create(t, frame, "project_roles_datasource_"+frame.UniqueResourcesID)

	roleKeys := []string{"role1_" + frame.UniqueResourcesID, "role2_" + frame.UniqueResourcesID}
	for _, key := range roleKeys {
		_, err := frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     key,
			DisplayName: key,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_project_roles" "default" {
  org_id     = "%s"
  project_id = "%s"
  role_key   = "nonexistent"
}
`, frame.OrgID, projectID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		nil,
		map[string]string{
			"role_keys.#": "0",
		},
	)
}

func checkRoleKeyMatch(frame *test_utils.OrgTestFrame, projectID, expectedKey string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.ListProjectRoles(frame, &management.ListProjectRolesRequest{
			ProjectId: projectID,
		})
		if err != nil {
			return err
		}

		for _, role := range resp.Result {
			if role.Key == expectedKey {
				return nil
			}
		}

		return fmt.Errorf("expected role key %s not found in project %s", expectedKey, projectID)
	}
}
