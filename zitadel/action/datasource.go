package action

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an action belonging to an organization.",
		Schema: map[string]*schema.Schema{
			ActionIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: helper.OrgIDResourceField,
			stateVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the state of the action",
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			ScriptVar: {
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
