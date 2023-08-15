package domain_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDomain(t *testing.T) {
	resourceName := "zitadel_domain"
	initialProperty := "initial.default.127.0.0.1.sslip.io"
	updatedProperty := "updated.default.127.0.0.1.sslip.io"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id          = "%s"
  name      = "%s"
  is_primary = false
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame),
		regexp.MustCompile(fmt.Sprintf(`^%s$|^%s$`, initialProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ConcatImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(_ *terraform.State) error {
			remoteResource, err := frame.ListOrgDomains(frame, &management.ListOrgDomainsRequest{
				Queries: []*org.DomainSearchQuery{{
					Query: &org.DomainSearchQuery_DomainNameQuery{
						DomainNameQuery: &org.DomainNameQuery{
							Name: expect,
						},
					},
				}},
			})
			if err != nil {
				return err
			}
			if len(remoteResource.GetResult()) == 0 {
				return fmt.Errorf("expected to find %s, but didn't: %w", expect, test_utils.ErrNotFound)
			}
			return nil
		}
	}
}
