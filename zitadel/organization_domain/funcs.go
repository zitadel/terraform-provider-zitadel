package organization_domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteOrganizationDomain(ctx, &org.DeleteOrganizationDomainRequest{
		OrganizationId: d.Get(OrganizationIDVar).(string),
		Domain:         d.Get(DomainVar).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get(OrganizationIDVar).(string)
	domain := d.Get(DomainVar).(string)

	_, err = client.AddOrganizationDomain(ctx, &org.AddOrganizationDomainRequest{
		OrganizationId: orgID,
		Domain:         domain,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	validationType := d.Get(ValidationTypeVar).(string)
	validationResp, err := client.GenerateOrganizationDomainValidation(ctx, &org.GenerateOrganizationDomainValidationRequest{
		OrganizationId: orgID,
		Domain:         domain,
		Type:           org.DomainValidationType(org.DomainValidationType_value[validationType]),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)
	if err := d.Set(ValidationTokenVar, validationResp.GetToken()); err != nil {
		return diag.Errorf("failed to set validation_token: %v", err)
	}
	if err := d.Set(ValidationURLVar, validationResp.GetUrl()); err != nil {
		return diag.Errorf("failed to set validation_url: %v", err)
	}

	if d.Get(VerifyVar).(bool) {
		_, err = client.VerifyOrganizationDomain(ctx, &org.VerifyOrganizationDomainRequest{
			OrganizationId: orgID,
			Domain:         domain,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return read(ctx, d, m)
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	if d.HasChange(VerifyVar) && d.Get(VerifyVar).(bool) {
		client, err := helper.GetOrgClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.VerifyOrganizationDomain(ctx, &org.VerifyOrganizationDomainRequest{
			OrganizationId: d.Get(OrganizationIDVar).(string),
			Domain:         d.Get(DomainVar).(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return read(ctx, d, m)
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get(OrganizationIDVar).(string)
	domain := d.Id()

	resp, err := client.ListOrganizationDomains(ctx, &org.ListOrganizationDomainsRequest{
		OrganizationId: orgID,
		Filters: []*org.DomainSearchFilter{
			{
				Filter: &org.DomainSearchFilter_DomainFilter{
					DomainFilter: &org.OrganizationDomainQuery{
						Domain: domain,
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
		return diag.Errorf("failed to get domain: %v", err)
	}

	if len(resp.Domains) == 0 {
		d.SetId("")
		return nil
	}

	remoteDomain := resp.Domains[0]

	if err := d.Set(DomainVar, remoteDomain.Domain); err != nil {
		return diag.Errorf("failed to set domain: %v", err)
	}
	if err := d.Set(IsVerifiedVar, remoteDomain.IsVerified); err != nil {
		return diag.Errorf("failed to set is_verified: %v", err)
	}
	if err := d.Set(IsPrimaryVar, remoteDomain.IsPrimary); err != nil {
		return diag.Errorf("failed to set is_primary: %v", err)
	}

	return nil
}

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := helper.GetID(d, OrganizationIDVar)
	domain := helper.GetID(d, DomainVar)

	resp, err := client.ListOrganizationDomains(ctx, &org.ListOrganizationDomainsRequest{
		OrganizationId: orgID,
		Filters: []*org.DomainSearchFilter{
			{
				Filter: &org.DomainSearchFilter_DomainFilter{
					DomainFilter: &org.OrganizationDomainQuery{
						Domain: domain,
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
		return diag.Errorf("failed to get domain: %v", err)
	}

	if len(resp.Domains) == 0 {
		d.SetId("")
		return nil
	}

	remoteDomain := resp.Domains[0]

	d.SetId(remoteDomain.Domain)
	if err := d.Set(DomainVar, remoteDomain.Domain); err != nil {
		return diag.Errorf("failed to set domain: %v", err)
	}
	if err := d.Set(IsVerifiedVar, remoteDomain.IsVerified); err != nil {
		return diag.Errorf("failed to set is_verified: %v", err)
	}
	if err := d.Set(IsPrimaryVar, remoteDomain.IsPrimary); err != nil {
		return diag.Errorf("failed to set is_primary: %v", err)
	}
	validationType := org.DomainValidationType_name[int32(remoteDomain.ValidationType)]
	if err := d.Set(ValidationTypeVar, validationType); err != nil {
		return diag.Errorf("failed to set validation_type: %v", err)
	}

	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetOrgClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	orgID := d.Get(OrganizationIDVar).(string)
	domainFilter := d.Get(DomainVar).(string)

	req := &org.ListOrganizationDomainsRequest{
		OrganizationId: orgID,
	}

	if domainFilter != "" {
		req.Filters = []*org.DomainSearchFilter{
			{
				Filter: &org.DomainSearchFilter_DomainFilter{
					DomainFilter: &org.OrganizationDomainQuery{
						Domain: domainFilter,
					},
				},
			},
		}
	}

	resp, err := client.ListOrganizationDomains(ctx, req)
	if err != nil {
		return diag.Errorf("failed to list domains: %v", err)
	}

	domains := make([]interface{}, len(resp.Domains))
	for i, domain := range resp.Domains {
		domainMap := map[string]interface{}{
			DomainVar:         domain.Domain,
			IsVerifiedVar:     domain.IsVerified,
			IsPrimaryVar:      domain.IsPrimary,
			ValidationTypeVar: org.DomainValidationType_name[int32(domain.ValidationType)],
			OrganizationIDVar: domain.OrganizationId,
		}
		domains[i] = domainMap
	}

	d.SetId(fmt.Sprintf("%s", orgID))
	return diag.FromErr(d.Set(domainsVar, domains))
}
