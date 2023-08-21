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
			helper.OrgIDVar: helper.OrgIDResourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			AppIDVar: {
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
			ExpirationDateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration date of the app key in the RFC3339 format",
				ForceNew:    true,
			},
			KeyDetailsVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the app key",
				Sensitive:   true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		ReadContext:   read,
		Importer: helper.ImportWithIDAndOptionalOrg(
			keyIDVar,
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(AppIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(KeyDetailsVar, helper.ConvertJSON, true),
		),
	}
}
