package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	memberOrgID  = "org_id"
	memberUserID = "user_id"
)

func GetMemberDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			memberOrgID: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the organization",
			},
			memberUserID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
			},
		},
	}
}

func readMembersOfOrg(ctx context.Context, domains *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgDomains(ctx, &management2.ListOrgDomainsRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	for i := range resp.Result {
		domain := resp.Result[i]

		values := map[string]interface{}{
			domainOrgIdVar:       domain.GetOrgId(),
			domainNameVar:        domain.GetDomainName(),
			domainIsVerified:     domain.GetIsVerified(),
			domainIsPrimary:      domain.GetIsPrimary(),
			domainValidationType: int(domain.GetValidationType().Number()),
		}
		domains.Add(values)
	}

	return nil
}
