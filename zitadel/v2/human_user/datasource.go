package human_user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a human user situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
		Schema: map[string]*schema.Schema{
			userIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			userStateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the user",
			},
			userNameVar: {
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
			displayNameVar: {
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
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
