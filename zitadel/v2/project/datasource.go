package project

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the project, which can then be granted to different organizations or users directly, containing different applications.",
		Schema: map[string]*schema.Schema{
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the project",
			},
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization in which the project is located",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the project",
			},
			roleAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "describes if roles of user should be added in token",
			},
			roleCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the user has at least one on this project",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			privateLabelingSettingVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines from where the private labeling should be triggered",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
