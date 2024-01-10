package application_api_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_api"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) (string, string) {
	return test_utils.CreateProjectDefaultDependency(t, "zitadel_application_api", frame.OrgID, application_api.ProjectIDVar, projectID, application_api.AppIDVar, func() (string, error) {
		apiApp, err := frame.AddAPIApp(frame, &management.AddAPIAppRequest{
			ProjectId:      projectID,
			Name:           name,
			AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
		})
		return apiApp.GetAppId(), err
	})
}
