package org_idp_gitlab_self_hosted

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_gitlab_self_hosted"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a GitLab Self Hosted IdP on the organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar:                  helper.OrgIDResourceField,
			idp_utils.NameVar:                idp_utils.NameResourceField,
			idp_utils.ClientIDVar:            idp_utils.ClientIDResourceField,
			idp_utils.ClientSecretVar:        idp_utils.ClientSecretResourceField,
			idp_utils.ScopesVar:              idp_utils.ScopesResourceField,
			idp_utils.IsLinkingAllowedVar:    idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar:   idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:      idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:        idp_utils.IsAutoUpdateResourceField,
			idp_gitlab_self_hosted.IssuerVar: idp_gitlab_self_hosted.IssuerResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: org_idp_utils.Delete,
		Importer:      &schema.ResourceImporter{StateContext: org_idp_utils.ImportIDPWithOrgAndSecret(idp_utils.ClientSecretVar)},
	}
}
