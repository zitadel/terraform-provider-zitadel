package test_utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	frame *InstanceTestFrame,
	resourceFunc func(string, string) string,
) {
	getProviderByIDResponse := new(admin.GetProviderByIDResponse)
	initialConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(initialProviderName, initialSecret))
	updatedNameConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProviderName, initialSecret))
	updatedClientSecretConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProviderName, updatedSecret))
	resource.Test(t, resource.TestCase{
		ProviderFactories: ZitadelProviderFactories(frame.ConfiguredProvider),
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
						CheckStateHasIDSet(frame.BaseTestFrame),
						CheckName(initialProviderName, getProviderByIDResponse),
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
					CheckName(updatedProviderName, getProviderByIDResponse),
				),
			}, { // Check updating client secret has a diff
				Config:             updatedClientSecretConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check client secret can be updated
				Config: updatedClientSecretConfig,
			}, { // Expect import error if client secret is not given
				ResourceName:  frame.TerraformName,
				ImportState:   true,
				ImportStateId: "12345",
				ExpectError:   regexp.MustCompile(`12345`),
			}, { // Expect importing works
				ResourceName: frame.TerraformName,
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					lastState := state.RootModule().Resources[frame.TerraformName].Primary
					return fmt.Sprintf("%s:%s", lastState.ID, importedSecret), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
				Check: func(state *terraform.State) error {
					// Check the client_secret is imported correctly
					currentState := state.RootModule().Resources[frame.TerraformName].Primary
					actual := currentState.Attributes["client_secret"]
					if actual != importedSecret {
						return fmt.Errorf("expected client_secret to be %s, but got %s", importedSecret, actual)
					}
					return nil
				},
			},
		},
	})
}
