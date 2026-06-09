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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the organization the project belongs to. Required because the v2 CreateProject API takes the organization as an explicit request field.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the project",
				/* Not necessary as long as only active projects are created
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return EnumValueValidation(projectStateVar, value, project.ProjectState_value)
				},*/
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
		Importer:      helper.ImportWithIDAndOptionalOrg(ProjectIDVar),
	}
}
