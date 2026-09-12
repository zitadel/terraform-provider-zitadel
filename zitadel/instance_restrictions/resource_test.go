package instance_restrictions_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
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

// TestAccInstanceRestrictionsRejectedByServer covers the report in issue 444:
// ZITADEL refuses allowed_languages that do not contain the instance default
// language with FailedPrecondition. The apply has to fail with that error and
// nothing may be persisted on the server.
func TestAccInstanceRestrictionsRejectedByServer(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_restrictions")
	rejectedConfig := frame.ProviderSnippet + `
resource "zitadel_instance_restrictions" "default" {
  allowed_languages                = ["de"]
  disallow_public_org_registration = true
}
`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      rejectedConfig,
				ExpectError: regexp.MustCompile(`The default language must be allowed`),
			},
			{
				Config: frame.ProviderSnippet,
				Check:  test_utils.CheckAMinute(checkRemoteAllowedLanguages(frame, nil)),
			},
		},
	})
}

// TestAccInstanceRestrictionsUpdateRejectedByServer covers the same rejection
// on the update path of a resource that is already in state.
func TestAccInstanceRestrictionsUpdateRejectedByServer(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_restrictions")
	config := func(disallow bool) string {
		return fmt.Sprintf(`
%s
resource "zitadel_instance_restrictions" "default" {
  disallow_public_org_registration = %t
}
`, frame.ProviderSnippet, disallow)
	}
	rejectedConfig := frame.ProviderSnippet + `
resource "zitadel_instance_restrictions" "default" {
  allowed_languages                = ["de"]
  disallow_public_org_registration = true
}
`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config(true),
				Check:  test_utils.CheckAMinute(checkRemoteProperty(frame)(true)),
			},
			{
				Config:      rejectedConfig,
				ExpectError: regexp.MustCompile(`The default language must be allowed`),
			},
			{
				Config: config(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					test_utils.CheckAMinute(checkRemoteProperty(frame)(false)),
					test_utils.CheckAMinute(checkRemoteAllowedLanguages(frame, nil)),
				),
			},
		},
	})
}

// TestAccInstanceRestrictionsCreateUnchanged creates the resource with a config
// that already matches the server. SetRestrictions answers that with success,
// so the create must not need to tolerate any error.
func TestAccInstanceRestrictionsCreateUnchanged(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_restrictions")
	config := func(disallow bool) string {
		return fmt.Sprintf(`
%s
resource "zitadel_instance_restrictions" "default" {
  disallow_public_org_registration = %t
}
`, frame.ProviderSnippet, disallow)
	}
	setRemoteProperty := func(value bool) func() {
		return func() {
			client, err := helper.GetAdminClient(context.Background(), frame.ClientInfo)
			if err != nil {
				t.Fatalf("failed to get client: %v", err)
			}
			if _, err := client.SetRestrictions(context.Background(), &admin.SetRestrictionsRequest{
				DisallowPublicOrgRegistration: &value,
			}); err != nil {
				t.Fatalf("setting instance restrictions failed: %v", err)
			}
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: setRemoteProperty(true),
				Config:    config(true),
				Check:     test_utils.CheckAMinute(checkRemoteProperty(frame)(true)),
			},
			{
				Config: config(false),
				Check:  test_utils.CheckAMinute(checkRemoteProperty(frame)(false)),
			},
		},
	})
}

func checkRemoteAllowedLanguages(frame *test_utils.InstanceTestFrame, expect []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := helper.GetAdminClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}
		resp, err := client.GetRestrictions(context.Background(), &admin.GetRestrictionsRequest{})
		if err != nil {
			return fmt.Errorf("getting instance restrictions failed: %w", err)
		}
		if actual := resp.GetAllowedLanguages(); !slices.Equal(actual, expect) {
			return fmt.Errorf("expected allowed languages %v, but got %v", expect, actual)
		}
		return nil
	}
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
