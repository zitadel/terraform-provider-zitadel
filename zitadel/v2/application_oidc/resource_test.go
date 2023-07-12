package application_oidc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccZITADELAppOIDC(t *testing.T) {
	resourceName := "zitadel_application_oidc"
	initialProperty := "initialname"
	updatedProperty := "updatedname"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	project, err := frame.AddProject(frame, &management.AddProjectRequest{
		Name: frame.UniqueResourcesID,
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(configProperty, _ interface{}) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id           = "%s"
  project_id       = "%s"
  name             = "%s"
  redirect_uris               = ["https://localhost.com"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris   = ["https://localhost.com"]
  app_type                    = "OIDC_APP_TYPE_WEB"
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_BASIC"
  version                     = "OIDC_VERSION_1_0"
  clock_skew                  = "0s"
  dev_mode                    = true
  access_token_type           = "OIDC_TOKEN_TYPE_BEARER"
  access_token_role_assertion = false
  id_token_role_assertion     = false
  id_token_userinfo_assertion = false
  additional_origins          = []
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, project.GetId(), configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(frame, project.GetId()),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, project.GetId()), updatedProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectId string) func(interface{}) resource.TestCheckFunc {
	return func(expect interface{}) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs := state.RootModule().Resources[frame.TerraformName]
			remoteResource, err := frame.GetAppByID(frame, &management.GetAppByIDRequest{AppId: rs.Primary.ID, ProjectId: projectId})
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
