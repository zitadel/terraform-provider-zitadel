package action_execution_request

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an action execution triggered by a request.",
		Schema: actionexecutionbase.WithTargetIDs(map[string]*schema.Schema{
			MethodVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					MethodVar,
					ServiceVar,
					AllVar,
				},
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					v := i.(string)
					if !regexp.MustCompile(`^/[\w./]+$`).MatchString(v) {
						return diag.Errorf("invalid method: %s. Must start with / and contain only letters, numbers, dots, slashes, and underscores", v)
					}
					return nil
				},
				Description: "The gRPC method to trigger on, e.g., `/zitadel.session.v2.SessionService/ListSessions`",
			},
			ServiceVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					MethodVar,
					ServiceVar,
					AllVar,
				},
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					v := i.(string)
					if !regexp.MustCompile(`^[\w.]+$`).MatchString(v) {
						return diag.Errorf("invalid service: %s. Must contain only letters, numbers, dots, and underscores", v)
					}
					return nil
				},
				Description: "The gRPC service to trigger on, e.g., `zitadel.session.v2.SessionService`",
			},
			AllVar: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					MethodVar,
					ServiceVar,
					AllVar,
				},
				Description: "Trigger on all requests.",
			},
		}),
		CreateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		DeleteContext: actionexecutionbase.NewDeleteExecution(buildCondition),
		ReadContext:   readExecution,
		UpdateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				m := regexp.MustCompile(`^(?:method:(/[-\w./]+)|service:([-\w.]+)|all:?)$`).FindStringSubmatch(d.Id())
				if m == nil {
					return nil, fmt.Errorf("invalid import ID: %s. Must be 'method:/pkg.Service/Method', 'service:pkg.Service', or 'all'/'all:'", d.Id())
				}
				if m[1] != "" {
					d.SetId("request" + m[1])
				} else if m[2] != "" {
					d.SetId("request/" + m[2])
				} else {
					d.SetId("request")
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
