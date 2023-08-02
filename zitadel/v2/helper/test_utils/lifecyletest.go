package test_utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func RunLifecyleTest[P comparable](
	t *testing.T,
	frame BaseTestFrame,
	resourceFunc func(initialProperty P, initialSecret string) string,
	initialProperty, updatedProperty P,
	initialSecret, updatedSecret string,
	checkRemoteProperty func(expect P) resource.TestCheckFunc,
	idPattern *regexp.Regexp,
	checkDestroy, checkImportState resource.TestCheckFunc,
	importStateIdFunc resource.ImportStateIdFunc,
	wrongImportID,
	secretAttribute string,
) {
	var importStateVerifyIgnore []string
	initialConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(initialProperty, initialSecret))
	updatedNameConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProperty, initialSecret))
	updatedSecretConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProperty, updatedSecret))
	steps := []resource.TestStep{
		{ // Check first plan has a diff
			Config:             initialConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan just works with PlanOnly set to true
			PlanOnly: true,
		}, { // Check resource is created
			Config: initialConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				CheckAMinute(checkRemoteProperty(initialProperty)),
				CheckStateHasIDSet(frame, idPattern),
			),
		}, { // Check updating name has a diff
			Config:             updatedNameConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan just works with PlanOnly set to true
			PlanOnly: true,
		}, { // Check remote state can be updated
			Config: updatedNameConfig,
			Check:  CheckAMinute(checkRemoteProperty(updatedProperty)),
		},
	}
	if secretAttribute != "" {
		steps = append(steps, resource.TestStep{ // Check that secret has a diff
			Config:             updatedSecretConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan just works with PlanOnly set to true
			PlanOnly: true,
		}, resource.TestStep{ // Check secret can be updated
			Config: updatedSecretConfig,
		})
		importStateVerifyIgnore = []string{secretAttribute}
	}
	if wrongImportID != "" {
		steps = append(steps, resource.TestStep{ // Expect import error if secret is not given
			ResourceName:  frame.TerraformName,
			ImportState:   true,
			ImportStateId: wrongImportID,
			ExpectError:   regexp.MustCompile(wrongImportID),
		})
	}
	if checkImportState != nil {
		steps = append(steps, resource.TestStep{ // Expect importing works
			ResourceName:            frame.TerraformName,
			ImportState:             true,
			ImportStateIdFunc:       importStateIdFunc,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: importStateVerifyIgnore,
			Check:                   checkImportState,
		})
	}
	resource.Test(t, resource.TestCase{
		CheckDestroy:             CheckAMinute(checkDestroy),
		Steps:                    steps,
		ProtoV6ProviderFactories: frame.v6ProviderFactories,
		ProtoV5ProviderFactories: frame.v5ProviderFactories,
	})
}
