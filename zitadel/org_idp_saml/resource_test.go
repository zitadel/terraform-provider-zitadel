package org_idp_saml_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPSAML(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_saml", "")
}

const minimalMetadataXML = `<?xml version=\"1.0\"?><md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://saml.example.com\"><md:IDPSSODescriptor WantAuthnRequestsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"><md:SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"https://saml.example.com/sso\"/></md:IDPSSODescriptor></md:EntityDescriptor>`

func TestAccOrgIdPSAMLNameIdFormat(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_saml")
	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_saml" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "%s"
  metadata_xml        = "%s"
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  name_id_format      = "SAML_NAME_ID_FORMAT_PERSISTENT"
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID, minimalMetadataXML)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name_id_format", "SAML_NAME_ID_FORMAT_PERSISTENT"),
				),
			},
		},
	})
}

func TestAccOrgIdPSAMLSignatureAlgorithm(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_saml")
	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_saml" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "%s"
  metadata_xml        = "%s"
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  with_signed_request = true
  signature_algorithm = "SAML_SIGNATURE_RSA_SHA512"
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID, minimalMetadataXML)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "signature_algorithm", "SAML_SIGNATURE_RSA_SHA512"),
					resource.TestCheckResourceAttr(frame.TerraformName, "with_signed_request", "true"),
				),
			},
		},
	})
}

func TestAccOrgIdPSAMLFederatedLogout(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_saml")
	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_saml" "default" {
  org_id                   = data.zitadel_org.default.id
  name                     = "%s"
  metadata_xml             = "%s"
  is_linking_allowed       = false
  is_creation_allowed      = true
  is_auto_creation         = false
  is_auto_update           = true
  federated_logout_enabled = true
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID, minimalMetadataXML)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "federated_logout_enabled", "true"),
				),
			},
		},
	})
}

func TestAccOrgIdPSAMLFieldUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_saml")
	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_saml" "default" {
  org_id                   = data.zitadel_org.default.id
  name                     = "%s"
  metadata_xml             = "%s"
  is_linking_allowed       = false
  is_creation_allowed      = true
  is_auto_creation         = false
  is_auto_update           = true
  name_id_format           = "SAML_NAME_ID_FORMAT_EMAIL_ADDRESS"
  federated_logout_enabled = false
  signature_algorithm      = "SAML_SIGNATURE_RSA_SHA256"
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID, minimalMetadataXML)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_saml" "default" {
  org_id                   = data.zitadel_org.default.id
  name                     = "%s"
  metadata_xml             = "%s"
  is_linking_allowed       = false
  is_creation_allowed      = true
  is_auto_creation         = false
  is_auto_update           = true
  name_id_format           = "SAML_NAME_ID_FORMAT_PERSISTENT"
  federated_logout_enabled = true
  signature_algorithm      = "SAML_SIGNATURE_RSA_SHA512"
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID, minimalMetadataXML)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name_id_format", "SAML_NAME_ID_FORMAT_EMAIL_ADDRESS"),
					resource.TestCheckResourceAttr(frame.TerraformName, "federated_logout_enabled", "false"),
					resource.TestCheckResourceAttr(frame.TerraformName, "signature_algorithm", "SAML_SIGNATURE_RSA_SHA256"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name_id_format", "SAML_NAME_ID_FORMAT_PERSISTENT"),
					resource.TestCheckResourceAttr(frame.TerraformName, "federated_logout_enabled", "true"),
					resource.TestCheckResourceAttr(frame.TerraformName, "signature_algorithm", "SAML_SIGNATURE_RSA_SHA512"),
				),
			},
		},
	})
}
