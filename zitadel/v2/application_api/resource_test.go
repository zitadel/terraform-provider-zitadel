package application_api_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccAppAPI(t *testing.T) {
	resourceName := "zitadel_application_api"
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
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		frame.OrgExampleDatasource,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "my_%s" {
  org_id           = "%s"
  project_id       = "%s"
  name             = "%s"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}`, resourceName, resourceName, frame.OrgID, project.GetId(), configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		false,
		checkRemoteProperty(frame, project.GetId()),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, project.GetId()), updatedProperty),
		nil, nil, "", "",
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
