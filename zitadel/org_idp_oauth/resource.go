package org_idp_oauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_oauth"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a generic OAuth2 IDP on the organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar:                    helper.OrgIDResourceField,
			idp_utils.NameVar:                  idp_utils.NameResourceField,
			idp_utils.ClientIDVar:              idp_utils.ClientIDResourceField,
			idp_utils.ClientSecretVar:          idp_utils.ClientSecretResourceField,
			idp_oauth.AuthorizationEndpointVar: idp_oauth.AuthorizationEndpointResourceField,
			idp_oauth.TokenEndpointVar:         idp_oauth.TokenEndpointResourceField,
			idp_oauth.UserEndpointVar:          idp_oauth.UserEndpointResourceField,
			idp_oauth.IdAttributeVar:           idp_oauth.IdAttributeResourceField,
			idp_utils.ScopesVar:                idp_utils.ScopesResourceField,
			idp_utils.IsLinkingAllowedVar:      idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar:     idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:        idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:          idp_utils.IsAutoUpdateResourceField,
			idp_utils.AutoLinkingVar:           idp_utils.AutoLinkingResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: org_idp_utils.Delete,
		Importer:      helper.ImportWithIDAndOptionalOrgAndSecret(idp_utils.IdpIDVar, idp_utils.ClientSecretVar),
	}
}
