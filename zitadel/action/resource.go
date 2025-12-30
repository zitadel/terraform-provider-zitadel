package action

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an action belonging to an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			stateVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the state of the action",
				/* Not necessary as long as only active actions are created
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return EnumValueValidation(actionState, value, action.ActionState_value)
				},*/
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			ScriptVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			timeoutVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "after which time the action will be terminated if not finished",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			allowedToFailVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "when true, the next action will be called even if this action fails",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithIDAndOptionalOrg(ActionIDVar),
	}
}
