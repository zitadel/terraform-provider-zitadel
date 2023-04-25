package test_utils

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func CheckStateHasIDSet(frame BaseTestFrame) resource.TestCheckFunc {
	// ZITADEL IDs have thirteen digits
	idPattern := regexp.MustCompile(`\d{13}`)
	return func(state *terraform.State) error {
		return resource.TestMatchResourceAttr(frame.TerraformName, "id", idPattern)(state)
	}
}
