package organization_domain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
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

	// Without validate_org_domains ZITADEL verifies the domain while adding it,
	// leaving no challenge to generate.
	requiresValidation, err := domainValidationRequired(ctx, clientinfo, orgID)
	if err != nil {
		return diag.Errorf("failed to get domain policy: %v", err)
	}

	verify := d.Get(VerifyVar).(bool)
	if requiresValidation && verify && !hasValidationType(d) {
		return verifyNeedsValidationTypeError(orgID)
	}

	_, err = client.AddOrganizationDomain(ctx, &org.AddOrganizationDomainRequest{
		OrganizationId: orgID,
		Domain:         domain,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain)

	verified := !requiresValidation
	if !verified && hasValidationType(d) {
		validationType := d.Get(ValidationTypeVar).(string)
		validationResp, err := client.GenerateOrganizationDomainValidation(ctx, &org.GenerateOrganizationDomainValidationRequest{
			OrganizationId: orgID,
			Domain:         domain,
			Type:           org.DomainValidationType(org.DomainValidationType_value[validationType]),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(ValidationTokenVar, validationResp.GetToken()); err != nil {
			return diag.Errorf("failed to set validation_token: %v", err)
		}
		if err := d.Set(ValidationURLVar, validationResp.GetUrl()); err != nil {
			return diag.Errorf("failed to set validation_url: %v", err)
		}
	}

	if !verified && verify {
		_, err = client.VerifyOrganizationDomain(ctx, &org.VerifyOrganizationDomainRequest{
			OrganizationId: orgID,
			Domain:         domain,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		verified = true
	}

	if err := d.Set(IsVerifiedVar, verified); err != nil {
		return diag.Errorf("failed to set %s: %v", IsVerifiedVar, err)
	}
	if err := d.Set(IsPrimaryVar, false); err != nil {
		return diag.Errorf("failed to set %s: %v", IsPrimaryVar, err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	if d.HasChange(VerifyVar) && d.Get(VerifyVar).(bool) {
		orgID := d.Get(OrganizationIDVar).(string)
		if !hasValidationType(d) {
			requiresValidation, err := domainValidationRequired(ctx, clientinfo, orgID)
			if err != nil {
				return diag.Errorf("failed to get domain policy: %v", err)
			}
			if requiresValidation {
				return verifyNeedsValidationTypeError(orgID)
			}
		}

		client, err := helper.GetOrgClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.VerifyOrganizationDomain(ctx, &org.VerifyOrganizationDomainRequest{
			OrganizationId: orgID,
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

// hasValidationType reports whether the configuration asks for a validation
// challenge. The unspecified enum member is rejected by the API.
func hasValidationType(d *schema.ResourceData) bool {
	validationType := d.Get(ValidationTypeVar).(string)
	return validationType != "" &&
		validationType != org.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED.String()
}

// verifyNeedsValidationTypeError reports that verify cannot be honoured because
// there is no validation challenge for ZITADEL to check.
func verifyNeedsValidationTypeError(orgID string) diag.Diagnostics {
	return diag.Errorf(
		"%s needs %s when the domain policy of organization %s has validate_org_domains enabled, "+
			"because ZITADEL verifies a domain against a published validation challenge",
		VerifyVar, ValidationTypeVar, orgID)
}

// domainValidationRequired reports whether the organization's domain policy
// requires org domains to be validated. It reads the policy through the v1
// management API because settings/v2 GetDomainSettings requires the policy.read
// permission while GetDomainPolicy only requires authenticated, and adding a
// domain itself only needs org.write.
func domainValidationRequired(ctx context.Context, clientinfo *helper.ClientInfo, orgID string) (bool, error) {
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return false, err
	}
	resp, err := client.GetDomainPolicy(helper.CtxSetOrgID(ctx, orgID), &management.GetDomainPolicyRequest{})
	if err != nil {
		return false, err
	}
	return resp.GetPolicy().GetValidateOrgDomains(), nil
}
