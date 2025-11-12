package action_execution_function

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetDatasource() *schema.Resource {
	return actionexecutionbase.NewActionExecutionDatasource(
		"Datasource representing an action execution triggered by a function.",
		"The ID of this resource. Must be one of: `preuserinfo`, `preaccesstoken`, `presamlresponse`",
		map[string]*schema.Schema{
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the function.",
			},
		},
		readExecution,
	)
}
