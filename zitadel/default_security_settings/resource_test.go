package default_security_settings_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultSecuritySettings(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_security_settings")
	resourceExample := `
resource "zitadel_default_security_settings" "default" {
  enable_impersonation = true
}
`
	updatedExample := `
resource "zitadel_default_security_settings" "default" {
  enable_impersonation = false
}
`
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		func(property bool, secret string) string {
			if property {
				return resourceExample
			}
			return updatedExample
		},
		true, false,
		"", "", "",
		false,
		checkRemoteProperty(*frame),
		regexp.MustCompile(`^default_security_settings$`),
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			resp, err := client.GetSecuritySettings(context.Background(), &settingsv2.GetSecuritySettingsRequest{})
			if err != nil {
				return fmt.Errorf("getting security settings failed: %w", err)
			}
			actual := resp.GetSettings().GetEnableImpersonation()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
