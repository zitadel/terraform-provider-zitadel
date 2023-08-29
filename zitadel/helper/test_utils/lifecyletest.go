package test_utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func RunLifecyleTest[P comparable](
	t *testing.T,
	frame BaseTestFrame,
	datasources []string,
	resourceFunc func(property P, secret string) string,
	exampleProperty, updatedProperty P,
	secretAttribute, exampleSecret, updatedSecret string,
	allowNonEmptyPlan bool,
	checkRemoteProperty func(expect P) resource.TestCheckFunc,
	idPattern *regexp.Regexp,
	checkDestroy resource.TestCheckFunc,
	importStateIdFunc resource.ImportStateIdFunc,
	importStateVerifyIgnore ...string,
) {
	exampleConfig := fmt.Sprintf("%s\n%s\n%s", frame.ProviderSnippet, strings.Join(datasources, "\n"), resourceFunc(exampleProperty, exampleSecret))
	updatedPropertyConfig := fmt.Sprintf("%s\n%s\n%s", frame.ProviderSnippet, strings.Join(datasources, "\n"), resourceFunc(updatedProperty, exampleSecret))
	steps := []resource.TestStep{
		{ // Check first plan has a diff
			Config:             exampleConfig,
			ExpectNonEmptyPlan: true,
			// ExpectNonEmptyPlan just works with PlanOnly set to true
			PlanOnly: true,
		}, { // Check resource is created
			Config: exampleConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				CheckAMinute(checkRemoteProperty(exampleProperty)),
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
			Config:                  updatedPropertyConfig,
			ResourceName:            frame.TerraformName,
			ImportState:             true,
			ImportStateIdFunc:       importStateIdFunc,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: importStateVerifyIgnore,
		})
	}
	if secretAttribute != "" {
		updatedSecretConfig := fmt.Sprintf("%s\n%s\n%s", frame.ProviderSnippet, strings.Join(datasources, "\n"), resourceFunc(updatedProperty, updatedSecret))
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
	resource.ParallelTest(t, resource.TestCase{
		CheckDestroy: CheckAMinute(checkDestroy),
		Steps:        steps,
		ErrorCheck: func(err error) error {
			if err != nil && allowNonEmptyPlan && os.Getenv("CI") == "true" && strings.Contains(err.Error(), "After applying this test step and performing a `terraform refresh`, the plan was not empty") {
				t.Logf("Ignoring non-empty plan error because we can't guarantee consistency: %s", err.Error())
				return nil
			}
			return err
		},
		ProtoV6ProviderFactories: frame.v6ProviderFactories,
	})
}
