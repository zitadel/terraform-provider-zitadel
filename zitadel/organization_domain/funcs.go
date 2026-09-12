package organization_domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	orgV2 "github.com/zitadel/zitadel-go/v3/pkg/client/org/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	// The domain exists from here on, so claim the ID before anything else can
	// fail. Otherwise an error in one of the calls below leaves a domain that
	// ZITADEL knows about and Terraform does not, which can only be cleaned up
	// by hand in the console.
	d.SetId(domain)

	// Unless the organizations effective domain policy requires org domains to
	// be validated, ZITADEL verifies the domain while adding it, so there is
	// nothing left to challenge and asking for a validation anyway is rejected
	// with ORG-HGw21. The policy is asked rather than the domain itself because
	// the projection a domain read goes through can still report a freshly
	// added and auto-verified domain as unverified.
	requiresValidation, err := domainValidationRequired(ctx, clientinfo, orgID)
	if err != nil {
		return diag.Errorf("failed to get domain policy of organization %s: %v", orgID, err)
	}
	verified := !requiresValidation

	if !verified && hasValidationType(d) {
		validationType := d.Get(ValidationTypeVar).(string)
		validationResp, err := client.GenerateOrganizationDomainValidation(ctx, &org.GenerateOrganizationDomainValidationRequest{
			OrganizationId: orgID,
			Domain:         domain,
			Type:           org.DomainValidationType(org.DomainValidationType_value[validationType]),
		})
		switch {
		case isAlreadyVerifiedError(err):
			verified = true
		case err != nil:
			return diag.FromErr(err)
		default:
			if err := d.Set(ValidationTokenVar, validationResp.GetToken()); err != nil {
				return diag.Errorf("failed to set validation_token: %v", err)
			}
			if err := d.Set(ValidationURLVar, validationResp.GetUrl()); err != nil {
				return diag.Errorf("failed to set validation_url: %v", err)
			}
		}
	}

	if !verified && d.Get(VerifyVar).(bool) {
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
	// A newly added domain is never the primary one; a later read picks up the
	// change if it is promoted afterwards.
	if err := d.Set(IsPrimaryVar, false); err != nil {
		return diag.Errorf("failed to set %s: %v", IsPrimaryVar, err)
	}
	return nil
}

// domainValidationRequired reports whether the organizations effective domain
// policy makes ZITADEL require ownership validation for added org domains. When
// it does not - the default - a domain is verified as part of being added.
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

// hasValidationType reports whether the configuration asks for a validation
// challenge. validation_type is optional, and the unspecified enum member is
// rejected by the API, so both are treated as "no challenge wanted".
func hasValidationType(d *schema.ResourceData) bool {
	validationType := d.Get(ValidationTypeVar).(string)
	return validationType != "" &&
		validationType != org.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED.String()
}

// alreadyVerifiedErrorID is the ZITADEL error ID for generating a validation
// challenge for a domain that is already verified.
const alreadyVerifiedErrorID = "ORG-HGw21"

// isAlreadyVerifiedError reports whether err is that rejection. create decides
// whether to generate a challenge from a projection that can lag behind the
// command that added the domain, so the error is handled as well as avoided.
// The ID is matched rather than the bare FailedPrecondition code, which the API
// also uses to reject requests that genuinely have to fail the apply.
func isAlreadyVerifiedError(err error) bool {
	if status.Code(err) != codes.FailedPrecondition {
		return false
	}
	return strings.Contains(status.Convert(err).Message(), alreadyVerifiedErrorID)
}

// fetchDomain returns the named domain of an organization, or nil when it is
// not (yet) listed.
func fetchDomain(ctx context.Context, client *orgV2.Client, orgID, domain string) (*org.Domain, error) {
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
	if err != nil {
		return nil, err
	}
	if len(resp.Domains) == 0 {
		return nil, nil
	}
	return resp.Domains[0], nil
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

	remoteDomain, err := fetchDomain(ctx, client, orgID, domain)
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get domain: %v", err)
	}

	if remoteDomain == nil {
		d.SetId("")
		return nil
	}

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

	remoteDomain, err := fetchDomain(ctx, client, orgID, domain)
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get domain: %v", err)
	}

	if remoteDomain == nil {
		d.SetId("")
		return nil
	}

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
