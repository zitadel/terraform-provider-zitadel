package application_v2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

// TestAccApplicationV2_OIDC exercises the unified zitadel_application_v2
// resource with the OIDC configuration block populated. Mirrors the v1
// TestAccAppOIDC pattern but verifies remote state through the v2
// GetApplication endpoint.
//
// We build the HCL inline rather than reading from
// examples/provider/resources/application_v2.tf because the test helper
// `test_utils.ReadExample` uses HCL's JustAttributes() which forbids
// nested blocks; the example file uses `oidc { ... }` block syntax.
func TestAccApplicationV2_OIDC(t *testing.T) {
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

  oidc {
    redirect_uris    = ["https://localhost.com/callback"]
    response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
    grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
    auth_method_type = "OIDC_AUTH_METHOD_TYPE_BASIC"
  }
}`, projectID, property)
		},
		"app_oidc_"+frame.UniqueResourcesID,
		"app_oidc_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
		// Computed-only / server-derived fields that won't be present in HCL
		// during the ImportStateVerify pass.
		"oidc.0.client_secret",
		"oidc.0.compliance_problems",
		"oidc.0.login_version",
	)
}

// TestAccApplicationV2_SAML exercises the unified resource with the SAML
// configuration variant, including the metadata_url oneof path and the
// shared login_version sub-block. This catches regressions in the SAML
// builder/flattener and the metadata oneof dispatch that the OIDC and
// API tests do not.
func TestAccApplicationV2_SAML(t *testing.T) {
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

  saml {
    metadata_url = "https://example.com/saml/metadata.xml"
    login_version {
      login_v2 {}
    }
  }
}`, projectID, property)
		},
		"app_saml_"+frame.UniqueResourcesID,
		"app_saml_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
		"saml.0.login_version",
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
			test_utils.ImportOrgId(frame),
		),
		"api.0.client_secret",
	)
}

// checkRemoteProperty validates the application's name via the v2 API,
// the same endpoint the resource itself reads from, so that the test
// covers the v2 wire format end-to-end.
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
