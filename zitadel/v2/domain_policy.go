package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	admin2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	domainPolicyOrgIdVar              = "org_id"
	domainPolicyUserLoginMustBeDomain = "user_login_must_be_domain"
	domainPolicyValidateOrgDomain     = "validate_org_domains"
	domainPolicySmtpSender            = "smtp_sender_address_matches_instance_domain"
)

func GetDomainPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom domain policy of an organization.",
		Schema: map[string]*schema.Schema{
			domainPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			domainPolicyUserLoginMustBeDomain: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "User login must be domain",
			},
			domainPolicyValidateOrgDomain: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Validate organization domains",
			},
			domainPolicySmtpSender: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "",
			},
		},
		ReadContext:   readDomainPolicy,
		CreateContext: createDomainPolicy,
		DeleteContext: deleteDomainPolicy,
		UpdateContext: updateDomainPolicy,
	}
}

func deleteDomainPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := d.Get(domainPolicyOrgIdVar).(string)

	_, err = client.ResetCustomDomainPolicyToDefault(ctx, &admin2.ResetCustomDomainPolicyToDefaultRequest{
		OrgId: org,
	})
	if err != nil {
		return diag.Errorf("failed to reset domain policy: %v", err)
	}
	return nil
}

func updateDomainPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := d.Get(domainPolicyOrgIdVar).(string)

	_, err = client.UpdateCustomDomainPolicy(ctx, &admin2.UpdateCustomDomainPolicyRequest{
		OrgId:                                  org,
		UserLoginMustBeDomain:                  d.Get(domainPolicyUserLoginMustBeDomain).(bool),
		ValidateOrgDomains:                     d.Get(domainPolicyValidateOrgDomain).(bool),
		SmtpSenderAddressMatchesInstanceDomain: d.Get(domainPolicySmtpSender).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update domain policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createDomainPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := d.Get(domainPolicyOrgIdVar).(string)

	_, err = client.AddCustomDomainPolicy(ctx, &admin2.AddCustomDomainPolicyRequest{
		OrgId:                                  org,
		UserLoginMustBeDomain:                  d.Get(domainPolicyUserLoginMustBeDomain).(bool),
		ValidateOrgDomains:                     d.Get(domainPolicyValidateOrgDomain).(bool),
		SmtpSenderAddressMatchesInstanceDomain: d.Get(domainPolicySmtpSender).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create domain policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readDomainPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(domainPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetDomainPolicy(ctx, &management2.GetDomainPolicyRequest{})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get domain policy: %v", err)
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		domainPolicyOrgIdVar:              policy.GetDetails().GetResourceOwner(),
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
