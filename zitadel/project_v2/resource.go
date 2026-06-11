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
				ValidateDiagFunc: helper.NonEmptyString(helper.OrgIDVar),
				Description:      "ID of the organization the project belongs to. Required because the v2 CreateProject API takes the organization as an explicit request field.",
			},
			NameVar: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: helper.NonEmptyString(NameVar),
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
				Description: "Whether the roles assigned to a user are asserted (added) in the access and ID tokens issued for this project.",
			},
			roleCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether ZITADEL checks that the authenticating user has at least one role granted on this project before issuing a token.",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether ZITADEL checks that the user's organization is granted access to this project before issuing a token.",
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
		// Use ConvertNonEmpty rather than the default ConvertID (which applies
		// a strict generated-id format check). The id is validated server-side
		// by the subsequent GetProject call, so accepting any non-empty id at
		// import time avoids rejecting a valid id on a format assumption. This
		// matches the application_v2 importer. Format stays
		// `<project_id[:org_id]>`.
		Importer: helper.ImportWithAttributes(
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertNonEmpty, false),
			helper.ImportOptionalOrgAttribute,
		),
	}
}
