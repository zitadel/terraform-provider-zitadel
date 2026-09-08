package idp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all identity providers on the instance, optionally filtered by name and type.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar:       idp_utils.NameFilterDataSourceField,
			idp_utils.NameMethodVar: idp_utils.NameMethodDataSourceField,
			idp_utils.TypeVar:       idp_utils.TypeFilterDataSourceField,
			idp_utils.IdpsVar:       idp_utils.IdpsDataSourceField,
		},
		ReadContext: list,
	}
}
