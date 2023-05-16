package test_utils

import (
	"fmt"
	"regexp"
	"time"

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

func RetryAMinute(check resource.TestCheckFunc) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		start := time.Now()
		for {
			err := check(state)
			if err == nil {
				return nil
			}
			if time.Since(start) > time.Minute {
				return fmt.Errorf("function failed after retrying for a minute: %w", err)
			}
			time.Sleep(time.Second)
		}
	}
}
