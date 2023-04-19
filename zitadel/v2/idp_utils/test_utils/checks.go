package test_utils

import (
	"fmt"

	"regexp"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ResponseProto interface {
	GetIdp() *idp.Provider
}

func CheckName(expect string, getProviderByIDResponse ResponseProto) resource.TestCheckFunc {
	return func(*terraform.State) error {
		actual := getProviderByIDResponse.GetIdp().GetName()
		if actual != expect {
			return fmt.Errorf("expected name %s, actual name: %s", expect, actual)
		}
		return nil
	}
}

func CheckStateHasIDSet(frame BaseTestFrame) resource.TestCheckFunc {
	// ZITADEL IDs have thirteen digits
	idPattern := regexp.MustCompile(`\d{13}`)
	return func(state *terraform.State) error {
		return resource.TestMatchResourceAttr(frame.TerraformName, "id", idPattern)(state)
	}
}

func CheckDestroy(ctx *InstanceTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := AssignGetProviderByIDResponse(ctx, new(admin.GetProviderByIDResponse))(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}

func AssignGetProviderByIDResponse(ctx *InstanceTestFrame, assign *admin.GetProviderByIDResponse) resource.TestCheckFunc {
	return func(state *terraform.State) (err error) {
		rs := state.RootModule().Resources[ctx.TerraformName]
		apiProvider, err := ctx.Client.GetProviderByID(ctx, &admin.GetProviderByIDRequest{Id: rs.Primary.ID})
		if err != nil {
			return err
		}
		*assign = *apiProvider //nolint:govet
		return nil
	}
}
