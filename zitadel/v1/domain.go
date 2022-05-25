package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	domainOrgIdVar       = "org_id"
	domainNameVar        = "name"
	domainIsVerified     = "is_verified"
	domainIsPrimary      = "is_primary"
	domainValidationType = "validation_type"
)

func GetDomainDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			domainNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the domain",
			},
			domainOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			domainIsVerified: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Is domain verified",
			},
			domainIsPrimary: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Is domain primary",
			},
			domainValidationType: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Validation type",
			},
		},
	}
}

func readDomainsOfOrg(ctx context.Context, domains *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
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
