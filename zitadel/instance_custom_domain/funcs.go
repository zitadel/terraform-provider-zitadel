package instance_custom_domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	instance "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(InstanceIDVar).(string)
	domain := d.Get(DomainVar).(string)

	_, err = client.AddCustomDomain(ctx, &instance.AddCustomDomainRequest{
		InstanceId:   instanceID,
		CustomDomain: domain,
	})
	if err != nil {
		return diag.Errorf("failed to add custom domain: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, domain))

	// The add call can return before list APIs reflect the new domain.
	if err := waitUntilCustomDomainExists(ctx, client, instanceID, domain); err != nil {
		return diag.Errorf("failed waiting for custom domain to be readable: %v", err)
	}

	return read(ctx, d, m)
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, expected instance_id/domain")
	}
	instanceID := parts[0]
	domain := parts[1]

	resp, err := client.ListCustomDomains(ctx, &instance.ListCustomDomainsRequest{
		InstanceId: instanceID,
	})
	if err != nil {
		return diag.Errorf("failed to list custom domains: %v", err)
	}

	found := false
	for _, domainEntry := range resp.GetDomains() {
		if equalsDomain(domainEntry.GetDomain(), domain) {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		InstanceIDVar: instanceID,
		DomainVar:     domain,
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s: %v", k, err)
		}
	}

	return nil
}

func waitUntilCustomDomainExists(ctx context.Context, client instance.InstanceServiceClient, instanceID, domain string) error {
	deadline := time.Now().Add(30 * time.Second)

	for {
		resp, err := client.ListCustomDomains(ctx, &instance.ListCustomDomainsRequest{
			InstanceId: instanceID,
		})
		if err != nil {
			return fmt.Errorf("failed to list custom domains: %w", err)
		}

		for _, domainEntry := range resp.GetDomains() {
			if equalsDomain(domainEntry.GetDomain(), domain) {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return fmt.Errorf(
				"custom domain %q was not found for instance %q before deadline %s",
				domain,
				instanceID,
				deadline.Format(time.RFC3339Nano),
			)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

func equalsDomain(a, b string) bool {
	return strings.EqualFold(strings.TrimSuffix(a, "."), strings.TrimSuffix(b, "."))
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetInstanceClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get(InstanceIDVar).(string)
	domain := d.Get(DomainVar).(string)

	_, err = client.RemoveCustomDomain(ctx, &instance.RemoveCustomDomainRequest{
		InstanceId:   instanceID,
		CustomDomain: domain,
	})
	if err != nil {
		return diag.Errorf("failed to remove custom domain: %v", err)
	}

	return nil
}
