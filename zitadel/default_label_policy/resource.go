package default_label_policy

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/policy"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default label policy.",
		Schema: map[string]*schema.Schema{
			PrimaryColorVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for primary color",
			},
			hideLoginNameSuffixVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this scope in https://zitadel.ch/docs/concepts#Reserved_Scopes",
			},
			warnColorVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for warn color",
			},
			backgroundColorVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for background color",
			},
			fontColorVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for font color",
			},
			primaryColorDarkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for primary color dark theme",
			},
			backgroundColorDarkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for background color dark theme",
			},
			warnColorDarkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for warn color dark theme",
			},
			fontColorDarkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for font color dark theme",
			},
			disableWatermarkVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "disable watermark",
			},
			LogoPathVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			LogoHashVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			logoURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			IconPathVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			IconHashVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			iconURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			LogoDarkPathVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			LogoDarkHashVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			logoURLDarkVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			IconDarkPathVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			IconDarkHashVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			iconURLDarkVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			FontPathVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			FontHashVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			fontURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			SetActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set the label policy active after creating/updating",
			},
			themeModeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "theme mode" + helper.DescriptionEnumValuesList(policy.ThemeMode_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(themeModeVar, value, policy.ThemeMode_value)
				},
				Default: policy.ThemeMode_THEME_MODE_AUTO.String(),
			},
		},
		ReadContext:   read,
		CreateContext: update,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      helper.ImportWithEmptyID(),
	}
}
