package instance_trusted_domains

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

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(InstanceIDVar).(string)

	resp, err := client.ListTrustedDomains(ctx, &instance.ListTrustedDomainsRequest{
		InstanceId: instanceID,
	})
	if err != nil {
		return diag.Errorf("failed to list trusted domains: %v", err)
	}

	domains := make([]string, 0)
	for _, domain := range resp.GetTrustedDomain() {
		domains = append(domains, domain.GetDomain())
	}

	if err := d.Set(DomainsVar, domains); err != nil {
		return diag.Errorf("failed to set domains: %v", err)
	}

	// Use instance ID as the datasource ID
	if instanceID != "" {
		d.SetId(instanceID)
	} else {
		d.SetId("current")
	}

	return nil
}
