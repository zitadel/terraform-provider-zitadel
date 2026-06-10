package project_v2

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	projectpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the project, which can then be granted to different organizations or users directly, containing different applications.",
		Schema: map[string]*schema.Schema{
			// org_id is Required here (unlike most resources where it is
			// optional) because the v2 CreateProject RPC takes the
			// organization as an explicit OrganizationId field in the
			// request body rather than as context metadata. An empty value
			// produces an invalid request that fails at apply time.
			helper.OrgIDVar: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: nonEmptyString(helper.OrgIDVar),
				Description:      "ID of the organization the project belongs to. Required because the v2 CreateProject API takes the organization as an explicit request field.",
			},
			NameVar: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: nonEmptyString(NameVar),
				Description:      "Name of the project",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the project",
			},
			roleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "describes if roles of user should be added in token",
			},
			roleCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ZITADEL checks if the user has at least one on this project",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			privateLabelingSettingVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines from where the private labeling should be triggered" + helper.DescriptionEnumValuesList(projectpb.PrivateLabelingSetting_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(privateLabelingSettingVar, value, projectpb.PrivateLabelingSetting_value)
				},
				Default: defaultPrivateLabelingSetting,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		// Use ConvertNonEmpty rather than the default ConvertID (which
		// enforces a strict ^\d{18}$). ZITADEL ids are not guaranteed to be
		// exactly 18 digits (they can be 19), so a strict check could reject
		// a valid project id on import. This matches the application_v2
		// importer. Format stays `<project_id[:org_id]>`.
		Importer: helper.ImportWithAttributes(
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertNonEmpty, false),
			helper.ImportOptionalOrgAttribute,
		),
	}
}

// nonEmptyString returns a ValidateDiagFunc that rejects empty strings, so a
// Required attribute set to "" fails at plan time instead of producing an
// invalid request at apply.
func nonEmptyString(attr string) schema.SchemaValidateDiagFunc {
	return func(value interface{}, _ cty.Path) diag.Diagnostics {
		if s, _ := value.(string); s == "" {
			return diag.Errorf("%s must not be empty", attr)
		}
		return nil
	}
}
