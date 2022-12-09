package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveOrgDomain(ctx, &management.RemoveOrgDomainRequest{
		Domain: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete domain: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(nameVar).(string)
	_, err = client.AddOrgDomain(ctx, &management.AddOrgDomainRequest{
		Domain: name,
	})
	if err != nil {
		return diag.Errorf("failed to create domain: %v", err)
	}
	d.SetId(name)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgDomains(ctx, &management.ListOrgDomainsRequest{
		Queries: []*org.DomainSearchQuery{
			{Query: &org.DomainSearchQuery_DomainNameQuery{
				DomainNameQuery: &org.DomainNameQuery{
					Name:   d.Id(),
					Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
				},
			},
			},
		},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to list domains")
	}

	if len(resp.Result) == 1 {
		domain := resp.Result[0]
		set := map[string]interface{}{
			nameVar:           domain.GetDomainName(),
			orgIDVar:          domain.GetOrgId(),
			isVerifiedVar:     domain.GetIsVerified(),
			isPrimaryVar:      domain.GetIsPrimary(),
			validationTypeVar: domain.GetValidationType().Number(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of domain: %v", k, err)
			}
		}
		d.SetId(domain.GetDomainName())
		return nil
	}

	d.SetId("")
	return nil
}
