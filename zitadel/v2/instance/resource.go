package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an instance in ZITADEL, which is the highest level.",
		Schema: map[string]*schema.Schema{
			instanceNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance",
			},
			firstOrgNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the first organization created on this instance",
			},
			firstOrgIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the first organization created on this instance",
			},
			customDomainVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerUserNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username of the owner of the instance",
			},
			ownerUserIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the owner of the instance",
			},
			ownerEmailVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom domain used for the instance",
			},
			ownerIsEmailVerifiedVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerFirstNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerLastNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerPreferredLanguageVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerPasswordVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			ownerPasswordChangeRequiredVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
			defaultLanguageVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom domain used for the instance",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
