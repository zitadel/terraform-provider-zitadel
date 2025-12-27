package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	instance "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientInfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(InstanceIDVar).(string)

	req := &instance.GetInstanceRequest{}
	if instanceID != "" {
		req.InstanceId = instanceID
	}

	inst, err := client.GetInstance(ctx, req)
	if err != nil {
		return diag.Errorf("failed to get instance: %v", err)
	}

	customDomainsReq := &instance.ListCustomDomainsRequest{}
	if instanceID != "" {
		customDomainsReq.InstanceId = instanceID
	}

	customDomainsResp, err := client.ListCustomDomains(ctx, customDomainsReq)
	if err != nil {
		return diag.Errorf("failed to list custom domains: %v", err)
	}

	trustedDomainsReq := &instance.ListTrustedDomainsRequest{}
	if instanceID != "" {
		trustedDomainsReq.InstanceId = instanceID
	}

	trustedDomainsResp, err := client.ListTrustedDomains(ctx, trustedDomainsReq)
	if err != nil {
		return diag.Errorf("failed to list trusted domains: %v", err)
	}

	customDomains := make([]string, len(customDomainsResp.Domains))
	for i, domain := range customDomainsResp.Domains {
		customDomains[i] = domain.Domain
	}

	trustedDomains := make([]string, len(trustedDomainsResp.TrustedDomain))
	for i, domain := range trustedDomainsResp.TrustedDomain {
		trustedDomains[i] = domain.Domain
	}

	var generatedDomain string
	var primaryDomain string

	for _, domain := range customDomainsResp.Domains {
		if domain.Generated {
			generatedDomain = domain.Domain
		}
		if domain.Primary {
			primaryDomain = domain.Domain
		}
	}

	if primaryDomain == "" && len(customDomains) > 0 {
		primaryDomain = customDomains[0]
	}

	set := map[string]interface{}{
		NameVar:            inst.Instance.Name,
		PrimaryDomainVar:   primaryDomain,
		GeneratedDomainVar: generatedDomain,
		CustomDomainsVar:   customDomains,
		TrustedDomainsVar:  trustedDomains,
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s: %v", k, err)
		}
	}

	if instanceID != "" {
		d.SetId(instanceID)
	} else {
		d.SetId(inst.Instance.Id)
	}

	return nil
}
