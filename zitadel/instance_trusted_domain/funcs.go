package instance_trusted_domain

import (
	"context"
	"fmt"
	"strings"

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

	_, err = client.AddTrustedDomain(ctx, &instance.AddTrustedDomainRequest{
		InstanceId:    instanceID,
		TrustedDomain: domain,
	})
	if err != nil {
		return diag.Errorf("failed to add trusted domain: %v", err)
	}

	// Set ID as composite: instance_id/domain (or just domain if instance_id is empty)
	if instanceID != "" {
		d.SetId(fmt.Sprintf("%s/%s", instanceID, domain))
	} else {
		d.SetId(domain)
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

	// Parse composite ID
	var instanceID, domain string
	parts := strings.Split(d.Id(), "/")
	if len(parts) == 2 {
		instanceID = parts[0]
		domain = parts[1]
	} else if len(parts) == 1 {
		// Instance context, just domain
		domain = parts[0]
		instanceID = d.Get(InstanceIDVar).(string)
	} else {
		return diag.Errorf("invalid ID format, expected instance_id/domain or domain")
	}

	// List all trusted domains and find ours
	resp, err := client.ListTrustedDomains(ctx, &instance.ListTrustedDomainsRequest{
		InstanceId: instanceID,
	})
	if err != nil {
		return diag.Errorf("failed to list trusted domains: %v", err)
	}

	// Find the domain in the list
	found := false
	for _, d := range resp.GetTrustedDomain() {
		if d.GetDomain() == domain {
			found = true
			break
		}
	}

	if !found {
		// Domain was deleted outside of Terraform
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		DomainVar: domain,
	}
	if instanceID != "" {
		set[InstanceIDVar] = instanceID
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s: %v", k, err)
		}
	}

	return nil
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

	_, err = client.RemoveTrustedDomain(ctx, &instance.RemoveTrustedDomainRequest{
		InstanceId:    instanceID,
		TrustedDomain: domain,
	})
	if err != nil {
		return diag.Errorf("failed to remove trusted domain: %v", err)
	}

	return nil
}
