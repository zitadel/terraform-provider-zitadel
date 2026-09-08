package org_idp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all identity providers available to an organization, optionally filtered by name, type and owner.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar:         helper.OrgIDDatasourceField,
			idp_utils.NameVar:       idp_utils.NameFilterDataSourceField,
			idp_utils.NameMethodVar: idp_utils.NameMethodDataSourceField,
			idp_utils.TypeVar:       idp_utils.TypeFilterDataSourceField,
			idp_utils.OwnerTypeVar:  idp_utils.OwnerTypeFilterDataSourceField,
			idp_utils.IdpsVar:       idp_utils.IdpsDataSourceField,
		},
		ReadContext: list,
	}
}
