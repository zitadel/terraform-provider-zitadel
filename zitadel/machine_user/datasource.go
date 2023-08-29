package machine_user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a serviceaccount situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			userStateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the user",
			},
			UserNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username",
			},
			loginNamesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Loginnames",
			},
			preferredLoginNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Preferred login name",
			},

			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the machine user",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the user",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access token type",
			}},
		ReadContext: read,
	}
}
