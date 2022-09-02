package label_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom label policy of an organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			primaryColorVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for primary color",
			},
			hideLoginNameSuffixVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this scope in https://docs.zitadel.ch/concepts#Reserved_Scopes",
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
			logoURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			iconURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			logoURLDarkVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			iconURLDarkVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			fontURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			setActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set the label policy active after creating/updating",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
