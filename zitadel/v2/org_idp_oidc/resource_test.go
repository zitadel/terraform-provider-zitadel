package org_idp_oidc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_oidc"
)

func TestAccOrgIDPOIDC(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_oidc")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org_idp_oidc.DisplayNameMappingVar, exampleAttributes).AsString()
	updatedProperty := idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL.String()
	exampleSecret := test_utils.AttributeValue(t, idp_utils.ClientSecretVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, updatedProperty,
		idp_utils.ClientSecretVar, exampleSecret, "an updated secret",
		true,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame), updatedProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, idp_utils.ClientSecretVar),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetOrgIDPByID(frame, &management.GetOrgIDPByIDRequest{
				Id: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := resp.GetIdp().GetOidcConfig().GetDisplayNameMapping().String()
			if expect != actual {
				return fmt.Errorf("expected jwt endpoint %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
