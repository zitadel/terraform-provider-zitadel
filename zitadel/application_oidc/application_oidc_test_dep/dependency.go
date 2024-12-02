package application_oidc_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) (template string, id string, clientId string) {
	template, id = test_utils.CreateDefaultDependency(t, "zitadel_application_oidc", application_oidc.AppIDVar, func() (string, error) {
		oidcApp, err := frame.AddOIDCApp(frame, &management.AddOIDCAppRequest{
			ProjectId: projectID,
			Name:      name,
		})
		clientId = oidcApp.GetClientId()
		return oidcApp.GetAppId(), err
	})
	return template, id, clientId
}
