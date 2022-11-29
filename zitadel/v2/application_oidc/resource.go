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
				Description: "Response type" + helper.DescriptionEnumValuesList(app.OIDCResponseType_name),
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
				Description: "Grant types" + helper.DescriptionEnumValuesList(app.OIDCGrantType_name),
				/* Not yet supported
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return enumValuesValidation(applicationGrantTypesVar, value, app.OIDCGrantType_value)
				},*/
			},
			appTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "App type" + helper.DescriptionEnumValuesList(app.OIDCAppType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(appTypeVar, value, app.OIDCAppType_value)
				},
				Default: app.OIDCAppType_name[0],
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auth method type" + helper.DescriptionEnumValuesList(app.OIDCAuthMethodType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(authMethodTypeVar, value, app.OIDCAuthMethodType_value)
				},
				Default: app.OIDCAuthMethodType_name[0],
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
				Description: "Version" + helper.DescriptionEnumValuesList(app.OIDCVersion_name),
				Default:     app.OIDCVersion_name[0],
			},
			devModeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dev mode",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type" + helper.DescriptionEnumValuesList(app.OIDCTokenType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(accessTokenTypeVar, value, app.OIDCTokenType_value)
				},
				Default: app.OIDCTokenType_name[0],
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
				Optional:    true,
				Description: "Clockskew",
				Default:     "0s",
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
