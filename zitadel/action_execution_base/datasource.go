package action_execution_base

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func NewActionExecutionDatasource(
	datasourceDescription string,
	idDescription string,
	specificSchema map[string]*schema.Schema,
	readFunc schema.ReadContextFunc,
) *schema.Resource {
	specificSchema[IDVar] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: idDescription,
	}
	specificSchema[TargetIDsVar] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "The list of target IDs to call.",
	}

	return &schema.Resource{
		Description: datasourceDescription,
		Schema:      specificSchema,
		ReadContext: readFunc,
	}
}
