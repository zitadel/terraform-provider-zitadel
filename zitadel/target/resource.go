package target

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"strings"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a target, which can be used in executions.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the target.",
			},
			EndpointVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The endpoint of the target.",
			},
			TargetTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the target. (REST_WEBHOOK, REST_CALL, REST_ASYNC)",
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					val := value.(string)
					validValues := []string{"REST_WEBHOOK", "REST_CALL", "REST_ASYNC"}
					for _, valid := range validValues {
						if val == valid {
							return nil
						}
					}
					return diag.Errorf("%s: invalid value %s, allowed values: %s", path, val, strings.Join(validValues, ", "))
				},
			},
			TimeoutVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timeout defines the duration until ZITADEL cancels the execution.",
			},
			InterruptOnErrorVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Define if any error stops the whole execution. Note: this is not used for REST_ASYNC target type.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithIDAndOptionalOrg(TargetIDVar),
	}
}
