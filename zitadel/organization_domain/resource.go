package organization_domain

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a domain of an organization in ZITADEL.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the organization",
			},
			DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name to be added to the organization",
			},
			ValidationTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of domain validation" + helper.DescriptionEnumValuesList(org.DomainValidationType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(ValidationTypeVar, value, org.DomainValidationType_value)
				},
			},
			VerifyVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger domain verification. Set to true after adding DNS/HTTP validation.",
			},
			IsVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the domain has been verified",
			},
			IsPrimaryVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the primary domain of the organization",
			},
			ValidationTokenVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Validation token for domain verification",
			},
			ValidationURLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL where validation file should be hosted for HTTP verification",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithOptionalOrg(helper.NewImportAttribute(DomainVar, helper.ConvertNonEmpty, false)),
	}
}
