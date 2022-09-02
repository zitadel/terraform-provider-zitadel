package project_role

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the project roles, which can be given as authorizations to users.",
		Schema: map[string]*schema.Schema{
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			keyVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used for project role",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name used for project role",
			},
			groupVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group used for project role",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
