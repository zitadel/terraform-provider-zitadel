package user_metadata_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/human_user/human_user_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/user_metadata"
)

func TestAccUserMetadata(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_metadata")
	userDep, userID := human_user_test_dep.Create(t, frame)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	keyProperty := test_utils.AttributeValue(t, user_metadata.KeyVar, exampleAttributes).AsString()
	exampleProperty := test_utils.AttributeValue(t, user_metadata.ValueVar, exampleAttributes).AsString()
	updatedProperty := "YW5vdGhlciB2YWx1ZQ=="
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, userDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame, userID, keyProperty),
		regexp.MustCompile(fmt.Sprintf(`^%s_%s$`, helper.ZitadelGeneratedIdPattern, keyProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, userID, keyProperty), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, user_metadata.UserIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, user_metadata.KeyVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, userID, key string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetUserMetadata(frame, &management.GetUserMetadataRequest{
				Id:  userID,
				Key: key,
			})
			if err != nil {
				return err
			}
			actual := helper.Base64Encode(resp.GetMetadata().GetValue())
			if expect != actual {
				return fmt.Errorf("expected role %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
