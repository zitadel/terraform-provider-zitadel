package instance_member_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/instance_member"
)

func TestAccInstanceMember(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_instance_member")
	userDep, userID := human_user_test_dep.Create(t, frame)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, instance_member.RolesVar, exampleAttributes).AsValueSlice()[0].AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "IAM_OWNER_VIEWER",
		"", "", "",
		true,
		checkRemoteProperty(*frame, userID),
		regexp.MustCompile(fmt.Sprintf("^%s_%s$", helper.ZitadelGeneratedIdPattern, helper.ZitadelGeneratedIdPattern)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		test_utils.ImportStateAttribute(frame.BaseTestFrame, instance_member.UserIDVar),
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
