package application_saml_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, name string) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_application_saml", application_saml.AppIDVar, func() (string, error) {
		app, err := frame.AddSAMLApp(frame, &management.AddSAMLAppRequest{
			ProjectId: projectID,
			Name:      name,
			Metadata:  &management.AddSAMLAppRequest_MetadataXml{MetadataXml: metadata(name)},
		})
		return app.GetAppId(), err
	})
}

func metadata(name string) []byte {
	return []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2024-01-26T17:48:38Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"" + name + "\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"http://example.com/saml/cas\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")
}
