package test_utils

import (
	"fmt"
	"regexp"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// ZITADEL IDs have thirteen digits
var ZITADEL_GENERATED_ID_REGEX = regexp.MustCompile(`\d{13}`)

func CheckStateHasIDSet(frame BaseTestFrame, idPattern *regexp.Regexp) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		return resource.TestMatchResourceAttr(frame.TerraformName, "id", idPattern)(state)
	}
}

func CheckAMinute(check resource.TestCheckFunc) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		return retryAMinute(func() error {
			return check(state)
		})
	}
}

func CheckIsNotFoundFromPropertyCheck(checkRemoteProperty func(string) resource.TestCheckFunc) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := checkRemoteProperty("")(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}

func retryAMinute(try func() error) error {
	start := time.Now()
	for {
		err := try()
		if err == nil {
			return nil
		}
		if time.Since(start) > time.Minute {
			return fmt.Errorf("function failed after retrying for a minute: %w", err)
		}
		time.Sleep(time.Second)
	}
}
