package pat_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/pat"
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
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	org_id              = "%s"
	user_id = "%s"
  	expiration_date = "%s"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, userID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		test_utils.ConcatImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, pat.UserIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, pat.TokenVar),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPersonalAccessTokenByIDs(frame, &management.GetPersonalAccessTokenByIDsRequest{
				UserId:  userID,
				TokenId: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := resp.GetToken().GetExpirationDate().AsTime().Format("2006-01-02T15:04:05Z")
			if expect != actual {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
