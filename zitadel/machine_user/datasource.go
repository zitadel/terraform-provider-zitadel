package machine_user

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
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

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing machine users situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			userIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all user IDs",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			UserNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username to filter by",
			},
			userNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by username" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(userNameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			emailVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email to filter by",
			},
			emailMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by email" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(emailMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			loginNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login name to filter by",
			},
			loginNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by login name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(loginNameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
		},
		ReadContext: list,
	}
}
