package user_grant_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org/org_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_role/project_role_test_dep"
)

func TestAccUserGrantDatasource(t *testing.T) {
	datasourceName := "zitadel_user_grant"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, projectID, roleKey)
	grant, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
		RoleKeys:  []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_grant" "default" {
  org_id   = "%s"
  user_id  = "%s"
  grant_id = "%s"
}
`, frame.OrgID, userID, grant.GetUserGrantId())

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep, userDep},
		nil,
		map[string]string{
			"project_id":  projectID,
			"role_keys.#": "1",
			"role_keys.0": roleKey,
		},
	)
}

func TestAccUserGrantsDatasource_All(t *testing.T) {
	datasourceName := "zitadel_user_grants"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	projectNames := []string{"project1_" + frame.UniqueResourcesID, "project2_" + frame.UniqueResourcesID, "project3_" + frame.UniqueResourcesID}
	for _, projectName := range projectNames {
		_, projectID := project_test_dep.Create(t, frame, projectName)
		project_role_test_dep.Create(t, frame, projectID, roleKey)
		_, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
			UserId:    userID,
			ProjectId: projectID,
			RoleKeys:  []string{roleKey},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_user_grants" "default" {
  org_id  = "%s"
  user_id = "%s"
}
`, frame.OrgID, userID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"ids.#": "3",
		},
	)
}

func TestAccUserGrantsDatasource_FilterByProject(t *testing.T) {
	datasourceName := "zitadel_user_grants"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	_, matchingProjectID := project_test_dep.Create(t, frame, "matching_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, matchingProjectID, roleKey)
	matching, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: matchingProjectID,
		RoleKeys:  []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, otherProjectID := project_test_dep.Create(t, frame, "other_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, otherProjectID, roleKey)
	_, err = frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: otherProjectID,
		RoleKeys:  []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_grants" "default" {
  org_id     = "%s"
  user_id    = "%s"
  project_id = "%s"
}
`, frame.OrgID, userID, matchingProjectID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": matching.GetUserGrantId(),
		},
	)
}

func TestAccUserGrantsDatasource_FilterByProjectGrant(t *testing.T) {
	datasourceName := "zitadel_user_grants"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	roleKey := "role_" + frame.UniqueResourcesID
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, projectID, roleKey)
	_, grantedOrgID, grantedFrame := org_test_dep.Create(t, frame, "granted_org")
	projectGrant, err := frame.AddProjectGrant(frame, &management.AddProjectGrantRequest{
		ProjectId:    projectID,
		GrantedOrgId: grantedOrgID,
		RoleKeys:     []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, userID := human_user_test_dep.Create(t, grantedFrame)
	matching, err := grantedFrame.AddUserGrant(grantedFrame, &management.AddUserGrantRequest{
		UserId:         userID,
		ProjectId:      projectID,
		ProjectGrantId: projectGrant.GetGrantId(),
		RoleKeys:       []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_grants" "default" {
  org_id           = "%s"
  user_id          = "%s"
  project_grant_id = "%s"
}
`, grantedOrgID, userID, projectGrant.GetGrantId())

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": matching.GetUserGrantId(),
		},
	)
}

func TestAccUserGrantsDatasource_FilterByRoleKey(t *testing.T) {
	datasourceName := "zitadel_user_grants"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	userDep, userID := human_user_test_dep.Create(t, frame)

	matchingRoleKey := "admin_" + frame.UniqueResourcesID
	_, matchingProjectID := project_test_dep.Create(t, frame, "matching_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, matchingProjectID, matchingRoleKey)
	matching, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: matchingProjectID,
		RoleKeys:  []string{matchingRoleKey},
	})
	if err != nil {
		t.Fatal(err)
	}
	otherRoleKey := "viewer_" + frame.UniqueResourcesID
	_, otherProjectID := project_test_dep.Create(t, frame, "other_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, otherProjectID, otherRoleKey)
	_, err = frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: otherProjectID,
		RoleKeys:  []string{otherRoleKey},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_grants" "default" {
  org_id   = "%s"
  user_id  = "%s"
  role_key = "%s"
}
`, frame.OrgID, userID, matchingRoleKey)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": matching.GetUserGrantId(),
		},
	)
}

func TestAccUserGrantsDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_user_grants"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, projectID, roleKey)
	_, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
		RoleKeys:  []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_grants" "default" {
  org_id   = "%s"
  user_id  = "%s"
  role_key = "nonexistent"
}
`, frame.OrgID, userID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"ids.#": "0",
		},
	)
}
