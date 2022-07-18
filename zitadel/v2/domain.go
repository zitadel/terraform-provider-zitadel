package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	domainOrgIdVar       = "org_id"
	domainNameVar        = "name"
	domainIsVerified     = "is_verified"
	domainIsPrimary      = "is_primary"
	domainValidationType = "validation_type"
)

func GetDomain() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			domainNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain",
				ForceNew:    true,
			},
			domainOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			domainIsVerified: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is domain verified",
			},
			domainIsPrimary: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is domain primary",
			},
			domainValidationType: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Validation type",
			},
		},
		ReadContext:   readDomain,
		CreateContext: createDomain,
		DeleteContext: deleteDomain,
	}
}

func deleteDomain(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(domainOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveOrgDomain(ctx, &management2.RemoveOrgDomainRequest{
		Domain: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete domain: %v", err)
	}
	return nil
}

func createDomain(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(domainOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(domainNameVar).(string)
	_, err = client.AddOrgDomain(ctx, &management2.AddOrgDomainRequest{
		Domain: name,
	})
	if err != nil {
		return diag.Errorf("failed to create domain: %v", err)
	}
	d.SetId(name)
	return nil
}

func readDomain(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(domainOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgDomains(ctx, &management2.ListOrgDomainsRequest{})
	if err != nil {
		return diag.Errorf("failed to read domain: %v", err)
	}

	set := map[string]interface{}{}
	domainStr := ""
	for i := range resp.Result {
		domain := resp.Result[i]
		if domain.GetDomainName() == d.Id() {
			domainStr = d.Id()
			set[domainNameVar] = domain.GetDomainName()
			set[domainOrgIdVar] = domain.GetOrgId()
			set[domainIsVerified] = domain.GetIsVerified()
			set[domainIsPrimary] = domain.GetIsPrimary()
			set[domainValidationType] = domain.GetValidationType().Number()
		}
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of domain: %v", k, err)
		}
	}
	d.SetId(domainStr)
	return nil
}
