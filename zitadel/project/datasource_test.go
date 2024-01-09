package project_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project/project_test_dep"
)

func TestAccProjectDatasource_ID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project")
	projectName := "project_datasource_" + frame.UniqueResourcesID
	projectDep, projectID := project_test_dep.Create(t, frame)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		projectDep,
		nil,
		map[string]string{
			"org_id":     frame.OrgID,
			"project_id": projectID,
			"name":       projectName,
		},
	)
}

func TestAccProjectsDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_projects"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, project.NameVar, attributes).AsString()
	projectName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:7], "\n")
	config = strings.Replace(config, exampleName, projectName, 1)
	_, projectID := project_test_dep.Create(t, frame)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		checkRemoteProperty(frame)(projectName),
		map[string]string{
			"project_ids.0": projectID,
			"project_ids.#": "1",
		},
	)
}
