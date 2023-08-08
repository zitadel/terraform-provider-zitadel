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
	checkEventualRemoteProperty func(expect P) resource.TestCheckFunc,
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
	steps := []resource.TestStep{{
		// Check first plan has a diff
		Config:             initialConfig,
		ExpectNonEmptyPlan: true,
		PlanOnly:           true,
	}, {
		// Check resource is created
		// Eventual consistency doesn't allow us to expect not empty plans directly after apply
		// Instead we await the remote property here with retries and expect an empty plan in the next step
		RefreshState: false,
		Config:       initialConfig,
		Check: resource.ComposeAggregateTestCheckFunc(
			CheckAMinute(checkEventualRemoteProperty(initialProperty)),
			CheckStateHasIDSet(frame, idPattern),
		),
	}, {
		// We expect an empty plan because we awaited eventual consistency above
		Config:             initialConfig,
		PlanOnly:           true,
		ExpectNonEmptyPlan: false,
	}, {
		// Check updating name has a diff
		Config:             updatedNameConfig,
		ExpectNonEmptyPlan: true,
		// ExpectNonEmptyPlan just works with PlanOnly set to true
		PlanOnly: true,
	}, {
		// Check remote state can be updated
		RefreshState: false,
		Config:       updatedNameConfig,
		Check:        CheckAMinute(checkEventualRemoteProperty(updatedProperty)),
	}, {
		// We expect an empty plan because we awaited eventual consistency above
		Config:             updatedNameConfig,
		ExpectNonEmptyPlan: false,
		// ExpectNonEmptyPlan just works with PlanOnly set to true
		PlanOnly: true,
	}}
	if secretAttribute != "" {
		steps = append(steps,
			// Check that secret has a diff
			resource.TestStep{
				Config:             updatedSecretConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
				// Check secret can be updated
			}, resource.TestStep{
				Config: updatedSecretConfig,
				// We can't exect consistency here, but we can also not query secrets, so we just skip refreshing the state
				RefreshState: false,
			},
		)
		importStateVerifyIgnore = []string{secretAttribute}
	}
	if wrongImportID != "" {
		steps = append(steps,
			// Expect import error if secret is not given
			resource.TestStep{
				ResourceName:  frame.TerraformName,
				ImportState:   true,
				ImportStateId: wrongImportID,
				ExpectError:   regexp.MustCompile(wrongImportID),
			})
	}
	if checkImportState != nil {
		steps = append(steps,
			// Expect importing works
			resource.TestStep{
				ResourceName:            frame.TerraformName,
				ImportState:             true,
				ImportStateIdFunc:       importStateIdFunc,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
				Check:                   checkImportState,
			})
	}
	resource.ParallelTest(t, resource.TestCase{
		CheckDestroy:             CheckAMinute(checkDestroy),
		Steps:                    steps,
		ProtoV6ProviderFactories: frame.v6ProviderFactories,
		ProtoV5ProviderFactories: frame.v5ProviderFactories,
	})
}
