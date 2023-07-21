package org_idp_jwt_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccOrgIDPJWT(t *testing.T) {
	resourceName := "zitadel_org_idp_jwt"
	initialProperty := "https://initialproperty.com"
	updatedProperty := "https://updatedproperty.com"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(configProperty, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id          = "%s"
  name            = "%s"
  styling_type  = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint  = "%s"
  issuer        = "https://google.com"
  keys_endpoint = "https://jwtendpoint.com/keys"
  header_name   = "x-auth-token"
  auto_register = false
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetOrgIDPByID(frame, &management.GetOrgIDPByIDRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetIdp().GetJwtConfig().GetJwtEndpoint()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
