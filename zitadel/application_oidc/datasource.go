package application_oidc

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an OIDC application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			AppIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			NameVar: {
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
			ClientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client ID",
				Sensitive:   true,
			},
			skipNativeAppSuccessPageVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Skip the successful login page on native apps and directly redirect the user to the callback.",
			},
			NoneCompliantVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "specifies whether the config is OIDC compliant. A production configuration SHOULD be compliant",
			},
			ComplianceProblemsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "lists the problems for non-compliancy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ComplianceProblemKeyVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine-readable identifier for the compliance problem",
						},
						ComplianceProblemMessageVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable localized message",
						},
					},
				},
			},
			BackChannelLogoutURIVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ZITADEL will use this URI to notify the application about terminated session according to the OIDC Back-Channel Logout",
			},
			LoginVersionVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specify the preferred login UI, where the user is redirected to for authentication. If unset, the login UI is chosen by the instance default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						LoginV1Var: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Login V1",
						},
						LoginV2Var: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Login V2",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									BaseURIVar: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Optionally specify a base uri of the login UI. If unspecified the default URI will be used.",
									},
								},
							},
						},
					},
				},
			},
		},
		ReadContext: read,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing multiple OIDC applications belonging to a project.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			appIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			nameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying applications by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
		},
		ReadContext: list,
	}
}
