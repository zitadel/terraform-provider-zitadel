package user_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "represents role grants",
		Schema: map[string]*schema.Schema{
			grantIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the usergrant",
			},
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
			},
			RoleKeysVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all roles for a user.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			projectNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the project.",
				Computed:    true,
			},
			roleStatusVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of role",
				Computed:    true,
			},
			userNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "username",
				Computed:    true,
			},
			emailVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "email of user",
				Computed:    true,
			},
			nameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "display name of user",
				Computed:    true,
			},
		},
		ReadContext: readDS,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "represents role grants",
		Schema: map[string]*schema.Schema{
			projectNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project.",
			},
			userGrantDataVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all usergrants.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grantIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "grantID",
						},
						UserIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "userid",
						},
						RoleKeysVar: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A set of all roles for a user.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						projectNameVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the project.",
							Computed:    true,
						},
						roleStatusVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status of role",
							Computed:    true,
						},
						userNameVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "username",
							Computed:    true,
						},
						emailVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "email of user",
							Computed:    true,
						},
						nameVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display name of user",
							Computed:    true,
						},
					},
				},
			},
		},
		ReadContext: list,
	}
}
