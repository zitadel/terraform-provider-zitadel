package login_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom login policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			allowUsernamePasswordVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user is allowed to login with his username and password",
			},
			allowRegisterVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a person is allowed to register a user on this organisation",
			},
			allowExternalIDPVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user is allowed to add a defined identity provider. E.g. Google auth",
			},
			forceMFAVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user MUST use a multi factor to log in",
			},
			forceMFALocalOnlyVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "if activated, ZITADEL only enforces MFA on local authentications. On authentications through MFA, ZITADEL won't prompt for MFA.",
			},
			passwordlessTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "defines if passwordless is allowed for users",
			},
			hidePasswordResetVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if password reset link should be shown in the login screen",
			},
			ignoreUnknownUsernamesVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if unknown username on login screen directly return an error or always display the password screen",
			},
			DefaultRedirectURIVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "defines where the user will be redirected to if the login is started without app context (e.g. from mail)",
			},
			passwordCheckLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			externalLoginCheckLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			mfaInitSkipLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			secondFactorCheckLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			multiFactorCheckLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			secondFactorsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "allowed second factors",
			},
			multiFactorsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "allowed multi factors",
			},
			idpsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "allowed idps to login or register",
			},
			allowDomainDiscovery: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "if set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organisation on success.",
			},
			disableLoginWithEmail: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "defines if user can additionally (to the loginname) be identified by their verified email address",
			},
			disableLoginWithPhone: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "defines if user can additionally (to the loginname) be identified by their verified phone number",
			},
		},
		CreateContext: create,
		UpdateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
