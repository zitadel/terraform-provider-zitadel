package org_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func checkRemoteProperty(frame *test_utils.OrgTestFrame, id func(state *terraform.State) string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.Admin.GetOrgByID(frame, &admin.GetOrgByIDRequest{Id: id(state)})
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
