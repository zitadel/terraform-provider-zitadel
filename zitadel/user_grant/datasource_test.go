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

func TestAccUserGrantsDatasource_All(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grants")
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	projectIDs := make([]string, 3)
	for i := range projectIDs {
		_, projectID := project_test_dep.Create(t, frame, fmt.Sprintf("user_grants_datasource_%d_%s", i, frame.UniqueResourcesID))
		project_role_test_dep.Create(t, frame, projectID, roleKey)
		addUserGrant(t, frame, userID, projectID, roleKey)
		projectIDs[i] = projectID
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
			"user_grants.#": "3",
		},
	)
}

func TestAccUserGrantsDatasource_FilterByProject(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grants")
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	_, matchingProjectID := project_test_dep.Create(t, frame, "user_grants_datasource_matching_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, matchingProjectID, roleKey)
	matchingGrantID := addUserGrant(t, frame, userID, matchingProjectID, roleKey)

	_, otherProjectID := project_test_dep.Create(t, frame, "user_grants_datasource_other_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, otherProjectID, roleKey)
	addUserGrant(t, frame, userID, otherProjectID, roleKey)

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
			"user_grants.#":              "1",
			"user_grants.0.id":           matchingGrantID,
			"user_grants.0.project_id":   matchingProjectID,
			"user_grants.0.project_name": "user_grants_datasource_matching_" + frame.UniqueResourcesID,
			"user_grants.0.role_keys.#":  "1",
			"user_grants.0.role_keys.0":  roleKey,
			"user_grants.0.state":        "USER_GRANT_STATE_ACTIVE",
		},
	)
}

func TestAccUserGrantsDatasource_FilterByRoleKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grants")
	userDep, userID := human_user_test_dep.Create(t, frame)

	matchingRoleKey := "admin_" + frame.UniqueResourcesID
	otherRoleKey := "viewer_" + frame.UniqueResourcesID
	_, matchingProjectID := project_test_dep.Create(t, frame, "user_grants_datasource_matching_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, matchingProjectID, matchingRoleKey)
	matchingGrantID := addUserGrant(t, frame, userID, matchingProjectID, matchingRoleKey)

	_, otherProjectID := project_test_dep.Create(t, frame, "user_grants_datasource_other_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, otherProjectID, otherRoleKey)
	addUserGrant(t, frame, userID, otherProjectID, otherRoleKey)

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
			"user_grants.#":             "1",
			"user_grants.0.id":          matchingGrantID,
			"user_grants.0.project_id":  matchingProjectID,
			"user_grants.0.role_keys.0": matchingRoleKey,
		},
	)
}

func TestAccUserGrantsDatasource_FilterByProjectGrant(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grants")

	roleKey := "role_" + frame.UniqueResourcesID
	_, projectID := project_test_dep.Create(t, frame, "user_grants_datasource_"+frame.UniqueResourcesID)
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
	projectGrantID := projectGrant.GetGrantId()

	_, userID := human_user_test_dep.Create(t, grantedFrame)
	resp, err := grantedFrame.AddUserGrant(grantedFrame, &management.AddUserGrantRequest{
		UserId:         userID,
		ProjectId:      projectID,
		ProjectGrantId: projectGrantID,
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
`, grantedOrgID, userID, projectGrantID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"user_grants.#":                  "1",
			"user_grants.0.id":               resp.GetUserGrantId(),
			"user_grants.0.project_id":       projectID,
			"user_grants.0.project_grant_id": projectGrantID,
			"user_grants.0.granted_org_id":   grantedOrgID,
			"user_grants.0.role_keys.0":      roleKey,
		},
	)
}

func TestAccUserGrantsDatasource_NoMatch(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_grants")
	userDep, userID := human_user_test_dep.Create(t, frame)

	roleKey := "role_" + frame.UniqueResourcesID
	_, projectID := project_test_dep.Create(t, frame, "user_grants_datasource_"+frame.UniqueResourcesID)
	project_role_test_dep.Create(t, frame, projectID, roleKey)
	addUserGrant(t, frame, userID, projectID, roleKey)

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
			"user_grants.#": "0",
		},
	)
}

func addUserGrant(t *testing.T, frame *test_utils.OrgTestFrame, userID, projectID, roleKey string) string {
	resp, err := frame.AddUserGrant(frame, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
		RoleKeys:  []string{roleKey},
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.GetUserGrantId()
}
