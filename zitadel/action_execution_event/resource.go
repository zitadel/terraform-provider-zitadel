package action_execution_event

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an action execution triggered by an event.",
		Schema: map[string]*schema.Schema{
			EventVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					GroupVar,
					AllVar,
				},
				Description: "The specific event to trigger on, e.g., `user.human.added`",
			},
			GroupVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					EventVar,
					AllVar,
				},
				Description: "The event group to trigger on, e.g., `user.human`",
			},
			AllVar: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					EventVar,
					GroupVar,
				},
				Description: "Trigger on all events.",
			},
			actionexecutionbase.TargetIDsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The list of target IDs to call.",
			},
		},
		CreateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		DeleteContext: actionexecutionbase.NewDeleteExecution(buildCondition),
		ReadContext:   readExecution,
		UpdateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				m := regexp.MustCompile(`^(?:event:([-\w.]+)|group:([-\w.]+)(?:\.\*)?|all:?)$`).FindStringSubmatch(d.Id())
				if m == nil {
					return nil, fmt.Errorf("invalid import ID: %s. Must be 'event:name', 'group:name' (optionally ending with '.*'), or 'all'/'all:'", d.Id())
				}
				if m[1] != "" {
					d.SetId("event/" + m[1])
				} else if m[2] != "" {
					d.SetId("event/" + m[2] + ".*")
				} else {
					d.SetId("event")
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
