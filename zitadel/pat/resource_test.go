package pat_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user/machine_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/pat"
)

func TestAccPersonalAccessToken(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_personal_access_token")
	userDep, userID := machine_user_test_dep.Create(t, frame, frame.UniqueResourcesID)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, pat.ExpirationDateVar, exampleAttributes).AsString()
	updatedProperty := "2051-01-01T00:00:00Z"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame, userID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID), ""),
		test_utils.ChainImportStateIdFuncs(
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
