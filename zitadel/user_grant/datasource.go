package user_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the authorization given to a user directly, including the given roles.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			grantIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user grant.",
			},
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user.",
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			projectGrantIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the granted project.",
			},
			RoleKeysVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of roles granted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: get,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all authorizations given to a user directly, which can be looked up in detail with the user grant datasource.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			grantIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all user grant IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user.",
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the project.",
			},
			projectGrantIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the granted project.",
			},
			roleKeyVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key of a granted role.",
			},
		},
		ReadContext: list,
	}
}
