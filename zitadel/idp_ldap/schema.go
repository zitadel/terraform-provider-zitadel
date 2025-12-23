package idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	ServersVar           = "servers"
	StartTLSVar          = "start_tls"
	BaseDNVar            = "base_dn"
	BindDNVar            = "bind_dn"
	BindPasswordVar      = "bind_password"
	UserBaseVar          = "user_base"
	UserObjectClassesVar = "user_object_classes"
	UserFiltersVar       = "user_filters"
	TimeoutVar           = "timeout"
	IdAttributeVar       = "id_attribute"

	FirstNameAttributeVar         = "first_name_attribute"
	LastNameAttributeVar          = "last_name_attribute"
	DisplayNameAttributeVar       = "display_name_attribute"
	NickNameAttributeVar          = "nick_name_attribute"
	PreferredUsernameAttributeVar = "preferred_username_attribute"
	EmailAttributeVar             = "email_attribute"
	EmailVerifiedAttributeVar     = "email_verified_attribute"
	PhoneAttributeVar             = "phone_attribute"
	PhoneVerifiedAttributeVar     = "phone_verified_attribute"
	PreferredLanguageAttributeVar = "preferred_language_attribute"
	AvatarURLAttributeVar         = "avatar_url_attribute"
	ProfileAttributeVar           = "profile_attribute"
)

var (
	ServersResourceField = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "Servers to try in order for establishing LDAP connections",
	}
	ServersDataSourceField = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "Servers to try in order for establishing LDAP connections",
	}
	StartTLSResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Wether to use StartTLS for LDAP connections",
	}
	StartTLSDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Wether to use StartTLS for LDAP connections",
	}
	BaseDNResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Base DN for LDAP connections",
	}
	BaseDNDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Base DN for LDAP connections",
	}
	BindDNResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Bind DN for LDAP connections",
	}
	BindDNDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Bind DN for LDAP connections",
	}
	BindPasswordResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Bind password for LDAP connections",
		Sensitive:   true,
	}
	BindPasswordDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Bind password for LDAP connections",
		Sensitive:   true,
	}
	UserBaseResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "User base for LDAP connections",
	}
	UserBaseDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User base for LDAP connections",
	}
	UserObjectClassesResourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "User object classes for LDAP connections",
	}
	UserObjectClassesDataSourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "User object classes for LDAP connections",
	}
	UserFiltersResourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "User filters for LDAP connections",
	}
	UserFiltersDataSourceField = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "User filters for LDAP connections",
	}
	TimeoutResourceField = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Timeout for LDAP connections",
		DiffSuppressFunc: helper.DurationDiffSuppress,
	}
	TimeoutDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Timeout for LDAP connections",
	}
	IdAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the id",
	}
	IdAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the id",
	}

	FirstNameAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the first name",
	}
	FirstNameAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the first name",
	}
	LastNameAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the last name",
	}
	LastNameAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the last name",
	}
	DisplayNameAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the display name",
	}
	DisplayNameAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the display name",
	}
	NickNameAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the nick name",
	}
	NickNameAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the nick name",
	}
	PreferredUsernameAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the preferred username",
	}
	PreferredUsernameAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the preferred username",
	}
	EmailAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the email",
	}
	EmailAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the email",
	}
	EmailVerifiedAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the email verified state",
	}
	EmailVerifiedAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the email verified state",
	}
	PhoneAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the phone",
	}
	PhoneAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the phone",
	}
	PhoneVerifiedAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the phone verified state",
	}
	PhoneVerifiedAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the phone verified state",
	}
	PreferredLanguageAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the preferred language",
	}
	PreferredLanguageAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the preferred language",
	}
	AvatarURLAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the avatar url",
	}
	AvatarURLAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the avatar url",
	}
	ProfileAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "User attribute for the profile",
	}
	ProfileAttributeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User attribute for the profile",
	}
)
