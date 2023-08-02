package org_idp_test_utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"
)

func RunOrgLifecyleTest(
	t *testing.T,
	frame test_utils.OrgTestFrame,
	resourceFunc func(string, string) string,
	secretAttribute string,
) {
	const importedSecret = "an_imported_secret"
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		resourceFunc,
		"an initial provider name", "an updated provider name",
		"an_initial_secret", "an_updated_secret",
		CheckProviderName(frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		CheckDestroy(frame),
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
