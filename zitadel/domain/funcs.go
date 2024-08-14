package domain

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	domainName := d.Id()
	if d.Get(isPrimaryVar).(bool) {
		resp, err := client.ListOrgDomains(helper.CtxWithOrgID(ctx, d), &management.ListOrgDomainsRequest{})
		if err != nil {
			return diag.FromErr(err)
		}
		for _, domain := range resp.Result {
			parts := strings.Split(clientinfo.Domain, ":")
			if domain.IsVerified && domain.DomainName != domainName && strings.HasSuffix(domain.GetDomainName(), parts[0]) {
				if _, err := client.SetPrimaryOrgDomain(helper.CtxWithOrgID(ctx, d), &management.SetPrimaryOrgDomainRequest{Domain: domain.DomainName}); err != nil {
					return diag.FromErr(err)
				}
				break
			}
		}
	}

	_, err = client.RemoveOrgDomain(helper.CtxWithOrgID(ctx, d), &management.RemoveOrgDomainRequest{
		Domain: domainName,
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(NameVar).(string)
	_, err = client.AddOrgDomain(helper.CtxWithOrgID(ctx, d), &management.AddOrgDomainRequest{
		Domain: name,
	})
	if err != nil {
		return diag.Errorf("failed to create domain: %v", err)
	}
	d.SetId(name)
	if d.Get(isPrimaryVar).(bool) {
		_, err = client.SetPrimaryOrgDomain(helper.CtxWithOrgID(ctx, d), &management.SetPrimaryOrgDomainRequest{Domain: name})
		if err != nil {
			return diag.Errorf("failed to set domain primary: %v", err)
		}
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(NameVar).(string)
	d.SetId(name)
	if d.HasChange(isPrimaryVar) {
		if d.Get(isPrimaryVar).(bool) {
			_, err = client.SetPrimaryOrgDomain(helper.CtxWithOrgID(ctx, d), &management.SetPrimaryOrgDomainRequest{Domain: name})
			if err != nil {
				return diag.Errorf("failed to set domain primary: %v", err)
			}
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgDomains(helper.CtxWithOrgID(ctx, d), &management.ListOrgDomainsRequest{
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
			NameVar:           domain.GetDomainName(),
			helper.OrgIDVar:   domain.GetOrgId(),
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
