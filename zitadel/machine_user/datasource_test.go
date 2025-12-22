package machine_user_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user/machine_user_test_dep"
)

func TestAccMachineUserDatasource_ID(t *testing.T) {
	datasourceName := "zitadel_machine_user"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleID := test_utils.AttributeValue(t, machine_user.UserIDVar, attributes).AsString()
	userName := "machine_user_datasource_" + frame.UniqueResourcesID
	_, userID := machine_user_test_dep.Create(t, frame, userName)
	config = strings.Replace(config, exampleID, userID, 1)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"org_id":    frame.OrgID,
			"user_id":   userID,
			"user_name": userName,
		},
	)
}

func TestAccMachineUsersDatasource_All(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	usernames := []string{"machine1_" + frame.UniqueResourcesID, "machine2_" + frame.UniqueResourcesID, "machine3_" + frame.UniqueResourcesID}
	for _, username := range usernames {
		_, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
			UserName: username,
			Name:     "Test Machine",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_machine_users" "default" {
  org_id    = "%s"
  user_name = "machine"
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

func TestAccMachineUsersDatasource_FilterByUsername(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	matchingUsername := "admin_machine_" + frame.UniqueResourcesID
	usernames := []string{matchingUsername, "user_machine_" + frame.UniqueResourcesID, "viewer_machine_" + frame.UniqueResourcesID}

	for _, username := range usernames {
		_, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
			UserName: username,
			Name:     "Test Machine",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_machine_users" "default" {
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

func TestAccMachineUsersDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	usernames := []string{"machine1_" + frame.UniqueResourcesID, "machine2_" + frame.UniqueResourcesID}
	for _, username := range usernames {
		_, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
			UserName: username,
			Name:     "Test Machine",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_machine_users" "default" {
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
