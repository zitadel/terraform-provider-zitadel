package pat_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccPersonalAccessToken(t *testing.T) {
	resourceName := "zitadel_personal_access_token"
	initialProperty := "2050-01-01T00:00:00Z"
	updatedProperty := "2051-01-01T00:00:00Z"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	user, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
		UserName: frame.UniqueResourcesID,
		Name:     "Don't care",
	})
	userID := user.GetUserId()
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(cfg, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	org_id              = "%s"
	user_id = "%s"
  	expiration_date = "%s"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, userID, cfg)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame, userID),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID string) func(interface{}) resource.TestCheckFunc {
	return func(expected interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPersonalAccessTokenByIDs(frame, &management.GetPersonalAccessTokenByIDsRequest{
				UserId:  userID,
				TokenId: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := resp.GetToken().GetExpirationDate().AsTime().Format("2006-01-02T15:04:05Z")
			if expected != actual {
				return fmt.Errorf("expected %s, but got %s", expected, actual)
			}
			return nil
		}
	}
}
