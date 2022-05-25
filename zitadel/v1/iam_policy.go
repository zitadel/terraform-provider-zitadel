package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	iamPolicyOrgIdVar  = "org_id"
	iamPolicyUserLogin = "user_login"
	iamPolicyIsDefault = "is_default"
)

func GetIAMPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			iamPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			iamPolicyUserLogin: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User login must be domain",
			},
			iamPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is this policy the default",
			},
		},
	}
}

func readIAMPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetOrgIAMPolicy(ctx, &management2.GetOrgIAMPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get iam policy: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		iamPolicyOrgIdVar:  policy.GetDetails().GetResourceOwner(),
		iamPolicyUserLogin: policy.GetUserLoginMustBeDomain(),
		iamPolicyIsDefault: policy.GetIsDefault(),
	}
	policies.Add(values)
	return nil
}
