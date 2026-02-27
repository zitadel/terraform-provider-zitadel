package org_idp_jwt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_jwt"
)

func TestAccOrgIDPJWT(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_jwt")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org_idp_jwt.JwtEndpointVar, exampleAttributes).AsString()
	updatedProperty := "https://example.com/updated"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		func(configProperty, _ string) string {
			return strings.Replace(resourceExample, exampleProperty, configProperty, 1)
		},
		exampleProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
	)
}

func TestAccOrgIDPJWTNameUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_jwt")

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_jwt" "default" {
  org_id        = data.zitadel_org.default.id
  name          = "initial_%s"
  styling_type  = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint  = "https://example.com/jwt"
  issuer        = "https://example.com"
  keys_endpoint = "https://example.com/keys"
  header_name   = "x-auth-token"
  auto_register = false
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_jwt" "default" {
  org_id        = data.zitadel_org.default.id
  name          = "updated_%s"
  styling_type  = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint  = "https://example.com/jwt"
  issuer        = "https://example.com"
  keys_endpoint = "https://example.com/keys"
  header_name   = "x-auth-token"
  auto_register = true
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name", "initial_"+frame.UniqueResourcesID),
					resource.TestCheckResourceAttr(frame.TerraformName, "auto_register", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "name", "updated_"+frame.UniqueResourcesID),
					resource.TestCheckResourceAttr(frame.TerraformName, "auto_register", "true"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetOrgIDPByID(frame, &management.GetOrgIDPByIDRequest{
				Id: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := resp.GetIdp().GetJwtConfig().GetJwtEndpoint()
			if expect != actual {
				return fmt.Errorf("expected jwt endpoint %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
