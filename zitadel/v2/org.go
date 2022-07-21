package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	admin2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"strings"
)

const (
	nameVar = "name"
)

func OrgResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the org",
				ForceNew:    true,
			},
		},
		CreateContext: createOrg,
		DeleteContext: deleteOrg,
		ReadContext:   readOrg,
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

	if errDiag := readOrg(ctx, d, m); errDiag != nil {
		return errDiag
	}

	if d.Id() == "" {
		resp, err := client.AddOrg(ctx, &management2.AddOrgRequest{
			Name: d.Get(nameVar).(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resp.GetId())
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

	resp, err := client.ListOrgs(ctx, &admin2.ListOrgsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "found orgs", map[string]interface{}{
		"orglist": resp.Result,
	})

	//id := d.Get("id").(string)
	name := d.Get(nameVar).(string)
	tflog.Debug(ctx, "check if org is existing", map[string]interface{}{
		//	"id":  id,
		"org": name,
	})

	for i := range resp.Result {
		org := resp.Result[i]

		if strings.Compare(org.GetName(), name) == 0 {
			d.SetId(org.GetId())

			tflog.Debug(ctx, "found org", map[string]interface{}{
				"id":  d.Get("id"),
				"org": name,
			})
			return nil
		}
	}

	d.SetId("")
	tflog.Debug(ctx, "org not found", map[string]interface{}{
		"org": name,
	})
	return nil
}
