package application_saml_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

func TestAccAppSAML(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_saml")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, application_saml.NameVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedproperty",
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_saml.ProjectIDVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func TestAccAppSAMLMetadataUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_saml")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	initialMetadata := `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     entityID="http://example.com/saml/metadata">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="http://example.com/saml/acs"
                                     index="1" />
    </md:SPSSODescriptor>
</md:EntityDescriptor>`

	updatedMetadata := `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     entityID="http://example.com/saml/metadata-updated">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="http://example.com/saml/acs-updated"
                                     index="1" />
    </md:SPSSODescriptor>
</md:EntityDescriptor>`

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_saml" "default" {
  org_id       = data.zitadel_org.default.id
  project_id   = %q
  name         = %q
  metadata_xml = %q
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID, initialMetadata)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_saml" "default" {
  org_id       = data.zitadel_org.default.id
  project_id   = %q
  name         = %q
  metadata_xml = %q
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID, updatedMetadata)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name", frame.UniqueResourcesID),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name", frame.UniqueResourcesID),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectId string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: frame.State(state).ID, ProjectId: projectId})
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
