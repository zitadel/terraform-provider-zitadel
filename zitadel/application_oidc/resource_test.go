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

	resourceConfig := fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = "app_login_v1_%s"
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v1 = true
  }
}`, projectID, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string { return resourceConfig },
		"app_login_v1_"+frame.UniqueResourcesID, "app_login_v1_updated_"+frame.UniqueResourcesID,
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

func TestAccAppOIDC_LoginV2_WithBaseURI(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	resourceConfig := fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = "app_login_v2_%s"
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v2 {
      base_uri = "https://custom-login.example.com"
    }
  }
}`, projectID, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string { return resourceConfig },
		"app_login_v2_"+frame.UniqueResourcesID, "app_login_v2_updated_"+frame.UniqueResourcesID,
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

func TestAccAppOIDC_LoginV2_WithoutBaseURI(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_oidc")
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	resourceConfig := fmt.Sprintf(`
resource "zitadel_application_oidc" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = %q
  name             = "app_login_v2_default_%s"
  redirect_uris    = ["https://localhost.com/callback"]
  response_types   = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types      = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]

  login_version {
    login_v2 {}
  }
}`, projectID, frame.UniqueResourcesID)

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep},
		func(property, secret string) string { return resourceConfig },
		"app_login_v2_default_"+frame.UniqueResourcesID, "app_login_v2_default_updated_"+frame.UniqueResourcesID,
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

func checkRemoteLoginVersion(frame *test_utils.OrgTestFrame, projectId, expectedVersion, expectedBaseURI string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "zitadel_application_oidc" {
				continue
			}

			remoteResource, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: rs.Primary.ID, ProjectId: projectId})
			if err != nil {
				return err
			}

			loginVersion := remoteResource.GetApp().GetOidcConfig().GetLoginVersion()
			if loginVersion == nil {
				return fmt.Errorf("login_version is nil")
			}

			switch expectedVersion {
			case "v1":
				if loginVersion.GetLoginV1() == nil {
					return fmt.Errorf("expected LoginV1, got %T", loginVersion.GetVersion())
				}
			case "v2":
				v2 := loginVersion.GetLoginV2()
				if v2 == nil {
					return fmt.Errorf("expected LoginV2, got %T", loginVersion.GetVersion())
				}
				actualBaseURI := v2.GetBaseUri()
				if expectedBaseURI != actualBaseURI {
					return fmt.Errorf("expected base_uri %q, got %q", expectedBaseURI, actualBaseURI)
				}
			default:
				return fmt.Errorf("unknown version %s", expectedVersion)
			}
		}
		return nil
	}
}
