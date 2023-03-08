package app_key

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/authn"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a app key",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			appIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the application",
				ForceNew:    true,
			},
			keyTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the app key" + helper.DescriptionEnumValuesList(authn.KeyType_name),
				ForceNew:    true,
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(keyTypeVar, value, authn.KeyType_value)
				},
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration date of the app key in the RFC3339 format",
				ForceNew:    true,
			},
			keyDetailsVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the app key",
				Sensitive:   true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
