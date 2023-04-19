package idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an LDAP IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the IDP",
			},
			idp_utils.ServersVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Servers to try in order for establishing LDAP connections",
			},
			idp_utils.StartTLSVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Wether to use StartTLS for LDAP connections",
			},
			idp_utils.BaseDNVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base DN for LDAP connections",
			},
			idp_utils.BindDNVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bind DN for LDAP connections",
			},
			idp_utils.BindPasswordVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bind password for LDAP connections",
				Sensitive:   true,
			},
			idp_utils.UserBaseVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User base for LDAP connections",
			},
			idp_utils.UserObjectClassesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "User object classes for LDAP connections",
			},
			idp_utils.UserFiltersVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "User filters for LDAP connections",
			},
			idp_utils.TimeoutVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timeout for LDAP connections",
			},
			idp_utils.IdAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the id",
			},
			idp_utils.FirstNameAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the first name",
			},
			idp_utils.LastNameAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the last name",
			},
			idp_utils.DisplayNameAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the display name",
			},
			idp_utils.NickNameAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the nick name",
			},
			idp_utils.PreferredUsernameAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the preferred username",
			},
			idp_utils.EmailAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the email",
			},
			idp_utils.EmailVerifiedAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the email verified state",
			},
			idp_utils.PhoneAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the phone",
			},
			idp_utils.PhoneVerifiedAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the phone verified state",
			},
			idp_utils.PreferredLanguageAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the preferred language",
			},
			idp_utils.AvatarURLAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the avatar url",
			},
			idp_utils.ProfileAttributeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User attribute for the profile",
			},
			idp_utils.IsLinkingAllowedVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable if users should be able to link an existing ZITADEL user with an external account",
			},
			idp_utils.IsCreationAllowedVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable if users should be able to create a new account in ZITADEL when using an external account",
			},
			idp_utils.IsAutoCreationVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable if a new account in ZITADEL should be created automatically when login with an external account",
			},
			idp_utils.IsAutoUpdateVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable if a the ZITADEL account fields should be updated automatically on each login",
			},
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: idp_utils.Delete,
		Importer:      &schema.ResourceImporter{StateContext: idp_utils.ImportIDPWithSecret(idp_utils.BindPasswordVar)},
	}
}
