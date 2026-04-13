package application_saml_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_saml/application_saml_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

// TestAccApplicationSAMLDatasource_LoginVersion reproduces a bug where the
// datasource crashes if the SAML app has login_version set, because the
// datasource schema does not declare the login_version field.
func TestAccApplicationSAMLDatasource_LoginVersion(t *testing.T) {
	datasourceName := "zitadel_application_saml"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	appName := "saml_ds_loginver_" + frame.UniqueResourcesID

	// Create a SAML app via the API
	_, appID := application_saml_test_dep.Create(t, frame, projectID, appName)

	// Set login_version on the app via the API so the server returns it on read
	_, err := frame.UpdateSAMLAppConfig(frame, &management.UpdateSAMLAppConfigRequest{
		ProjectId: projectID,
		AppId:     appID,
		Metadata:  &management.UpdateSAMLAppConfigRequest_MetadataXml{MetadataXml: []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"" + appName + "\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://example.com/saml/cas\" index=\"1\" />\n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")},
		LoginVersion: &app.LoginVersion{
			Version: &app.LoginVersion_LoginV1{
				LoginV1: &app.LoginV1{},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to set login_version on SAML app: %v", err)
	}

	// Now try to read via the datasource — this should crash if login_version
	// is not in the datasource schema
	config := fmt.Sprintf(`
%s
%s
data "zitadel_application_saml" "test" {
  org_id     = data.zitadel_org.default.id
  project_id = %q
  app_id     = %q
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, appID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.zitadel_application_saml.test", "name", appName),
					resource.TestCheckResourceAttr("data.zitadel_application_saml.test", "login_version.#", "1"),
					resource.TestCheckResourceAttr("data.zitadel_application_saml.test", "login_version.0.login_v1", "true"),
				),
			},
		},
	})
}
