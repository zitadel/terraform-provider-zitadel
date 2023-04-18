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
	initialConfig := fmt.Sprintf(`%s
resource "%s" "%s" {
  org_id              = "%s"
  name                = "aninitialprovidername"
  client_id           = "aclientid"
  client_secret       = "a secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, zitadelCtx.providerSnippet, zitadelCtx.terraformType, zitadelCtx.terraformID, zitadelCtx.orgID)
	updatedNameConfig := fmt.Sprintf(`%s
resource "%s" "%s" {
  org_id              = "%s"
  name                = "anupdatedprovidername"
  client_id           = "aclientid"
  client_secret       = "a secret"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, zitadelCtx.providerSnippet, zitadelCtx.terraformType, zitadelCtx.terraformID, zitadelCtx.orgID)
	updatedSecret := "another secret"
	importedSecret := "an imported secret"
	updatedClientSecretConfig := fmt.Sprintf(`%s
resource "%s" "%s" {
  org_id              = "%s"
  name                = "anupdatedprovidername"
  client_id           = "aclientid"
  client_secret       = "%s"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, zitadelCtx.providerSnippet, zitadelCtx.terraformType, zitadelCtx.terraformID, zitadelCtx.orgID, updatedSecret)
	resource.Test(t, resource.TestCase{
		ProviderFactories: zitadelProviderFactories(ctx),
		CheckDestroy:      checkDestroy(ctx),
		Steps: []resource.TestStep{
			{ // Check first plan has a diff
				Config:             initialConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check resource is created
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					assignGetProviderByIDResponse(ctx, getProviderByIDResponse),
					resource.ComposeAggregateTestCheckFunc(
						checkStateHasIDSet(ctx),
						checkName("aninitialprovidername", getProviderByIDResponse),
					),
				),
			}, { // Check updating name has a diff
				Config:             updatedNameConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check name can be updated
				Config: updatedNameConfig,
				Check: resource.ComposeTestCheckFunc(
					assignGetProviderByIDResponse(ctx, getProviderByIDResponse),
					checkName("anupdatedprovidername", getProviderByIDResponse),
				),
			}, { // Check updating client secret has a diff
				Config:             updatedClientSecretConfig,
				ExpectNonEmptyPlan: true,
				// ExpectNonEmptyPlan just works with PlanOnly set to true
				PlanOnly: true,
			}, { // Check client secret can be updated
				Config: updatedClientSecretConfig,
			}, { // Expect import error if client secret is not given
				ResourceName:  zitadelCtx.terraformName,
				ImportState:   true,
				ImportStateId: "123:456",
				ExpectError:   regexp.MustCompile(`123:456`),
			}, { // Expect importing works
				ResourceName: zitadelCtx.terraformName,
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					lastState := state.RootModule().Resources[zitadelCtx.terraformName].Primary
					return fmt.Sprintf("%s:%s:%s", lastState.Attributes["org_id"], lastState.ID, importedSecret), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
				Check: func(state *terraform.State) error {
					// Check the client_secret is imported correctly
					currentState := state.RootModule().Resources[zitadelCtx.terraformName].Primary
					actual := currentState.Attributes["client_secret"]
					if actual != importedSecret {
						return fmt.Errorf("expected client_secret to be %s, but got %s", importedSecret, actual)
					}
					return nil
				},
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

func checkClientSecret(ctx context.Context, expectClientSecret string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		zitadelCtx := fromZitadelContext(ctx)
		currentState := state.RootModule().Resources[zitadelCtx.terraformName].Primary
		actual := currentState.Attributes["client_secret"]
		if actual != expectClientSecret {
			return fmt.Errorf("expected client_secret to be %s, but got %s", expectClientSecret, actual)
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
