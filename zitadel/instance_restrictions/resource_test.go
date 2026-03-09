package instance_restrictions_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceRestrictions(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_restrictions")
	resourceExample := `
resource "zitadel_instance_restrictions" "default" {
  disallow_public_org_registration = true
}
`
	updatedExample := `
resource "zitadel_instance_restrictions" "default" {
  disallow_public_org_registration = false
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
		checkRemoteProperty(frame),
		regexp.MustCompile(`^instance_restrictions$`),
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetAdminClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			resp, err := client.GetRestrictions(context.Background(), &admin.GetRestrictionsRequest{})
			if err != nil {
				return fmt.Errorf("getting instance restrictions failed: %w", err)
			}
			actual := resp.GetDisallowPublicOrgRegistration()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
