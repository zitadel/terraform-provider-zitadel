package human_user

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a human user situated under an organization, which then can be authorized through memberships or direct grants on other resources.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			userStateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the user",
				/* Not necessary as long as only active users are created
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return EnumValueValidation(userStateVar, value.(string), user.UserState_value)
				},*/
			},
			userNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username",
			},
			loginNamesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Loginnames",
				ForceNew:    true,
			},
			preferredLoginNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Preferred login name",
				ForceNew:    true,
			},

			firstNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "First name of the user",
			},
			lastNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Last name of the user",
			},
			nickNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Nick name of the user",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name of the user",
				Computed:    true,
			},
			preferredLanguageVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred language of the user",
				Default:     defaultPreferredLanguage,
			},
			genderVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gender of the user" + helper.DescriptionEnumValuesList(user.Gender_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(genderVar, value.(string), user.Gender_value)
				},
				Default: defaultGenderString,
			},
			emailVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email of the user",
			},
			isEmailVerifiedVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the email verified of the user, can only be true if password of the user is set",
			},
			phoneVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone of the user",
			},
			isPhoneVerifiedVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the phone verified of the user",
			},
			initialPasswordVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Initially set password for the user, not changeable after creation",
				Sensitive:   true,
				ForceNew:    true,
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		CustomizeDiff: customdiff.All(
			customdiff.IfValue(displayNameVar, func(ctx context.Context, value, meta interface{}) bool {
				if value == "" {
					return true
				}
				return false
			}, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
				return diff.SetNew(displayNameVar, defaultDisplayName(diff.Get(firstNameVar).(string), diff.Get(lastNameVar).(string)))
			}),
			customdiff.IfValue(genderVar, func(ctx context.Context, value, meta interface{}) bool {
				if value == "" {
					return true
				}
				return false
			}, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
				return diff.SetNew(genderVar, defaultGenderString)
			}),
			customdiff.IfValue(preferredLanguageVar, func(ctx context.Context, value, meta interface{}) bool {
				if value == "" {
					return true
				}
				return false
			}, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
				return diff.SetNew(preferredLanguageVar, defaultPreferredLanguage)
			}),
		),
		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
