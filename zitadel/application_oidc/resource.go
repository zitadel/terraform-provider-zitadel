package application_oidc

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an OIDC application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			NameVar: {
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
				Description: "Response type" + helper.DescriptionEnumValuesList(apppb.OIDCResponseType_name),
			},
			grantTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Grant types" + helper.DescriptionEnumValuesList(apppb.OIDCGrantType_name),
			},
			appTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "App type" + helper.DescriptionEnumValuesList(apppb.OIDCApplicationType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(appTypeVar, value, apppb.OIDCApplicationType_value)
				},
				Default: apppb.OIDCApplicationType_name[0],
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auth method type" + helper.DescriptionEnumValuesList(apppb.OIDCAuthMethodType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(authMethodTypeVar, value, apppb.OIDCAuthMethodType_value)
				},
				Default: apppb.OIDCAuthMethodType_name[0],
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
				Description: "Version" + helper.DescriptionEnumValuesList(apppb.OIDCVersion_name),
				Default:     apppb.OIDCVersion_name[0],
				ForceNew:    true,
			},
			devModeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dev mode",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type" + helper.DescriptionEnumValuesList(apppb.OIDCTokenType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(accessTokenTypeVar, value, apppb.OIDCTokenType_value)
				},
				Default: apppb.OIDCTokenType_name[0],
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Clockskew",
				Default:          "0s",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			additionalOriginsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Additional origins",
			},
			ClientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated ID for this config",
				Sensitive:   true,
			},
			ClientSecretVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated secret for this config",
				Sensitive:   true,
			},
			skipNativeAppSuccessPageVar: {
				Type:        schema.TypeBool,
				Optional:    true,
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
				Optional:    true,
				Computed:    true,
				Description: "ZITADEL will use this URI to notify the application about terminated session according to the OIDC Back-Channel Logout",
			},
			LoginVersionVar: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Specify the preferred login UI, where the user is redirected to for authentication. If unset, the login UI is chosen by the instance default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						LoginV1Var: {
							Type:          schema.TypeBool,
							Optional:      true,
							Description:   "Login V1",
							ConflictsWith: []string{LoginVersionVar + ".0." + LoginV2Var},
						},
						LoginV2Var: {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Description:   "Login V2",
							ConflictsWith: []string{LoginVersionVar + ".0." + LoginV1Var},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									BaseURIVar: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optionally specify a base uri of the login UI. If unspecified the default URI will be used.",
									},
								},
							},
						},
					},
				},
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer: helper.ImportWithIDAndOptionalOrg(
			AppIDVar,
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(ClientIDVar, helper.ConvertNonEmpty, true),
			helper.NewImportAttribute(ClientSecretVar, helper.ConvertNonEmpty, true),
		),
	}
}
