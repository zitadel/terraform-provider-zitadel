package application_v2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_v2"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

// createOIDCAppForDatasource is a small inline test helper that creates
// a real OIDC application via the v2 API. The package does not have a
// dedicated application_v2_test_dep package because the dependency we
// need here is small enough to inline, and the v2 API does not require
// the indirection through the v1 management client.
func createOIDCAppForDatasource(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) string {
	t.Helper()
	client, err := helper.GetAppV2Client(frame.Context, frame.ClientInfo)
	if err != nil {
		t.Fatalf("failed to get app v2 client: %v", err)
	}
	resp, err := client.CreateApplication(frame.Context, &apppb.CreateApplicationRequest{
		ProjectId: projectID,
		Name:      name,
		ApplicationType: &apppb.CreateApplicationRequest_OidcConfiguration{
			OidcConfiguration: &apppb.CreateOIDCApplicationRequest{
				RedirectUris:  []string{"https://localhost.com/callback"},
				ResponseTypes: []apppb.OIDCResponseType{apppb.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:    []apppb.OIDCGrantType{apppb.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create application via v2 API: %v", err)
	}
	return resp.GetApplicationId()
}

// TestAccApplicationV2Datasource_ID looks up an application by app_id
// through the singular zitadel_application_v2 datasource and verifies
// the returned attributes.
func TestAccApplicationV2Datasource_ID(t *testing.T) {
	datasourceName := "zitadel_application_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleID := test_utils.AttributeValue(t, application_v2.AppIDVar, attributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	appName := "application_v2_datasource_" + frame.UniqueResourcesID
	appID := createOIDCAppForDatasource(t, frame, projectID, appName)
	config = strings.Replace(config, exampleID, appID, 1)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		nil,
		map[string]string{
			"app_id":     appID,
			"project_id": projectID,
			"name":       appName,
		},
	)
}

// TestAccApplicationsV2Datasource_Name_Match verifies that the list
// datasource returns exactly one matching app_id when the name filter
// matches a freshly created application.
func TestAccApplicationsV2Datasource_Name_Match(t *testing.T) {
	datasourceName := "zitadel_applications_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, application_v2.NameVar, attributes).AsString()
	appName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// Trim down to just the list datasource block; the example also
	// includes a for_each singular lookup and an output we do not need.
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	// Replace example references with live values: the example references
	// data.zitadel_project.default.id for project_id; we substitute the
	// real project id directly so the test does not depend on an
	// additional data source.
	config = strings.Replace(config, "data.zitadel_project.default.id", fmt.Sprintf("%q", projectID), 1)
	config = strings.Replace(config, exampleName, appName, 1)
	appID := createOIDCAppForDatasource(t, frame, projectID, appName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		checkRemoteDatasourceProperty(frame, appID)(appName),
		map[string]string{
			"app_ids.0": appID,
			"app_ids.#": "1",
		},
	)
}

// TestAccApplicationsV2Datasource_Name_Mismatch verifies the list filter
// returns zero results when no application name matches.
func TestAccApplicationsV2Datasource_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_applications_v2"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, application_v2.NameVar, attributes).AsString()
	appName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	config = strings.Join(strings.Split(config, "\n")[0:5], "\n")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	config = strings.Replace(config, "data.zitadel_project.default.id", fmt.Sprintf("%q", projectID), 1)
	config = strings.Replace(config, exampleName, "no-such-application", 1)
	appID := createOIDCAppForDatasource(t, frame, projectID, appName)
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		checkRemoteDatasourceProperty(frame, appID)(appName),
		map[string]string{
			"app_ids.#": "0",
		},
	)
}

// checkRemoteDatasourceProperty hits the v2 GetApplication endpoint so
// the datasource tests verify behaviour against the same API the
// resource itself uses, not a v1 management call.
func checkRemoteDatasourceProperty(frame *test_utils.OrgTestFrame, id string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetAppV2Client(frame.Context, frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get app v2 client: %w", err)
			}
			resp, err := client.GetApplication(frame.Context, &apppb.GetApplicationRequest{ApplicationId: id})
			if err != nil {
				return err
			}
			actual := resp.GetApplication().GetName()
			if actual != expect {
				return fmt.Errorf("expected application name %q, got %q", expect, actual)
			}
			return nil
		}
	}
}
