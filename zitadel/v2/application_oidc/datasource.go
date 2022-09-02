package application_oidc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an OIDC application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			appIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "orgID of the application",
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the application",
			},
			redirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "RedirectURIs",
			},
			responseTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Response type",
			},
			grantTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Grant types",
			},
			appTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "App type",
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auth method type",
			},
			postLogoutRedirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Post logout redirect URIs",
			},
			versionVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version",
			},
			devModeVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Dev mode",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access token type",
			},
			accessTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Access token role assertion",
			},
			idTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ID token role assertion",
			},
			idTokenUserinfoAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Token userinfo assertion",
			},
			clockSkewVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Clockskew",
			},
			additionalOriginsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Additional origins",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
