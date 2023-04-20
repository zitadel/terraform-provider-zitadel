package idp_utils

const (
	IdpIDVar                 = "id"
	NameVar                  = "name"
	ClientIDVar              = "client_id"
	ClientSecretVar          = "client_secret"
	ScopesVar                = "scopes"
	IsLinkingAllowedVar      = "is_linking_allowed"
	IsCreationAllowedVar     = "is_creation_allowed"
	IsAutoCreationVar        = "is_auto_creation"
	IsAutoUpdateVar          = "is_auto_update"
	AuthorizationEndpointVar = "authorization_endpoint"
	TokenEndpointVar         = "token_endpoint"
	UserEndpointVar          = "user_endpoint"
	IssuerVar                = "issuer"
	TenantTypeVar            = "tenant_type"
	TenantIDVar              = "tenant_id"
	EmailVerifiedVar         = "email_verified"
	// ServersVar is the first LDAP specific provider config property
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
	// FirstNameAttributeVar is the first LDAP specific user config property
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