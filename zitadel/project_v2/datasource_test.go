package project_v2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	projectpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_v2"
)

// TestAccProjectV2Datasource_ID looks up an existing project by project_id
// via the v2 GetProject API and verifies the returned attributes.
func TestAccProjectV2Datasource_ID(t *testing.T) {
	datasourceName := "zitadel_project_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleID := test_utils.AttributeValue(t, project_v2.ProjectIDVar, attributes).AsString()
	projectName := "project_v2_datasource_" + frame.UniqueResourcesID
	_, projectID := project_test_dep.Create(t, frame, projectName)
	config = strings.Replace(config, exampleID, projectID, 1)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"org_id":     frame.OrgID,
			"project_id": projectID,
			"name":       projectName,
		},
	)
}

// TestAccProjectsV2Datasource_Name_Match exercises the v2 ListProjects
// filter path: a name filter that matches an existing project should
// return exactly one project_id. This also covers the v1->v2 text query
// method conversion done in project_v2/funcs.go.
func TestAccProjectsV2Datasource_Name_Match(t *testing.T) {
	datasourceName := "zitadel_projects_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, project_v2.NameVar, attributes).AsString()
	projectName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// Strip everything after the first block so we only configure the
	// list datasource itself (the example file also defines a singular
	// datasource via for_each and an output, which we don't need here).
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, projectName, 1)
	_, projectID := project_test_dep.Create(t, frame, projectName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		checkRemoteDatasourceProperty(frame, projectID)(projectName),
		map[string]string{
			"project_ids.0": projectID,
			"project_ids.#": "1",
		},
	)
}

// TestAccProjectsV2Datasource_Name_Mismatch verifies the filter returns
// nothing when no project matches the supplied name.
func TestAccProjectsV2Datasource_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_projects_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, project_v2.NameVar, attributes).AsString()
	projectName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	config = strings.Replace(config, exampleName, "no-such-project", 1)
	_, projectID := project_test_dep.Create(t, frame, projectName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		checkRemoteDatasourceProperty(frame, projectID)(projectName),
		map[string]string{
			"project_ids.#": "0",
		},
	)
}

// checkRemoteDatasourceProperty hits the v2 GetProject endpoint and
// compares the live project name against the expected value, so the
// test verifies the same wire format the resource uses.
func checkRemoteDatasourceProperty(frame *test_utils.OrgTestFrame, id string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetProjectV2Client(frame.Context, frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get project v2 client: %w", err)
			}
			resp, err := client.GetProject(frame.Context, &projectpb.GetProjectRequest{ProjectId: id})
			if err != nil {
				return err
			}
			actual := resp.GetProject().GetName()
			if actual != expect {
				return fmt.Errorf("expected project name %q, got %q", expect, actual)
			}
			return nil
		}
	}
}
