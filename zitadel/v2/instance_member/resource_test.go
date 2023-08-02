package instance_member_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccInstanceMember(t *testing.T) {
	resourceName := "zitadel_instance_member"
	initialProperty := "IAM_OWNER"
	updatedProperty := "IAM_OWNER_VIEWER"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	user, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
		UserName: frame.UniqueResourcesID,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName: "Don't",
			LastName:  "Care",
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           "dont@care.com",
			IsEmailVerified: true,
		},
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
	user_id = "%s"
  	roles   = ["%s"]
}`, resourceName, frame.UniqueResourcesID, userID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame, userID),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.Admin.ListIAMMembers(frame, &admin.ListIAMMembersRequest{
				Queries: []*member.SearchQuery{{
					Query: &member.SearchQuery_UserIdQuery{UserIdQuery: &member.UserIDQuery{UserId: userID}},
				}},
			})
			if err != nil {
				return err
			}
			if len(resp.Result) == 0 || len(resp.Result[0].Roles) == 0 {
				return fmt.Errorf("expected 1 user with 1 role, but got %d: %w", len(resp.Result), test_utils.ErrNotFound)
			}
			actual := resp.Result[0].Roles[0]
			if expect != actual {
				return fmt.Errorf("expected role %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
