package idp_google_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func Create(t *testing.T, frame test_utils.BaseTestFrame, client admin.AdminServiceClient) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_idp_google", idp_utils.IdpIDVar, func() (string, error) {
		i, err := client.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
			Name:         "Google " + frame.UniqueResourcesID,
			ClientId:     "dummy",
			ClientSecret: "dummy",
		})
		return i.GetId(), err
	})
}
