package action_execution

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewActionExecutionDatasource(
	datasourceDescription string,
	idDescription string,
	specificSchema map[string]*schema.Schema,
	readFunc schema.ReadContextFunc,
) *schema.Resource {

	baseSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: idDescription,
		},
		TargetIDsVar: {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed:    true,
			Description: "The list of target IDs to call.",
		},
	}

	for key, val := range specificSchema {
		baseSchema[key] = val
	}

	return &schema.Resource{
		Description: datasourceDescription,
		Schema:      baseSchema,
		ReadContext: readFunc,
	}
}
