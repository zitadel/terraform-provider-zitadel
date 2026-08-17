package application_oidc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

func TestAccAppOIDC(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, application_oidc.NameVar, exampleAttributes).AsString()
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedproperty",
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ProjectIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientSecretVar),
		),
	)
}

func TestAccAppOIDC_LoginV1(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string {
			return fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = %q
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v1 = true
  }
}`, projectID, property)
		},
		"app_login_v1_"+frame.UniqueResourcesID,
		"app_login_v1_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ProjectIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientSecretVar),
		),
		"compliance_problems",
	)
}

func TestAccAppOIDC_LoginV2_WithBaseURI(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string {
			return fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = %q
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v2 {
      base_uri = "https://custom-login.example.com"
    }
  }
}`, projectID, property)
		},
		"app_login_v2_"+frame.UniqueResourcesID,
		"app_login_v2_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ProjectIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientSecretVar),
		),
		"compliance_problems",
	)
}

func TestAccAppOIDC_LoginV2_WithoutBaseURI(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string {
			return fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = %q
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v2 {}
  }
}`, projectID, property)
		},
		"app_login_v2_default_"+frame.UniqueResourcesID,
		"app_login_v2_default_updated_"+frame.UniqueResourcesID,
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ProjectIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_oidc.ClientSecretVar),
		),
		"compliance_problems", "login_version",
	)
}

func TestAccAppOIDCConfigUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_oidc" "default" {
  org_id                      = data.zitadel_org.default.id
  project_id                  = "%s"
  name                        = "%s"
  redirect_uris               = ["https://localhost.com/callback"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  dev_mode                    = false
  access_token_role_assertion = false
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_oidc" "default" {
  org_id                      = data.zitadel_org.default.id
  project_id                  = "%s"
  name                        = "%s"
  redirect_uris               = ["https://localhost.com/callback", "https://localhost.com/callback2"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  dev_mode                    = true
  access_token_role_assertion = true
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "dev_mode", "false"),
					resource.TestCheckResourceAttr(frame.TerraformName, "access_token_role_assertion", "false"),
					resource.TestCheckResourceAttr(frame.TerraformName, "redirect_uris.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "dev_mode", "true"),
					resource.TestCheckResourceAttr(frame.TerraformName, "access_token_role_assertion", "true"),
					resource.TestCheckResourceAttr(frame.TerraformName, "redirect_uris.#", "2"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectId string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: frame.State(state).ID, ProjectId: projectId})
			if err != nil {
				return err
			}
			actual := remoteResource.GetApp().GetName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}

func TestAccAppOIDC_NativeAppLinks(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	const (
		teamID      = "ABCDE12345"
		bundleID    = "com.example.app"
		packageName = "com.example.app"
		fingerprint = "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"
	)

	config := func(links string) string {
		return fmt.Sprintf(`
%s
%s
resource "zitadel_application_oidc" "default" {
  org_id         = data.zitadel_org.default.id
  project_id     = "%s"
  name           = "%s"
  redirect_uris  = ["https://localhost.com/callback"]
  response_types = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types    = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
%s
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID, links)
	}

	withLinks := config(fmt.Sprintf(`
  ios {
    team_id   = %q
    bundle_id = %q
  }
  android {
    package_name             = %q
    sha256_cert_fingerprints = [%q]
  }
`, teamID, bundleID, packageName, fingerprint))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: withLinks,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "ios.0.team_id", teamID),
					resource.TestCheckResourceAttr(frame.TerraformName, "ios.0.bundle_id", bundleID),
					resource.TestCheckResourceAttr(frame.TerraformName, "android.0.package_name", packageName),
					resource.TestCheckResourceAttr(frame.TerraformName, "android.0.sha256_cert_fingerprints.#", "1"),
					resource.TestCheckResourceAttr(frame.TerraformName, "android.0.sha256_cert_fingerprints.0", fingerprint),
					checkRemoteAppLinks(frame, projectID, teamID, bundleID, packageName, []string{fingerprint}),
				),
			},
			{
				// removing the blocks must clear the config server-side, not leave it stale
				Config: config(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "ios.#", "0"),
					resource.TestCheckResourceAttr(frame.TerraformName, "android.#", "0"),
					checkRemoteAppLinks(frame, projectID, "", "", "", nil),
				),
			},
		},
	})
}

func checkRemoteAppLinks(frame *test_utils.OrgTestFrame, projectID, teamID, bundleID, packageName string, fingerprints []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		remote, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: frame.State(state).ID, ProjectId: projectID})
		if err != nil {
			return err
		}
		oidc := remote.GetApp().GetOidcConfig()
		if got := oidc.GetIos().GetTeamId(); got != teamID {
			return fmt.Errorf("expected ios team_id %q, but got %q", teamID, got)
		}
		if got := oidc.GetIos().GetBundleId(); got != bundleID {
			return fmt.Errorf("expected ios bundle_id %q, but got %q", bundleID, got)
		}
		if got := oidc.GetAndroid().GetPackageName(); got != packageName {
			return fmt.Errorf("expected android package_name %q, but got %q", packageName, got)
		}
		got := oidc.GetAndroid().GetSha256CertFingerprints()
		if len(got) != len(fingerprints) {
			return fmt.Errorf("expected %d android fingerprints, but got %d", len(fingerprints), len(got))
		}
		for i := range fingerprints {
			if got[i] != fingerprints[i] {
				return fmt.Errorf("expected android fingerprint %q, but got %q", fingerprints[i], got[i])
			}
		}
		return nil
	}
}
