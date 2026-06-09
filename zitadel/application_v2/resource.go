package application_v2

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// GetResource returns the unified zitadel_application_v2 resource.
//
// The Zitadel v2 application API collapses what were three separate v1
// management endpoints (AddOIDCApp / AddAPIApp / AddSAMLApp) into one
// CreateApplication call where the per-type configuration is carried in a
// protobuf oneof. The Terraform surface mirrors that: a single resource
// with three mutually exclusive nested blocks — oidc{}, saml{}, api{} —
// enforced via ExactlyOneOf.
func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an application belonging to a project. Exposes the unified Zitadel Application v2 API; the application type is selected by populating exactly one of the oidc/saml/api configuration blocks.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project this application belongs to.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application.",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the application.",
			},

			oidcBlockVar: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					oidcBlockVar, samlBlockVar, apiBlockVar,
				},
				Description: "OIDC configuration. Mutually exclusive with `saml` and `api`.",
				Elem:        oidcConfigSchema(),
			},
			samlBlockVar: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					oidcBlockVar, samlBlockVar, apiBlockVar,
				},
				Description: "SAML configuration. Mutually exclusive with `oidc` and `api`.",
				Elem:        samlConfigSchema(),
			},
			apiBlockVar: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ExactlyOneOf: []string{
					oidcBlockVar, samlBlockVar, apiBlockVar,
				},
				Description: "API configuration. Mutually exclusive with `oidc` and `saml`.",
				Elem:        apiConfigSchema(),
			},
		},
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Importer: helper.ImportWithIDAndOptionalOrg(
			AppIDVar,
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertID, false),
		),
	}
}

func oidcConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			redirectURIsVar: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Allowed redirect URIs.",
			},
			responseTypesVar: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Response types" + helper.DescriptionEnumValuesList(apppb.OIDCResponseType_name),
			},
			grantTypesVar: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Grant types" + helper.DescriptionEnumValuesList(apppb.OIDCGrantType_name),
			},
			appTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application type" + helper.DescriptionEnumValuesList(apppb.OIDCApplicationType_name),
				ValidateDiagFunc: func(value interface{}, _ cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(appTypeVar, value, apppb.OIDCApplicationType_value)
				},
				Default: apppb.OIDCApplicationType_name[0],
			},
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auth method type" + helper.DescriptionEnumValuesList(apppb.OIDCAuthMethodType_name),
				ValidateDiagFunc: func(value interface{}, _ cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(authMethodTypeVar, value, apppb.OIDCAuthMethodType_value)
				},
				Default: apppb.OIDCAuthMethodType_name[0],
			},
			postLogoutRedirectURIsVar: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Post-logout redirect URIs.",
			},
			versionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OIDC version" + helper.DescriptionEnumValuesList(apppb.OIDCVersion_name),
				Default:     apppb.OIDCVersion_name[0],
			},
			devModeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Development mode (relaxes redirect-URI validation).",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type" + helper.DescriptionEnumValuesList(apppb.OIDCTokenType_name),
				ValidateDiagFunc: func(value interface{}, _ cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(accessTokenTypeVar, value, apppb.OIDCTokenType_value)
				},
				Default: apppb.OIDCTokenType_name[0],
			},
			accessTokenRoleAssertionVar: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			idTokenRoleAssertionVar: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			idTokenUserinfoAssertionVar: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			clockSkewVar: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "0s",
				Description:      "Allowed clock skew (Go duration string, e.g. `5s`).",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			additionalOriginsVar: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Additional allowed origins.",
			},
			skipNativeAppSuccessPageVar: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			backChannelLogoutURIVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Back-channel logout URI used by ZITADEL to notify the application of terminated sessions (OIDC Back-Channel Logout). Computed if not set, so the server-side default flows back into state.",
			},
			loginVersionVar: loginVersionSchema(oidcBlockVar),

			clientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated client ID.",
			},
			clientSecretVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated client secret (only set on create when the auth method requires one).",
			},
			noneCompliantVar: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			complianceProblemsVar: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						complianceKeyVar:     {Type: schema.TypeString, Computed: true},
						complianceMessageVar: {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func samlConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			metadataXMLVar: {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{samlBlockVar + ".0." + metadataXMLVar, samlBlockVar + ".0." + metadataURLVar},
				Description:  "SAML metadata as raw XML. Mutually exclusive with `metadata_url`. Marked sensitive because SAML metadata documents commonly embed signing/encryption certificates.",
			},
			metadataURLVar: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{samlBlockVar + ".0." + metadataXMLVar, samlBlockVar + ".0." + metadataURLVar},
				Description:  "URL from which SAML metadata can be fetched. Mutually exclusive with `metadata_xml`.",
			},
			loginVersionVar: loginVersionSchema(samlBlockVar),
		},
	}
}

func apiConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			authMethodTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API auth method type" + helper.DescriptionEnumValuesList(apppb.APIAuthMethodType_name),
				ValidateDiagFunc: func(value interface{}, _ cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(authMethodTypeVar, value, apppb.APIAuthMethodType_value)
				},
				Default: apppb.APIAuthMethodType_name[0],
			},
			clientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated client ID.",
			},
			clientSecretVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated client secret (only returned on create).",
			},
		},
	}
}

// loginVersionSchema returns the shared login_version sub-block. The
// parentBlock argument is the top-level resource attribute that contains
// this login_version (e.g. "oidc" or "saml"); it's used to build absolute
// ConflictsWith paths so that login_v1 and login_v2 cannot both be set.
func loginVersionSchema(parentBlock string) *schema.Schema {
	v1Path := parentBlock + ".0." + loginVersionVar + ".0." + loginV1Var
	v2Path := parentBlock + ".0." + loginVersionVar + ".0." + loginV2Var
	return &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Login UI version to use for this application. Exactly one of `login_v1` and `login_v2` may be set. Computed so that the server-side default flows back into state when the user omits this block.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				loginV1Var: {
					Type:          schema.TypeBool,
					Optional:      true,
					Description:   "Use the legacy Login UI (V1).",
					ConflictsWith: []string{v2Path},
				},
				loginV2Var: {
					Type:          schema.TypeList,
					MaxItems:      1,
					Optional:      true,
					Description:   "Use the Login UI V2.",
					ConflictsWith: []string{v1Path},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							baseURIVar: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Optional base URI of a custom Login UI V2. If unset, the instance default is used.",
							},
						},
					},
				},
			},
		},
	}
}
