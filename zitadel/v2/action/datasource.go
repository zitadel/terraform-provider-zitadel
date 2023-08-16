package action

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an action belonging to an organization.",
		Schema: map[string]*schema.Schema{
			helper.ResourceIDVar: helper.ResourceIDDatasourceField,
			helper.OrgIDVar:      helper.OrgIDDatasourceField,
			stateVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the state of the action",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			scriptVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			timeoutVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "after which time the action will be terminated if not finished",
			},
			allowedToFailVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "when true, the next action will be called even if this action fails",
			},
		},
		ReadContext: read,
	}
}
