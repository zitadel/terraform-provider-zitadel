package idp_apple_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func Create(t *testing.T, frame test_utils.BaseTestFrame, client admin.AdminServiceClient) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_idp_apple", idp_utils.IdpIDVar, func() (string, error) {
		// Create a minimal Apple IDP for testing
		// Note: In real scenarios, you would use actual Apple credentials
		i, err := client.AddAppleProvider(frame, &admin.AddAppleProviderRequest{
			Name:       "Apple " + frame.UniqueResourcesID,
			ClientId:   "com.example.test",
			TeamId:     "TEST123456",
			KeyId:      "KEY1234567",
			PrivateKey: []byte("dummy-private-key-for-testing"),
		})
		return i.GetId(), err
	})
}
