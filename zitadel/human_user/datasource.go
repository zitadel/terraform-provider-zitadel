package human_user

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a human user situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
		Schema: map[string]*schema.Schema{
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
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
			firstNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "First name of the user",
			},
			lastNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last name of the user",
			},
			nickNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Nick name of the user",
			},
			DisplayNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the user",
			},
			preferredLanguageVar: {
				Type:        schema.TypeString,
				Description: "Preferred language of the user",
				Computed:    true,
			},
			genderVar: {
				Type:        schema.TypeString,
				Description: "Gender of the user",
				Computed:    true,
			},
			emailVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user",
			},
			isEmailVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the email verified of the user, can only be true if password of the user is set",
			},
			phoneVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Phone of the user",
			},
			isPhoneVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the phone verified of the user",
			},
		},
		ReadContext: readFunc(true),
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing human users situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
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
			firstNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name to filter by",
			},
			firstNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by first name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(firstNameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			lastNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Last name to filter by",
			},
			lastNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by last name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(lastNameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			nickNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Nick name to filter by",
			},
			nickNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by nick name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nickNameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			DisplayNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name to filter by",
			},
			displayNameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying users by display name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(displayNameMethodVar, value, object.TextQueryMethod_value)
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
