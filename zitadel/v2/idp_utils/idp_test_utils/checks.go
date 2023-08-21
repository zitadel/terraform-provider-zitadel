package idp_test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func CheckCreationAllowed(frame test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expectAllowed bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteProvider, err := frame.Client.GetProviderByID(frame, &admin.GetProviderByIDRequest{Id: frame.State(state).ID})
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

func CheckDestroy(frame test_utils.InstanceTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := CheckCreationAllowed(frame)(true)(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}
