package action_execution_request

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an action execution triggered by a request.",
		Schema: map[string]*schema.Schema{
			MethodVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					ServiceVar,
					AllVar,
				},
				Description: "The gRPC method to trigger on, e.g., `/zitadel.session.v2.SessionService/ListSessions`",
			},
			ServiceVar: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					MethodVar,
					AllVar,
				},
				Description: "The gRPC service to trigger on, e.g., `zitadel.session.v2.SessionService`",
			},
			AllVar: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					MethodVar,
					ServiceVar,
				},
				Description: "Trigger on all requests.",
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
					return nil, fmt.Errorf("invalid import ID: %s. Must be in format 'subtype:value' (e.g., 'method:/zitadel.session.v2.SessionService/ListSessions') or 'all'", id)
				}

				subType := parts[0]
				value := parts[1]
				internalID := ""

				switch subType {
				case "method":
					internalID = "request" + value
					if err := d.Set(MethodVar, value); err != nil {
						return nil, err
					}
				case "service":
					internalID = "request/" + value
					if err := d.Set(ServiceVar, value); err != nil {
						return nil, err
					}
				case "all":
					internalID = "request"
					if err := d.Set(AllVar, true); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid import subtype: %s. Must be 'method', 'service', or 'all'", subType)
				}

				d.SetId(internalID)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
