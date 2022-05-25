package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	privacyPolicyOrgIdVar    = "org_id"
	privacyPolicyTOSLink     = "tos_link"
	privacyPolicyPrivacyLink = "privacy_link"
	privacyPolicyIsDefault   = "is_default"
	privacyPolicyHelpLink    = "help_link"
)

func GetPrivacyPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			privacyPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			privacyPolicyTOSLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			privacyPolicyPrivacyLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			privacyPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			privacyPolicyHelpLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func readPrivacyPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPrivacyPolicy(ctx, &management2.GetPrivacyPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		privacyPolicyOrgIdVar:    policy.GetDetails().GetResourceOwner(),
		privacyPolicyTOSLink:     policy.GetTosLink(),
		privacyPolicyPrivacyLink: policy.GetPrivacyLink(),
		privacyPolicyIsDefault:   policy.GetIsDefault(),
		privacyPolicyHelpLink:    policy.GetHelpLink(),
	}
	policies.Add(values)
	return nil
}
