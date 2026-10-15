package idp_test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
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

// CheckLatestChangedEventSecret checks whether the latest eventType event of
// the IDP in state carries secretField. Unchanged secrets must not be re-sent
// on updates, because ZITADEL records every sent secret as a secret change.
func CheckLatestChangedEventSecret(frame test_utils.InstanceTestFrame, eventType, secretField string) func(bool) resource.TestCheckFunc {
	return func(expectSecret bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			idpID := frame.State(state).ID
			resp, err := frame.Client.ListEvents(frame, &admin.ListEventsRequest{
				AggregateId: frame.InstanceID,
				EventTypes:  []string{eventType},
				Limit:       1000,
			})
			if err != nil {
				return err
			}
			for _, event := range resp.GetEvents() {
				fields := event.GetPayload().GetFields()
				if fields["id"].GetStringValue() != idpID {
					continue
				}
				_, hasSecret := fields[secretField]
				if hasSecret != expectSecret {
					return fmt.Errorf("expected latest %s event to carry %s: %t, but got: %t", eventType, secretField, expectSecret, hasSecret)
				}
				return nil
			}
			return fmt.Errorf("no %s event found for IDP %s: %w", eventType, idpID, test_utils.ErrNotFound)
		}
	}
}
