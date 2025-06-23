package action_target

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a target, which can be used in executions.",
		Schema: map[string]*schema.Schema{
			TargetIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: helper.OrgIDResourceField,
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the target.",
			},
			EndpointVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint of the target.",
			},
			TargetTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the target.",
			},
			TimeoutVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timeout defines the duration until ZITADEL cancels the execution.",
			},
			InterruptOnErrorVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Define if any error stops the whole execution.",
			},
		},
		ReadContext: read,
	}
}
