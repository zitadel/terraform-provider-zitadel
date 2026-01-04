package idp_apple

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an Apple IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar:              idp_utils.NameResourceField,
			idp_utils.ClientIDVar:          idp_utils.ClientIDResourceField,
			TeamIDVar:                      TeamIDResourceField,
			KeyIDVar:                       KeyIDResourceField,
			PrivateKeyVar:                  PrivateKeyResourceField,
			idp_utils.ScopesVar:            idp_utils.ScopesResourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateResourceField,
			idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: idp_utils.Delete,
		Importer:      helper.ImportWithIDAndOptionalSecret(idp_utils.IdpIDVar, PrivateKeyVar),
	}
}
