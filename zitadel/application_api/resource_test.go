package application_api_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_api"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project/project_test_dep"
)

func TestAccAppAPI(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_api")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, application_api.NameVar, exampleAttributes).AsString()
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
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_api.ProjectIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_api.ClientIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_api.ClientSecretVar),
		),
	)
}

func TestAccAppAPIAuthMethodTypeUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_api")
	_, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_api" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = "%s"
  name             = "%s"
  auth_method_type = "API_AUTH_METHOD_TYPE_BASIC"
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_application_api" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = "%s"
  name             = "%s"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, projectID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "auth_method_type", "API_AUTH_METHOD_TYPE_BASIC"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "auth_method_type", "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"),
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
