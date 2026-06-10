package application_v2_test

import (
	"fmt"
	"regexp"
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
// configuration variant, using the metadata_xml branch of the SAML
// metadata oneof. This catches regressions in the SAML builder/flattener
// and the application-type oneof dispatch that the OIDC and API tests do
// not.
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
    metadata_xml = <<EOT
%s
EOT
  }
}`, projectID, property, samlSPMetadata)
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
		// The server normalises the stored metadata XML and computes
		// login_version, so ignore both on import verification.
		"saml.0.metadata_xml",
		"saml.0.login_version",
	)
}

// samlSPMetadata is a minimal but valid SAML 2.0 Service Provider metadata
// document. Zitadel parses and validates the metadata on create, so the
// SAML acceptance test cannot use a placeholder URL or arbitrary string.
const samlSPMetadata = `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://example.com/saml/metadata">
  <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://example.com/saml/acs" index="0" isDefault="true"/>
  </md:SPSSODescriptor>
</md:EntityDescriptor>`

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

// TestAccApplicationV2_ImportWithSecret verifies the custom importer's
// secret-preservation branch: importing an OIDC app with the import id
// <app_id:org_id:client_secret> seeds the secret into the oidc block and it
// survives in state (ImportStateVerify does not ignore client_secret here).
func TestAccApplicationV2_ImportWithSecret(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_v2")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	name := "app_oidc_importsecret_" + frame.UniqueResourcesID
	config := fmt.Sprintf(`
%s
%s
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
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttrSet(frame.TerraformName, "oidc.0.client_secret"),
			},
			{ // import passing the secret as the third segment; it must persist
				Config:       config,
				ResourceName: frame.TerraformName,
				ImportState:  true,
				ImportStateIdFunc: test_utils.ChainImportStateIdFuncs(
					test_utils.ImportResourceId(frame.BaseTestFrame),
					test_utils.ImportOrgId(frame),
					test_utils.ImportStateAttribute(frame.BaseTestFrame, "oidc.0.client_secret"),
				),
				ImportStateVerify: true,
				// client_secret is the field under test and is intentionally
				// NOT ignored; these two are server-computed and may differ.
				ImportStateVerifyIgnore: []string{"oidc.0.compliance_problems", "oidc.0.login_version"},
			},
		},
	})
}

// TestAccApplicationV2_ImportSecretRejectedForSAML verifies that supplying a
// client_secret segment when importing a SAML application fails with a clear
// error (SAML applications have no client secret).
func TestAccApplicationV2_ImportSecretRejectedForSAML(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_v2")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	name := "app_saml_importrej_" + frame.UniqueResourcesID
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_application_v2" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = %q
  name       = %q
  saml {
    metadata_xml = <<EOT
%s
EOT
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, name, samlSPMetadata)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:       config,
				ResourceName: frame.TerraformName,
				ImportState:  true,
				ImportStateIdFunc: test_utils.ChainImportStateIdFuncs(
					test_utils.ImportResourceId(frame.BaseTestFrame),
					test_utils.ImportOrgId(frame),
					func(*terraform.State) (string, error) { return "somesecret", nil },
				),
				ExpectError: regexp.MustCompile(`neither an OIDC nor an API application`),
			},
		},
	})
}

// checkRemoteProperty validates the application's name via the v2 API,
// the same endpoint the resource itself reads from, so that the test
// covers the v2 wire format end-to-end.
func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			id := frame.State(state).ID
			// After a destroy the resource is gone from state and the ID
			// is empty. The v2 GetApplication RPC validates the ID length
			// and returns InvalidArgument (not NotFound) for an empty ID,
			// which the not-found destroy assertion would not recognise.
			// Treat the empty ID as the application being absent.
			if id == "" {
				return test_utils.ErrNotFound
			}
			client, err := helper.GetAppV2Client(frame.Context, frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get app v2 client: %w", err)
			}
			resp, err := client.GetApplication(frame.Context, &apppb.GetApplicationRequest{
				ApplicationId: id,
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
