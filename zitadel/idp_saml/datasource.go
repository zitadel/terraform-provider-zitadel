package idp_saml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a SAML IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:             idp_utils.IdPIDDataSourceField,
			idp_utils.NameVar:              idp_utils.NameDataSourceField,
			BindingVar:                     BindingDatasourceField,
			MetadataXMLVar:                 MetadataXMLDatasourceField,
			WithSignedRequestVar:           WithSignedRequestDatasourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateDataSourceField,
			idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingDataSourceField,
		},
		ReadContext: read,
	}
}
