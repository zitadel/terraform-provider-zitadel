package app_key_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/app_key"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccAppKey(t *testing.T) {
	resourceName := "zitadel_application_key"
	initialProperty := "2500-01-01T08:45:00Z"
	updatedProperty := "2501-01-01T08:45:00Z"
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
	app, err := frame.AddOIDCApp(frame, &management.AddOIDCAppRequest{
		ProjectId:      project.GetId(),
		Name:           frame.UniqueResourcesID,
		AuthMethodType: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
	})
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id          = "%s"
  project_id      = "%s"
  app_id          = "%s"
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "%s"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, project.GetId(), app.GetAppId(), configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, project.GetId(), app.GetAppId()),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, project.GetId(), app.GetAppId()), updatedProperty),
		test_utils.ImportStateId(frame.BaseTestFrame, helper.OrgIDVar, app_key.ProjectIDVar, app_key.AppIDVar, app_key.KeyDetailsVar),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, projectId, appId string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetAppKey(frame, &management.GetAppKeyRequest{KeyId: frame.State(state).ID, ProjectId: projectId, AppId: appId})
			if err != nil {
				return err
			}
			actual := remoteResource.GetKey().GetExpirationDate().AsTime().Format("2006-01-02T15:04:05Z")
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
