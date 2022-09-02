package application_api

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an API application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			appIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "orgID of the application",
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the application",
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auth method type",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
