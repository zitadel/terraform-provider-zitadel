package application_v2

const (
	AppIDVar      = "app_id"
	appIDsVar     = "app_ids"
	ProjectIDVar  = "project_id"
	NameVar       = "name"
	nameMethodVar = "name_method"
	stateVar      = "state"

	// Top-level config block selectors (ExactlyOneOf).
	oidcBlockVar = "oidc"
	samlBlockVar = "saml"
	apiBlockVar  = "api"

	// OIDC config keys (mirror the v1 application_oidc surface).
	redirectURIsVar             = "redirect_uris"
	responseTypesVar            = "response_types"
	grantTypesVar               = "grant_types"
	appTypeVar                  = "app_type"
	authMethodTypeVar           = "auth_method_type"
	postLogoutRedirectURIsVar   = "post_logout_redirect_uris"
	versionVar                  = "version"
	devModeVar                  = "dev_mode"
	accessTokenTypeVar          = "access_token_type"
	accessTokenRoleAssertionVar = "access_token_role_assertion"
	idTokenRoleAssertionVar     = "id_token_role_assertion"
	idTokenUserinfoAssertionVar = "id_token_userinfo_assertion"
	clockSkewVar                = "clock_skew"
	additionalOriginsVar        = "additional_origins"
	skipNativeAppSuccessPageVar = "skip_native_app_success_page"
	backChannelLogoutURIVar     = "back_channel_logout_uri"
	loginVersionVar             = "login_version"
	loginV1Var                  = "login_v1"
	loginV2Var                  = "login_v2"
	baseURIVar                  = "base_uri"

	// Computed OIDC outputs.
	clientIDVar           = "client_id"
	clientSecretVar       = "client_secret"
	noneCompliantVar      = "none_compliant"
	complianceProblemsVar = "compliance_problems"
	complianceKeyVar      = "key"
	complianceMessageVar  = "message"

	// SAML config keys.
	metadataXMLVar = "metadata_xml"
	metadataURLVar = "metadata_url"
)
