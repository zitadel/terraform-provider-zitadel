package test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckDestroy(frame *OrgTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := AssignGetProviderByIDResponse(frame, new(management.GetProviderByIDResponse))(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}

func AssignGetProviderByIDResponse(frame *OrgTestFrame, assign *management.GetProviderByIDResponse) resource.TestCheckFunc {
	return func(state *terraform.State) (err error) {
		rs := state.RootModule().Resources[frame.TerraformName]
		apiProvider, err := frame.GetProviderByID(frame, &management.GetProviderByIDRequest{Id: rs.Primary.ID})
		if err != nil {
			return err
		}
		*assign = *apiProvider //nolint:govet
		return nil
	}
}
