package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	actionOrgId         = "org_id"
	actionState         = "state"
	actionName          = "name"
	actionScript        = "script"
	actionTimeout       = "timeout"
	actionAllowedToFail = "allowed_to_fail"
)

func GetActionDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			actionOrgId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the organization",
			},
			actionState: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the state of the action",
			},
			actionName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			actionScript: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			actionTimeout: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "after which time the action will be terminated if not finished",
			},
			actionAllowedToFail: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "when true, the next action will be called even if this action fails",
			},
		},
	}
}

func readActionsOfOrg(ctx context.Context, actions *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActions(ctx, &management2.ListActionsRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	for i := range resp.Result {
		action := resp.Result[i]

		values := map[string]interface{}{
			actionOrgId:         action.GetDetails().GetResourceOwner(),
			actionState:         action.GetState(),
			actionName:          action.GetName(),
			actionScript:        action.GetScript(),
			actionTimeout:       action.GetTimeout().String(),
			actionAllowedToFail: action.GetAllowedToFail(),
		}
		actions.Add(values)
	}

	return nil
}
