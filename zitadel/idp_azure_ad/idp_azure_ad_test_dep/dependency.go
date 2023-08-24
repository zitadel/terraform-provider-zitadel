package idp_azure_ad_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
)

func Create(t *testing.T, frame test_utils.BaseTestFrame, client admin.AdminServiceClient) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_idp_azure_ad", idp_utils.IdpIDVar, func() (string, error) {
		i, err := client.AddAzureADProvider(frame, &admin.AddAzureADProviderRequest{
			Name: "Azure AD " + frame.UniqueResourcesID,
			Tenant: &idp.AzureADTenant{
				Type: &idp.AzureADTenant_TenantType{
					TenantType: idp.AzureADTenantType_AZURE_AD_TENANT_TYPE_COMMON,
				},
			},
			ClientId:     "dummy",
			ClientSecret: "dummy",
		})
		return i.GetId(), err
	})
}
