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

// TestAccApplicationV2_OIDC exercises the unified zitadel_application_v2
// resource with the OIDC configuration block populated. It mirrors the v1
// TestAccAppOIDC but verifies remote state through the v2 GetApplication
// endpoint.
func TestAccApplicationV2_OIDC(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_v2")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, application_v2.NameVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	// Replace the example file's data-source project_id with the live project
	// ID we just created. We do this once up front so the resourceFunc only
	// has to substitute the application name.
	staticExample := strings.ReplaceAll(
		resourceExample,
		"data.zitadel_project.default.id",
		fmt.Sprintf("%q", projectID),
	)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		test_utils.ReplaceAll(staticExample, exampleProperty, ""),
		exampleProperty, "updatedproperty",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_v2.ProjectIDVar),
			test_utils.ImportOrgId(frame),
		),
		// Computed-only / server-derived fields that won't be present in HCL
		// during the ImportStateVerify pass.
		"oidc.0.client_secret",
		"oidc.0.compliance_problems",
	)
}

// TestAccApplicationV2_API exercises the unified resource with the API
// configuration variant, ensuring the oneof dispatch on
// CreateApplicationRequest.application_type works for non-OIDC payloads.
func TestAccApplicationV2_API(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_v2")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, _ string) string {
			return fmt.Sprintf(`
resource "zitadel_application_v2" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = %q
  name       = %q

  api {
    auth_method_type = "API_AUTH_METHOD_TYPE_BASIC"
  }
}`, projectID, property)
		},
		"app_api_"+frame.UniqueResourcesID,
		"app_api_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_v2.ProjectIDVar),
			test_utils.ImportOrgId(frame),
		),
		"api.0.client_secret",
	)
}

// checkRemoteProperty validates the application's name via the v2 API — the
// same endpoint the resource itself reads from — so that we cover the v2
// wire format end-to-end.
func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetAppV2Client(frame.Context, frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get app v2 client: %w", err)
			}
			resp, err := client.GetApplication(frame.Context, &apppb.GetApplicationRequest{
				ApplicationId: frame.State(state).ID,
			})
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
