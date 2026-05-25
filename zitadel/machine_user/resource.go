package machine_user

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Machine user (service account) under an organization. Backward-compatible: tries the user/v2 API first and falls back to the management API, so it works with both ZITADEL 3.x and 4.x.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of this resource. Optionally set a custom unique ID. If omitted, ZITADEL will generate one.",
			},
			userStateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the user",
				/* Not necessary as long as only active users are created
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return EnumValueValidation(userStateVar, value.(string), user.UserState_value)
				},*/
			},
			UserNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username",
			},
			loginNamesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Loginnames",
			},
			preferredLoginNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Preferred login name",
			},
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the machine user",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the user",
			},
			accessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type" + helper.DescriptionEnumValuesList(user.AccessTokenType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(accessTokenTypeVar, value, user.AccessTokenType_value)
				},
				Default: defaultAccessTokenType,
			},
			WithSecretVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Generate machine secret, only applicable if creation or change from false",
			},
			clientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the client ID if withSecret is true",
				Sensitive:   true,
			},
			clientSecretVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the client secret if withSecret is true",
				Sensitive:   true,
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		Importer: helper.ImportWithIDAndOptionalOrg(
			UserIDVar,
			helper.NewImportAttribute(WithSecretVar, helper.ConvertBool, false),
			helper.NewImportAttribute(clientIDVar, helper.ConvertNonEmpty, true),
			helper.NewImportAttribute(clientSecretVar, helper.ConvertNonEmpty, true),
		),
	}
}
