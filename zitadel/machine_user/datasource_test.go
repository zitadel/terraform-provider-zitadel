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

func TestAccMachineUsersDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_machine_users"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, machine_user.UserNameVar, attributes).AsString()
	userName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, userName, 1)
	_, userID := machine_user_test_dep.Create(t, frame, userName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
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
	userName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, "mismatch", 1)
	_, userID := machine_user_test_dep.Create(t, frame, userName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
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
