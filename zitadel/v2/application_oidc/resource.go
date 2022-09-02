package application_oidc

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an OIDC application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "orgID of the application",
				ForceNew:    true,
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			redirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "RedirectURIs",
			},
			responseTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Response type",
				/* Not yet supported
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return enumValuesValidation(applicationAuthMethodTypeVar, value, app.OIDCResponseType_value)
				},*/
			},
			grantTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Grant types",
				/* Not yet supported
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return enumValuesValidation(applicationGrantTypesVar, value, app.OIDCGrantType_value)
				},*/
			},
			appTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App type",
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(appTypeVar, value, app.OIDCAppType_value)
				},
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth method type",
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(authMethodTypeVar, value, app.OIDCAuthMethodType_value)
				},
			},
			postLogoutRedirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Post logout redirect URIs",
			},
			versionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version",
			},
			devModeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dev mode",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type",
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(accessTokenTypeVar, value, app.OIDCTokenType_value)
				},
			},
			accessTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Access token role assertion",
			},
			idTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ID token role assertion",
			},
			idTokenUserinfoAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Token userinfo assertion",
			},
			clockSkewVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Clockskew",
			},
			additionalOriginsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Additional origins",
			},
			clientID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated ID for this config",
				Sensitive:   true,
			},
			clientSecret: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated secret for this config",
				Sensitive:   true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
