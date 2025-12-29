package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	instance "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")

	clientInfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(InstanceIDVar).(string)

	instanceReq := &instance.GetInstanceRequest{}
	if instanceID != "" {
		instanceReq.InstanceId = instanceID
	}

	inst, err := client.GetInstance(ctx, instanceReq)
	if err != nil {
		return diag.Errorf("failed to get instance: %v", err)
	}

	if inst == nil || inst.Instance == nil {
		return diag.Errorf("instance not found")
	}

	trustedDomainsReq := &instance.ListTrustedDomainsRequest{}
	if instanceID != "" {
		trustedDomainsReq.InstanceId = instanceID
	}

	trustedDomainsResp, err := client.ListTrustedDomains(ctx, trustedDomainsReq)
	if err != nil {
		return diag.Errorf("failed to list trusted domains: %v", err)
	}

	customDomains := make([]map[string]interface{}, 0)
	if inst.Instance.CustomDomains != nil {
		customDomains = make([]map[string]interface{}, len(inst.Instance.CustomDomains))
		for i, domain := range inst.Instance.CustomDomains {
			if domain != nil {
				customDomains[i] = map[string]interface{}{
					"domain":    domain.Domain,
					"primary":   domain.Primary,
					"generated": domain.Generated,
				}
			}
		}
	}

	trustedDomains := make([]string, 0)
	if trustedDomainsResp != nil && trustedDomainsResp.TrustedDomain != nil {
		trustedDomains = make([]string, len(trustedDomainsResp.TrustedDomain))
		for i, domain := range trustedDomainsResp.TrustedDomain {
			if domain != nil {
				trustedDomains[i] = domain.Domain
			}
		}
	}

	set := map[string]interface{}{
		NameVar:           inst.Instance.Name,
		VersionVar:        inst.Instance.Version,
		StateVar:          inst.Instance.State.String(),
		CustomDomainsVar:  customDomains,
		TrustedDomainsVar: trustedDomains,
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
