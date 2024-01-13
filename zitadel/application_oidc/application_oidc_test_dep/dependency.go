package application_oidc_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_oidc"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_application_oidc", application_oidc.AppIDVar, func() (string, error) {
		oidcApp, err := frame.AddOIDCApp(frame, &management.AddOIDCAppRequest{
			ProjectId: projectID,
			Name:      name,
		})
		return oidcApp.GetAppId(), err
	})
}
