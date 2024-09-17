package org_idp_test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func CheckCreationAllowed(frame test_utils.OrgTestFrame) func(bool) resource.TestCheckFunc {
	return func(expectAllowed bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteProvider, err := frame.GetProviderByID(frame, &management.GetProviderByIDRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteProvider.GetIdp().GetConfig().GetOptions().GetIsCreationAllowed()
			if actual != expectAllowed {
				return fmt.Errorf("expected creation allowed to be %t, but got %t", expectAllowed, actual)
			}
			return nil
		}
	}
}

func CheckDestroy(frame test_utils.OrgTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := CheckCreationAllowed(frame)(false)(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}
