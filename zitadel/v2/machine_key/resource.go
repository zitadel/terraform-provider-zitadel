package machine_key

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a machine key",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			userIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			keyTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the machine key" + helper.DescriptionEnumValuesList(authn.KeyType_name),
				ForceNew:    true,
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(keyTypeVar, value, authn.KeyType_value)
				},
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration date of the machine key in the RFC3339 format",
				ForceNew:    true,
				Computed:    true,
			},
			keyDetailsVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the machine key",
				Sensitive:   true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
