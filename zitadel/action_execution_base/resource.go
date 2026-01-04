package action_execution_base

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func WithTargetIDs(specific map[string]*schema.Schema) map[string]*schema.Schema {
	specific[TargetIDsVar] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "The list of target IDs to call.",
	}
	return specific
}
