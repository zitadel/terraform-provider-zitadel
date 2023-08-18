package app_key_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccAppKey(t *testing.T) {
	resourceName := "zitadel_application_key"
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
	apiApp, err := frame.AddAPIApp(frame, &management.AddAPIAppRequest{
		ProjectId:      project.GetId(),
		Name:           frame.UniqueResourcesID,
		AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
	})
	resourceExample, exampleAttributes := frame.ReadExample(t, test_utils.Resources, frame.ResourceType)
	projectDatasourceExample, _ := frame.ReadExample(t, test_utils.Datasources, "project")
	projectDatasourceExample = strings.Replace(projectDatasourceExample, "123456789012345678", project.GetId(), 1)
	appDatasourceExample, _ := frame.ReadExample(t, test_utils.Datasources, "application_api")
	appDatasourceExample = strings.Replace(appDatasourceExample, "123456789012345678", apiApp.GetAppId(), 1)
	exampleProperty := test_utils.AttributeValue(t, "expiration_date", exampleAttributes).AsString()
	updatedProperty := "2501-01-01T08:45:00Z"
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf("%s\n%s\n%s\n%s", frame.OrgExampleDatasource, projectDatasourceExample, appDatasourceExample, strings.Replace(resourceExample, exampleProperty, configProperty, 1))
		},
		exampleProperty, updatedProperty,
		"", "",
		false,
		checkRemoteProperty(frame, project.GetId(), apiApp.GetAppId()),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, project.GetId(), apiApp.GetAppId()), updatedProperty),
		nil, nil, "", "",
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
