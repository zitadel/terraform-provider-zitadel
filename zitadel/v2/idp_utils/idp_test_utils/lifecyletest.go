package idp_test_utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func RunInstanceIDPLifecyleTest(t *testing.T, resourceName, secretAttribute string) {
	const importedSecret = "an_imported_secret"
	frame := test_utils.NewInstanceTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	nameProperty := test_utils.AttributeValue(t, idp_utils.NameVar, exampleAttributes).AsString()
	// Using a unique name makes the test idempotent on failures
	resourceExample = strings.Replace(resourceExample, nameProperty, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, idp_utils.IsCreationAllowedVar, exampleAttributes).True()
	exampleSecret := test_utils.AttributeValue(t, secretAttribute, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		true, false,
		exampleSecret, "an_updated_secret",
		false,
		CheckCreationAllowed(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		CheckDestroy(*frame),
		func(state *terraform.State) error {
			// Check the secret is imported correctly
			actual := frame.State(state).Attributes[secretAttribute]
			if actual != importedSecret {
				return fmt.Errorf("expected %s to be %s, but got %s", secretAttribute, importedSecret, actual)
			}
			return nil
		},
		func(state *terraform.State) (string, error) {
			lastState := frame.State(state)
			return fmt.Sprintf("%s:%s", lastState.ID, importedSecret), nil
		},
		"12345",
		secretAttribute,
	)
}
