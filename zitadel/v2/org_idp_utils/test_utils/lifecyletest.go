package test_utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils/test_utils"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	initialProviderName = "an initial provider name"
	updatedProviderName = "an updated provider name"
	initialSecret       = "an initial secret"
	updatedSecret       = "an updated secret"
	importedSecret      = "an imported secret"
)

func RunBasicLifecyleTest(
	t *testing.T,
	frame *OrgTestFrame,
	resourceFunc func(string, string) string,
	secretAttribute string,
) {
	getProviderByIDResponse := new(management.GetProviderByIDResponse)
	initialConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(initialProviderName, initialSecret))
	updatedNameConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProviderName, initialSecret))
	updatedClientSecretConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProviderName, updatedSecret))
	resource.Test(t, resource.TestCase{
		ProviderFactories: test_utils.ZitadelProviderFactories(frame.ConfiguredProvider),
		CheckDestroy:      CheckDestroy(frame),
		Steps: []resource.TestStep{
			{ // Check first plan has a diff
				Config:             initialConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check resource is created
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					AssignGetProviderByIDResponse(frame, getProviderByIDResponse),
					resource.ComposeAggregateTestCheckFunc(
						test_utils.CheckStateHasIDSet(frame.BaseTestFrame),
						test_utils.CheckName(initialProviderName, getProviderByIDResponse),
					),
				),
			}, { // Check updating name has a diff
				Config:             updatedNameConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check name can be updated
				Config: updatedNameConfig,
				Check: resource.ComposeTestCheckFunc(
					AssignGetProviderByIDResponse(frame, getProviderByIDResponse),
					test_utils.CheckName(updatedProviderName, getProviderByIDResponse),
				),
			}, { // Check updating secret has a diff
				Config:             updatedClientSecretConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check secret can be updated
				Config: updatedClientSecretConfig,
			}, { // Expect import error if secret is not given
				ResourceName:  frame.TerraformName,
				ImportState:   true,
				ImportStateId: "123:456",
				ExpectError:   regexp.MustCompile(`123:456`),
			}, { // Expect importing works
				ResourceName: frame.TerraformName,
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					lastState := state.RootModule().Resources[frame.TerraformName].Primary
					return fmt.Sprintf("%s:%s:%s", lastState.Attributes["org_id"], lastState.ID, importedSecret), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{secretAttribute},
				Check: func(state *terraform.State) error {
					// Check the secretAttribute is imported correctly
					currentState := state.RootModule().Resources[frame.TerraformName].Primary
					actual := currentState.Attributes[secretAttribute]
					if actual != importedSecret {
						return fmt.Errorf("expected %s to be %s, but got %s", secretAttribute, importedSecret, actual)
					}
					return nil
				},
			},
		},
	})
}
