package trigger_actions

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing triggers, when actions get started",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			FlowTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the flow to which the action triggers belong" + helper.DescriptionEnumValuesList(FlowTypes()),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(FlowTypeVar, value, helper.EnumValueMap(FlowTypes()))
				},
				ForceNew: true,
			},
			TriggerTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trigger type on when the actions get triggered" + helper.DescriptionEnumValuesList(TriggerTypes()),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(TriggerTypeVar, value, helper.EnumValueMap(TriggerTypes()))
				},
				ForceNew: true,
			},
			actionsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "IDs of the triggered actions",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer: helper.ImportWithEmptyID(
			helper.NewImportAttribute(FlowTypeVar, helper.ConvertNonEmpty, false),
			helper.NewImportAttribute(TriggerTypeVar, helper.ConvertNonEmpty, false),
			helper.ImportOptionalOrgAttribute,
		),
	}
}

func FlowTypes() map[int32]string {
	return map[int32]string{
		1: "FLOW_TYPE_EXTERNAL_AUTHENTICATION",
		2: "FLOW_TYPE_CUSTOMISE_TOKEN",
		3: "FLOW_TYPE_INTERNAL_AUTHENTICATION",
		4: "FLOW_TYPE_SAML_RESPONSE",
	}
}
func TriggerTypes() map[int32]string {
	return map[int32]string{
		1: "TRIGGER_TYPE_POST_AUTHENTICATION",
		2: "TRIGGER_TYPE_PRE_CREATION",
		3: "TRIGGER_TYPE_POST_CREATION",
		4: "TRIGGER_TYPE_PRE_USERINFO_CREATION",
		5: "TRIGGER_TYPE_PRE_ACCESS_TOKEN_CREATION",
		6: "TRIGGER_TYPE_PRE_SAML_RESPONSE_CREATION",
	}
}
