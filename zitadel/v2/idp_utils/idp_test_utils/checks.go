package idp_test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckProviderName(frame test_utils.InstanceTestFrame) func(interface{}) resource.TestCheckFunc {
	return func(expectName interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs := state.RootModule().Resources[frame.TerraformName]
			remoteProvider, err := frame.Client.GetProviderByID(frame, &admin.GetProviderByIDRequest{Id: rs.Primary.ID})
			if err != nil {
				return err
			}
			actual := remoteProvider.GetIdp().GetName()
			if actual != expectName {
				return fmt.Errorf("expected name %s, actual name: %s", expectName, actual)
			}
			return nil
		}
	}
}

func CheckDestroy(frame test_utils.InstanceTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := CheckProviderName(frame)("")(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}
