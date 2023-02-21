package default_login_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default login policy.",
		Schema: map[string]*schema.Schema{
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
			defaultRedirectURIVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "defines where the user will be redirected to if the login is started without app context (e.g. from mail)",
			},
			passwordCheckLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			externalLoginCheckLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			mfaInitSkipLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			secondFactorCheckLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			multiFactorCheckLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
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
		CreateContext: update,
		UpdateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
