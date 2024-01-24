package application_saml_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/application_saml"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_application_saml", application_saml.AppIDVar, func() (string, error) {
		app, err := frame.AddSAMLApp(frame, &management.AddSAMLAppRequest{
			ProjectId: projectID,
			Name:      name,
			Metadata:  &management.AddSAMLAppRequest_MetadataXml{MetadataXml: []byte("metadata")},
		})
		return app.GetAppId(), err
	})
}
