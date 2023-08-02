package org_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccOrg(t *testing.T) {
	resourceName := "zitadel_org"
	initialProperty := "initialorgname"
	updatedProperty := "updatedorgname"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  name            = "%s"
}`, resourceName, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.Admin.GetOrgByID(frame, &admin.GetOrgByIDRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetOrg().GetName()
			if remoteResource.GetOrg().GetState() == org.OrgState_ORG_STATE_REMOVED {
				return fmt.Errorf("org is removed: %w", test_utils.ErrNotFound)
			}
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
