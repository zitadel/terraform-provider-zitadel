package idp_saml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a SAML IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar:              idp_utils.NameResourceField,
			BindingVar:                     BindingResourceField,
			WithSignedRequestVar:           WithSignedRequestResourceField,
			MetadataXMLVar:                 MetadataXMLResourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: idp_utils.Delete,
		Importer:      helper.ImportWithID(idp_utils.IdpIDVar),
	}
}
