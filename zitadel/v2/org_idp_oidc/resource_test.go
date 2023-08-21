package org_idp_oidc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_oidc"
)

func TestAccOrgIDPJWT(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_oidc")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org_idp_oidc.DisplayNameMappingVar, exampleAttributes).AsString()
	updatedProperty := idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL.String()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		func(configProperty, _ string) string {
			return strings.Replace(resourceExample, exampleProperty, configProperty, 1)
		},
		exampleProperty, updatedProperty,
		"", "",
		true,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame), updatedProperty),
		nil, nil, "", "",
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
