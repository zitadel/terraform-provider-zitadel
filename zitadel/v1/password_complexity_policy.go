package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	passwordCompPolicyOrgIdVar     = "org_id"
	passwordCompPolicyMinLength    = "min_length"
	passwordCompPolicyHasUppercase = "has_uppercase"
	passwordCompPolicyHasLowercase = "has_lowercase"
	passwordCompPolicyHasNumber    = "has_number"
	passwordCompPolicyHasSymbol    = "has_symbol"
	passwordCompPolicyIsDefault    = "is_default"
)

func GetPasswordComplexityPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			passwordCompPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			passwordCompPolicyMinLength: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimal length for the password",
			},
			passwordCompPolicyHasUppercase: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain an upper case letter",
			},
			passwordCompPolicyHasLowercase: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a lower case letter",
			},
			passwordCompPolicyHasNumber: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a number",
			},
			passwordCompPolicyHasSymbol: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a symbol. E.g. \"$\"",
			},
			passwordCompPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
		},
	}
}

func readPasswordComplexityPolicyPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPasswordComplexityPolicy(ctx, &management2.GetPasswordComplexityPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		passwordCompPolicyOrgIdVar:     policy.GetDetails().GetResourceOwner(),
		passwordCompPolicyMinLength:    policy.GetMinLength(),
		passwordCompPolicyHasUppercase: policy.GetHasUppercase(),
		passwordCompPolicyHasLowercase: policy.GetHasLowercase(),
		passwordCompPolicyHasNumber:    policy.GetHasNumber(),
		passwordCompPolicyHasSymbol:    policy.GetHasSymbol(),
		passwordCompPolicyIsDefault:    policy.GetIsDefault(),
	}
	policies.Add(values)
	return nil
}
