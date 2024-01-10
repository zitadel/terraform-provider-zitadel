package machine_user_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/machine_user"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/machine_user/machine_user_test_dep"
)

func TestAccMachineUserDatasource_ID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_machine_user")
	userName := "machine_user_datasource_" + frame.UniqueResourcesID
	userDep, userID := machine_user_test_dep.Create(t, frame, userName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		userDep,
		nil,
		map[string]string{
			"org_id":    frame.OrgID,
			"user_id":   userID,
			"user_name": userName,
		},
	)
}

func TestAccMachineUsersDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, machine_user.UserNameVar, attributes).AsString()
	exampleOrg := test_utils.AttributeValue(t, helper.OrgIDVar, attributes).AsString()
	userName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, userName, 1)
	config = strings.Replace(config, exampleOrg, frame.OrgID, 1)
	_, userID := machine_user_test_dep.Create(t, frame, userName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		checkRemoteDatasourceProperty(frame, userID)(userName),
		map[string]string{
			"user_ids.0": userID,
			"user_ids.#": "1",
		},
	)
}

func TestAccMachineUsersDatasources_ID_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, machine_user.UserNameVar, attributes).AsString()
	exampleOrg := test_utils.AttributeValue(t, helper.OrgIDVar, attributes).AsString()
	userName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, "mismatch", 1)
	config = strings.Replace(config, exampleOrg, frame.OrgID, 1)
	_, userID := machine_user_test_dep.Create(t, frame, userName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		checkRemoteDatasourceProperty(frame, userID)(userName),
		map[string]string{
			"user_ids.#": "0",
		},
	)
}

func checkRemoteDatasourceProperty(frame *test_utils.OrgTestFrame, id string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetUserByID(frame, &management.GetUserByIDRequest{Id: id})
			if err != nil {
				return err
			}
			actual := remoteResource.GetUser().GetUserName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
