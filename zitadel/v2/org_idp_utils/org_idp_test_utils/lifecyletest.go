package org_idp_test_utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"
)

func RunOrgLifecyleTest(t *testing.T, resourceName, secretAttribute string) {
	const importedSecret = "an_imported_secret"
	frame := test_utils.NewOrgTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	nameProperty := test_utils.AttributeValue(t, idp_utils.NameVar, exampleAttributes).AsString()
	// Using a unique name makes the test idempotent on failures
	resourceExample = strings.Replace(resourceExample, nameProperty, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, idp_utils.IsCreationAllowedVar, exampleAttributes).True()
	exampleSecret := test_utils.AttributeValue(t, secretAttribute, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		true, false,
		exampleSecret, "an_updated_secret",
		false,
		CheckCreationAllowed(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		CheckDestroy(*frame),
		func(state *terraform.State) error {
			// Check the secretAttribute is imported correctly
			actual := frame.State(state).Attributes[secretAttribute]
			if actual != importedSecret {
				return fmt.Errorf("expected %s to be %s, but got %s", secretAttribute, importedSecret, actual)
			}
			return nil
		},
		func(state *terraform.State) (string, error) {
			lastState := frame.State(state)
			return fmt.Sprintf("%s:%s:%s", lastState.Attributes[org_idp_utils.OrgIDVar], lastState.ID, importedSecret), nil
		},
		"123:456",
		secretAttribute,
	)
}
