package human_user_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccHumanUsersDatasource_All(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	usernames := []string{"user1_" + frame.UniqueResourcesID, "user2_" + frame.UniqueResourcesID, "user3_" + frame.UniqueResourcesID}
	for _, username := range usernames {
		_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
			UserName: username,
			Profile: &management.ImportHumanUserRequest_Profile{
				FirstName: "Test",
				LastName:  "User",
			},
			Email: &management.ImportHumanUserRequest_Email{
				Email:           username + "@example.com",
				IsEmailVerified: true,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id     = "%s"
  user_name  = "user"
  user_name_method = "TEXT_QUERY_METHOD_CONTAINS"
}
`, frame.OrgID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{},
	)
}

func TestAccHumanUsersDatasource_FilterByUsername(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	matchingUsername := "admin_" + frame.UniqueResourcesID
	usernames := []string{matchingUsername, "user_" + frame.UniqueResourcesID, "viewer_" + frame.UniqueResourcesID}

	for _, username := range usernames {
		_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
			UserName: username,
			Profile: &management.ImportHumanUserRequest_Profile{
				FirstName: "Test",
				LastName:  "User",
			},
			Email: &management.ImportHumanUserRequest_Email{
				Email:           username + "@example.com",
				IsEmailVerified: true,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id    = "%s"
  user_name = "%s"
}
`, frame.OrgID, matchingUsername)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		checkUserExists(frame, matchingUsername),
		map[string]string{
			"user_ids.#": "1",
		},
	)
}

func TestAccHumanUsersDatasource_FilterByEmail(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	username := "emailuser_" + frame.UniqueResourcesID
	email := "test_" + frame.UniqueResourcesID + "@example.com"

	_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
		UserName: username,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName: "Test",
			LastName:  "User",
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           email,
			IsEmailVerified: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id = "%s"
  email  = "%s"
}
`, frame.OrgID, email)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"user_ids.#": "1",
		},
	)
}

func TestAccHumanUsersDatasource_FilterByFirstName(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	username := "nameuser_" + frame.UniqueResourcesID
	firstName := "UniqueFirst_" + frame.UniqueResourcesID

	_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
		UserName: username,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName: firstName,
			LastName:  "User",
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           username + "@example.com",
			IsEmailVerified: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id     = "%s"
  first_name = "%s"
}
`, frame.OrgID, firstName)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"user_ids.#": "1",
		},
	)
}

func TestAccHumanUsersDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	usernames := []string{"user1_" + frame.UniqueResourcesID, "user2_" + frame.UniqueResourcesID}
	for _, username := range usernames {
		_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
			UserName: username,
			Profile: &management.ImportHumanUserRequest_Profile{
				FirstName: "Test",
				LastName:  "User",
			},
			Email: &management.ImportHumanUserRequest_Email{
				Email:           username + "@example.com",
				IsEmailVerified: true,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id    = "%s"
  user_name = "nonexistent"
}
`, frame.OrgID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"user_ids.#": "0",
		},
	)
}

func TestAccHumanUsersDatasource_AllNoFilter(t *testing.T) {
	datasourceName := "zitadel_human_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	usernames := []string{"user1_" + frame.UniqueResourcesID, "user2_" + frame.UniqueResourcesID, "user3_" + frame.UniqueResourcesID}
	for _, username := range usernames {
		_, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
			UserName: username,
			Profile: &management.ImportHumanUserRequest_Profile{
				FirstName: "Test",
				LastName:  "User",
			},
			Email: &management.ImportHumanUserRequest_Email{
				Email:           username + "@example.com",
				IsEmailVerified: true,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_human_users" "default" {
  org_id = "%s"
}
`, frame.OrgID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{},
	)
}

func checkUserExists(frame *test_utils.OrgTestFrame, expectedUsername string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.ListUsers(frame, &management.ListUsersRequest{})
		if err != nil {
			return err
		}

		for _, user := range resp.Result {
			if user.UserName == expectedUsername {
				return nil
			}
		}

		return fmt.Errorf("expected user %s not found", expectedUsername)
	}
}
