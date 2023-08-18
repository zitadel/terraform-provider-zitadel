package action_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccAction(t *testing.T) {
	resourceName := "zitadel_action"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	resourceExample, exampleAttributes := frame.ReadExample(t, test_utils.Resources, frame.ResourceType)
	// name must be unique
	nameAttribute := test_utils.AttributeValue(t, "name", exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, nameAttribute, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, "script", exampleAttributes).AsString()
	updatedProperty := "updatedproperty"
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf("%s\n%s", frame.OrgExampleDatasource, strings.Replace(resourceExample, exampleProperty, configProperty, 1))
		},
		exampleProperty, updatedProperty,
		"", "",
		true,
		checkRemoteProperty(frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAction(frame, &management.GetActionRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetAction().GetScript()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
