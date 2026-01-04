package instance_features

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the feature flags configuration for an instance.",
		Schema: map[string]*schema.Schema{
			loginDefaultOrgVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The login UI will use the settings of the default org (and not from the instance) if no organization context is set",
			},
			userSchemaVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "User Schemas allow to manage data schemas of user. If the flag is enabled, you'll be able to use the new API and its features. Note that it is still in an early stage.",
			},
			oidcTokenExchangeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable the experimental `urn:ietf:params:oauth:grant-type:token-exchange` grant type for the OIDC token endpoint. Token exchange can be used to request tokens with a lesser scope or impersonate other users. See the security policy to allow impersonation on an instance.",
			},
			improvedPerformanceVar: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						improvedPerformanceProjectGrant,
						improvedPerformanceProject,
						improvedPerformanceUserGrant,
						improvedPerformanceOrgDomainVerified,
					}, false),
				},
				Description: "Improves performance of specified execution paths. Possible values: IMPROVED_PERFORMANCE_PROJECT_GRANT, IMPROVED_PERFORMANCE_PROJECT, IMPROVED_PERFORMANCE_USER_GRANT, IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED",
			},
			debugOidcParentErrorVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Return parent errors to OIDC clients for debugging purposes. Parent errors may contain sensitive data or unwanted details about the system status of zitadel. Only enable if really needed.",
			},
			oidcSingleV1SessionTerminationVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the flag is enabled, you'll be able to terminate a single session from the login UI by providing an id_token with a `sid` claim as id_token_hint on the end_session endpoint. Note that currently all sessions from the same user agent (browser) are terminated in the login UI. Sessions managed through the Session API already allow the termination of single sessions.",
			},
			enableBackChannelLogoutVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the flag is enabled, you'll be able to use the OIDC Back-Channel Logout to be notified in your application about terminated user sessions.",
			},
			loginV2Var: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						loginV2RequiredVar: {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Require that all users must use the new login UI. If enabled, all users will be redirected to the login V2 regardless of the application's preference.",
						},
						loginV2BaseURIVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optionally specify a base uri of the login UI. If unspecified the default URI will be used.",
						},
					},
				},
				Description: "Specify the login UI for all users and applications regardless of their preference.",
			},
			permissionCheckV2Var: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable a newer, more performant, permission check used for v2 and v3 resource based APIs.",
			},
			consoleUseV2UserAPIVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this is enabled the console web client will use the new User v2 API for certain calls",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
