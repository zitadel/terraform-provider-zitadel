package test_utils

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func RunLifecyleTest[P comparable](
	t *testing.T,
	frame BaseTestFrame,
	resourceFunc func(initialProperty P, initialSecret string) string,
	initialProperty, updatedProperty P,
	secretAttribute, initialSecret, updatedSecret string,
	allowNonEmptyPlan bool,
	checkRemoteProperty func(expect P) resource.TestCheckFunc,
	idPattern *regexp.Regexp,
	checkDestroy resource.TestCheckFunc,
	importStateIdFunc resource.ImportStateIdFunc,
) {
	initialConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(initialProperty, initialSecret))
	updatedPropertyConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProperty, initialSecret))
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
			Config:             updatedPropertyConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan just works with PlanOnly set to true
			PlanOnly: true,
		}, { // Check remote state can be updated
			Config: updatedPropertyConfig,
			Check:  CheckAMinute(checkRemoteProperty(updatedProperty)),
		},
	}
	if importStateIdFunc != nil {
		steps = append(steps, resource.TestStep{ // Expect importing works
			Config:            updatedPropertyConfig,
			ResourceName:      frame.TerraformName,
			ImportState:       true,
			ImportStateIdFunc: importStateIdFunc,
			ImportStateVerify: true,
		})
	}

	if secretAttribute != "" {
		updatedSecretConfig := fmt.Sprintf("%s\n%s", frame.ProviderSnippet, resourceFunc(updatedProperty, updatedSecret))
		steps = append(steps, resource.TestStep{ // Check that secret has a diff
			Config:             updatedSecretConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan only works with PlanOnly set to true
			PlanOnly: true,
		}, resource.TestStep{ // Check secret can be updated
			Config: updatedSecretConfig,
		})
	}
	resource.ParallelTest(t, resource.TestCase{
		CheckDestroy: CheckAMinute(checkDestroy),
		Steps:        steps,
		ErrorCheck: func(err error) error {
			if err != nil && allowNonEmptyPlan && strings.Contains(err.Error(), "After applying this test step and performing a `terraform refresh`, the plan was not empty") {
				t.Logf("Ignoring non-empty plan error because we can't guarantee consistency: %s", err.Error())
				return nil
			}
			return err
		},
		ProtoV6ProviderFactories: frame.v6ProviderFactories,
		ProtoV5ProviderFactories: frame.v5ProviderFactories,
	})
}
