package test_utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func RunDatasourceTest(
	t *testing.T,
	frame BaseTestFrame,
	config string,
	awaitRemoteResource resource.TestCheckFunc,
	expectProperties map[string]string,
) {
	var checks []resource.TestCheckFunc
	if awaitRemoteResource != nil {
		checks = append(checks, CheckAMinute(awaitRemoteResource))
	}
	for k, v := range expectProperties {
		checks = append(checks, resource.TestCheckResourceAttr("data."+frame.TerraformName, k, v))
	}
	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf("%s\n%s", frame.ProviderSnippet, config),
			Check:  resource.ComposeAggregateTestCheckFunc(checks...),
		}},
		ProtoV6ProviderFactories: frame.v6ProviderFactories,
		ProtoV5ProviderFactories: frame.v5ProviderFactories,
	})
}
