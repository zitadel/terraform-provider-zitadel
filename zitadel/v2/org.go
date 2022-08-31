package v2

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	nameVar = "name"
)

func OrgResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the org",
			},
		},
		CreateContext: createOrg,
		DeleteContext: deleteOrg,
		ReadContext:   readOrg,
		UpdateContext: updateOrg,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func deleteOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	/*client, ok := m.(*management.Client)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	_, err := client.DeactivateOrg(ctx, &management2.DeactivateOrgRequest{})
	if err != nil {
		return diag.FromErr(err)
	}*/
	return nil
}

func createOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOrg(ctx, &management.AddOrgRequest{
		Name: d.Get(nameVar).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetId())

	return nil
}

func updateOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(actionOrgId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrg(ctx, &management.UpdateOrgRequest{
		Name: d.Get(nameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update org: %v", err)
	}
	return nil
}

func readOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgs(ctx, &admin.ListOrgsRequest{})
	if err != nil {
		return diag.Errorf("error while listing orgs: %v", err)
	}
	tflog.Debug(ctx, "found orgs", map[string]interface{}{
		"orglist": resp.Result,
	})

	orgID := d.Id()
	tflog.Debug(ctx, "check if org is existing", map[string]interface{}{
		"id": orgID,
	})

	for i := range resp.Result {
		org := resp.Result[i]
		if strings.Compare(org.GetId(), orgID) == 0 {
			d.SetId(orgID)
			tflog.Debug(ctx, "found org", map[string]interface{}{
				"id": orgID,
			})
			if err := d.Set(nameVar, org.GetName()); err != nil {
				return diag.Errorf("failed to set %s of org: %v", nameVar, err)
			}
			return nil
		}
	}

	d.SetId("")
	tflog.Debug(ctx, "org not found", map[string]interface{}{})
	return nil
}
