package application_api_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_api"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_api/application_api_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project/project_test_dep"
)

func TestAccApplicationAPIDatasource_ID(t *testing.T) {
	datasourceName := "zitadel_application_api"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleID := test_utils.AttributeValue(t, application_api.AppIDVar, attributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	appName := "application_api_datasource_" + frame.UniqueResourcesID
	_, appID, clientID := application_api_test_dep.Create(t, frame, projectID, appName)
	config = strings.Replace(config, exampleID, appID, 1)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		nil,
		map[string]string{
			"org_id":     frame.OrgID,
			"project_id": projectID,
			"app_id":     appID,
			"name":       appName,
			"client_id":  clientID,
		},
	)
}

func TestAccApplicationAPIsDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_application_apis"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, application_api.NameVar, attributes).AsString()
	appName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:6], "\n")
	config = strings.Replace(config, exampleName, appName, 1)
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	_, appID, _ := application_api_test_dep.Create(t, frame, projectID, appName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		checkRemoteDatasourceProperty(frame, projectID, appID)(appName),
		map[string]string{
			"app_ids.0": appID,
			"app_ids.#": "1",
		},
	)
}

func TestAccApplicationAPIsDatasources_ID_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_application_apis"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, application_api.NameVar, attributes).AsString()
	appName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:6], "\n")
	config = strings.Replace(config, exampleName, "mismatch", 1)
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	_, appID, _ := application_api_test_dep.Create(t, frame, projectID, appName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		checkRemoteDatasourceProperty(frame, projectID, appID)(appName),
		map[string]string{
			"app_ids.#": "0",
		},
	)
}

func checkRemoteDatasourceProperty(frame *test_utils.OrgTestFrame, projectId, id string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: id, ProjectId: projectId})
			if err != nil {
				return err
			}
			actual := remoteResource.GetApp().GetName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
