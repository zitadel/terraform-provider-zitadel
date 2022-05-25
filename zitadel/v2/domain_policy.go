package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	domainPolicyOrgIdVar              = "org_id"
	domainPolicyUserLoginMustBeDomain = "user_login_must_be_domain"
	domainPolicyIsDefault             = "is_default"
	domainPolicyValidateOrgDomain     = "validate_org_domains"
)

func GetDomainPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			domainPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			domainPolicyUserLoginMustBeDomain: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User login must be domain",
			},
			domainPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is this policy the default",
			},
			domainPolicyValidateOrgDomain: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Validate organization domains",
			},
		},
		ReadContext: readDomainPolicy,
	}
}

func readDomainPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetDomainPolicy(ctx, &management2.GetDomainPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get domain policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		domainPolicyOrgIdVar:              policy.GetDetails().GetResourceOwner(),
		domainPolicyIsDefault:             policy.GetIsDefault(),
		domainPolicyUserLoginMustBeDomain: policy.GetUserLoginMustBeDomain(),
		domainPolicyValidateOrgDomain:     policy.GetValidateOrgDomains(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of domain: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
