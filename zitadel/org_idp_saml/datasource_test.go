package org_idp_saml_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpSamlDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_saml")
	resourceDep := fmt.Sprintf(`
resource "zitadel_org_idp_saml" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "%s"
  metadata_xml        = "<?xml version=\"1.0\"?><md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://saml.example.com\"><md:IDPSSODescriptor WantAuthnRequestsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"><md:SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"https://saml.example.com/sso\"/></md:IDPSSODescriptor></md:EntityDescriptor>"
  binding             = "SAML_BINDING_UNSPECIFIED"
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_org_idp_saml" "default" {
  org_id = data.zitadel_org.default.id
  id     = zitadel_org_idp_saml.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
