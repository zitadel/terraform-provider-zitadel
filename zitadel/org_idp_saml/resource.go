package org_idp_saml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a SAML IdP on the organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar:                           helper.OrgIDResourceField,
			idp_utils.NameVar:                         idp_utils.NameResourceField,
			idp_saml.BindingVar:                       idp_saml.BindingResourceField,
			idp_saml.MetadataXMLVar:                   idp_saml.MetadataXMLResourceField,
			idp_saml.MetadataURLVar:                   idp_saml.MetadataURLResourceField,
			idp_saml.WithSignedRequestVar:             idp_saml.WithSignedRequestResourceField,
			idp_saml.NameIdFormatVar:                  idp_saml.NameIdFormatResourceField,
			idp_saml.TransientMappingAttributeNameVar: idp_saml.TransientMappingAttributeNameResourceField,
			idp_saml.FederatedLogoutEnabledVar:        idp_saml.FederatedLogoutEnabledResourceField,
			idp_saml.SignatureAlgorithmVar:            idp_saml.SignatureAlgorithmResourceField,
			idp_utils.IsLinkingAllowedVar:             idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar:            idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:               idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:                 idp_utils.IsAutoUpdateResourceField,
			idp_utils.AutoLinkingVar:                  idp_utils.AutoLinkingResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: org_idp_utils.Delete,
		Importer:      helper.ImportWithIDAndOptionalOrg(idp_utils.IdpIDVar),
	}
}
