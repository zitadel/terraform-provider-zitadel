package org_idp_github_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAccZITADELOrgIdPGitHub(t *testing.T) {
	ctx, err := toZitadelContext()
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	zitadelCtx := fromZitadelContext(ctx)
	getProviderByIDResponse := new(management.GetProviderByIDResponse)
	resource.Test(t, resource.TestCase{
		ProviderFactories: zitadelProviderFactories(ctx),
		CheckDestroy:      checkDestroy(ctx),
		Steps: []resource.TestStep{
			{ // Check resource can be created
				Config: fmt.Sprintf(`
resource "zitadel_org_idp_github" "%s" {
  org_id              = "%s"
  name                = "aninitialprovidername"
  client_id           = "aclientid"
  client_secret       = "a secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}
%s
`, zitadelCtx.terraformID, zitadelCtx.orgID, zitadelCtx.providerSnippet),
				RefreshState: false,
				//				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					assignGetProviderByIDResponse(ctx, getProviderByIDResponse),
					resource.ComposeAggregateTestCheckFunc(
						checkStateHasIDSet(ctx),
						checkName("aninitialprovidername", getProviderByIDResponse),
					),
				),
			}, { // Check resource can be updated
				Config: fmt.Sprintf(`
resource "zitadel_org_idp_github" "%s" {
  org_id              = "%s"
  name                = "anupdatedprovidername"
  client_id           = "aclientid"
  client_secret       = "a secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}
%s
`, zitadelCtx.terraformID, zitadelCtx.orgID, zitadelCtx.providerSnippet),
				//				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					assignGetProviderByIDResponse(ctx, getProviderByIDResponse),
					checkName("anupdatedprovidername", getProviderByIDResponse),
				),
			}, { // Check client secret can be updated
				Config: fmt.Sprintf(`
resource "zitadel_org_idp_github" "%s" {
  org_id              = "%s"
  name                = "anupdatedprovidername"
  client_id           = "aclientid"
  client_secret       = "another secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}
%s
`, zitadelCtx.terraformID, zitadelCtx.orgID, zitadelCtx.providerSnippet),
				//				ExpectNonEmptyPlan: true,
			}, { // No changes produce an empty plan
				Config: fmt.Sprintf(`
resource "zitadel_org_idp_github" "%s" {
  org_id              = "%s"
  name                = "anupdatedprovidername"
  client_id           = "aclientid"
  client_secret       = "another secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}
%s
`, zitadelCtx.terraformID, zitadelCtx.orgID, zitadelCtx.providerSnippet),
				PlanOnly: true,
				//				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func checkName(expect string, getProviderByIDResponse *management.GetProviderByIDResponse) resource.TestCheckFunc {
	return func(*terraform.State) error {
		actual := getProviderByIDResponse.GetIdp().GetName()
		if getProviderByIDResponse.GetIdp().GetName() != expect {
			return fmt.Errorf("expected name %s, actual name: %s", expect, actual)
		}
		return nil
	}
}

func checkDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		err := assignGetProviderByIDResponse(ctx, new(management.GetProviderByIDResponse))(state)
		if status.Code(err) != codes.NotFound {
			return fmt.Errorf("expected not found error but got: %w", err)
		}
		return nil
	}
}

func assignGetProviderByIDResponse(ctx context.Context, assign *management.GetProviderByIDResponse) resource.TestCheckFunc {
	return func(state *terraform.State) (err error) {
		zitadelCtx := fromZitadelContext(ctx)
		rs := state.RootModule().Resources[zitadelCtx.terraformName]
		apiProvider, err := zitadelCtx.client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: rs.Primary.ID})
		if err != nil {
			return err
		}
		*assign = *apiProvider //nolint:govet
		return nil
	}
}

func checkStateHasIDSet(ctx context.Context) resource.TestCheckFunc {
	// ZITADEL IDs have thirteen digits
	idPattern := regexp.MustCompile(`\d{13}`)
	return func(state *terraform.State) error {
		zitadelCtx := fromZitadelContext(ctx)
		return resource.TestMatchResourceAttr(zitadelCtx.terraformName, "id", idPattern)(state)
	}
}

func zitadelProviderFactories(ctx context.Context) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"zitadel": func() (*schema.Provider, error) {
			zitadelCtx := fromZitadelContext(ctx)
			return zitadelCtx.zitadelProvider, nil
		},
	}
}
