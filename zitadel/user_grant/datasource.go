package user_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the authorizations given to a user directly, including the given roles.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the project to filter user grants by",
			},
			projectGrantIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the granted project to filter user grants by",
			},
			roleKeyVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role key to filter user grants by",
			},
			userGrantsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user grants",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						idVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the user grant",
						},
						projectIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the project",
						},
						projectNameVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the project",
						},
						projectGrantIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the granted project",
						},
						grantedOrgIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the organization the project is granted to",
						},
						RoleKeysVar: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of roles granted",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						stateVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the user grant",
						},
					},
				},
			},
		},
		ReadContext: list,
	}
}
