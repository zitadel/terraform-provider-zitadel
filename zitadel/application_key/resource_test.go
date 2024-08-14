package application_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_api/application_api_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_key"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project/project_test_dep"
)

func TestAccAppKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_application_key")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, application_key.ExpirationDateVar, exampleAttributes).AsString()
	updatedProperty := "2501-01-01T08:45:00Z"
	projectDep, projectID := project_test_dep.Create(t, frame, frame.UniqueResourcesID)
	appDep, appID, _ := application_api_test_dep.Create(t, frame, projectID, frame.UniqueResourcesID)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, projectDep, appDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, projectID, appID),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, projectID, appID), updatedProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_key.ProjectIDVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_key.AppIDVar),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, application_key.KeyDetailsVar),
		),
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
