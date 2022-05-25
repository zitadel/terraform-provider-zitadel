package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	lockoutPolicyOrgIdVar            = "org_id"
	lockoutPolicyMaxPasswordAttempts = "max_password_attempts"
	lockoutPolicyIsDefault           = "is_default"
)

func GetLockoutPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			lockoutPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			lockoutPolicyMaxPasswordAttempts: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correct or the password is reset.",
			},
			lockoutPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
		},
	}
}

func readLockoutPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLockoutPolicy(ctx, &management2.GetLockoutPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		lockoutPolicyOrgIdVar:            policy.GetDetails().GetResourceOwner(),
		lockoutPolicyIsDefault:           policy.GetIsDefault(),
		lockoutPolicyMaxPasswordAttempts: policy.GetMaxPasswordAttempts(),
	}
	policies.Add(values)
	return nil
}
