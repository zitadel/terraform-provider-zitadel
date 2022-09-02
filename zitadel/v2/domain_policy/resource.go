package domain_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom domain policy of an organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			userLoginMustBeDomainVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "User login must be domain",
			},
			validateOrgDomainVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Validate organization domains",
			},
			smtpSenderVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
