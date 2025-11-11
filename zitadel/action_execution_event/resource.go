package action_execution_event

import (
	"context"
	"fmt"
	"strings"

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
		CreateContext: actionexecutionbase.NewSetExecution(buildCondition),
		DeleteContext: actionexecutionbase.NewDeleteExecution(buildCondition),
		ReadContext:   readExecution,
		UpdateContext: actionexecutionbase.NewSetExecution(buildCondition),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				parts := strings.SplitN(id, ":", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid import ID: %s. Must be in format 'subtype:value' (e.g., 'event:user.human.added') or 'all'", id)
				}

				subType := parts[0]
				value := parts[1]
				internalID := ""

				switch subType {
				case "event":
					internalID = "event/" + value
					if err := d.Set(EventVar, value); err != nil {
						return nil, err
					}
				case "group":
					internalID = "event/" + value
					if !strings.HasSuffix(internalID, ".*") {
						internalID += ".*"
					}
					if err := d.Set(GroupVar, value); err != nil {
						return nil, err
					}
				case "all":
					internalID = "event"
					if err := d.Set(AllVar, true); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid import subtype: %s. Must be 'event', 'group', or 'all'", subType)
				}

				d.SetId(internalID)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
