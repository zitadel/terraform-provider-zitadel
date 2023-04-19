package idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an LDAP IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			idp_utils.NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the IDP",
			},
			idp_utils.ServersVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Servers to try in order for establishing LDAP connections",
			},
			idp_utils.StartTLSVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Wether to use StartTLS for LDAP connections",
			},
			idp_utils.BaseDNVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base DN for LDAP connections",
			},
			idp_utils.BindDNVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bind DN for LDAP connections",
			},
			idp_utils.BindPasswordVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bind password for LDAP connections",
				Sensitive:   true,
			},
			idp_utils.UserBaseVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User base for LDAP connections",
			},
			idp_utils.UserObjectClassesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "User object classes for LDAP connections",
			},
			idp_utils.UserFiltersVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "User filters for LDAP connections",
			},
			idp_utils.TimeoutVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timeout for LDAP connections",
			},
			idp_utils.IdAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the id",
			},
			idp_utils.FirstNameAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the first name",
			},
			idp_utils.LastNameAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the last name",
			},
			idp_utils.DisplayNameAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the display name",
			},
			idp_utils.NickNameAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the nick name",
			},
			idp_utils.PreferredUsernameAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the preferred username",
			},
			idp_utils.EmailAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the email",
			},
			idp_utils.EmailVerifiedAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the email verified state",
			},
			idp_utils.PhoneAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the phone",
			},
			idp_utils.PhoneVerifiedAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the phone verified state",
			},
			idp_utils.PreferredLanguageAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the preferred language",
			},
			idp_utils.AvatarURLAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the avatar url",
			},
			idp_utils.ProfileAttributeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User attribute for the profile",
			},
			idp_utils.IsLinkingAllowedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if users should be able to link an existing ZITADEL user with an external account",
			},
			idp_utils.IsCreationAllowedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if users should be able to create a new account in ZITADEL when using an external account",
			},
			idp_utils.IsAutoCreationVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if a new account in ZITADEL should be created automatically when login with an external account",
			},
			idp_utils.IsAutoUpdateVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if a the ZITADEL account fields should be updated automatically on each login",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: idp_utils.ImportIDPWithSecret(idp_utils.ClientSecretVar)},
	}
}
