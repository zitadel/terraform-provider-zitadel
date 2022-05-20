package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	orgVar     = "org"
	nameVar    = "name"
	issuerVar  = "issuer"
	addressVar = "address"
	projectVar = "project"
	tokenVar   = "token"
	usersVar   = "users"
)

func GetOrgDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			orgVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the organization",
			},
			nameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the organization",
			},
			issuerVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ISSUER", ""),
			},
			addressVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ADDRESS", ""),
			},
			projectVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROJECT", ""),
			},
			tokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},
			usersVar: {
				Type:        schema.TypeSet,
				Elem:        GetUserDatasource(),
				Computed:    true,
				Description: "List of ID of users in organization",
			},
		},
		ReadContext: readOrg,
	}
}

func readOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, err := GetClientInfo(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := getManagementClient(clientinfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetMyOrg(ctx, &management2.GetMyOrgRequest{})
	if err != nil {
		return diag.Errorf("failed to get org: %v", err)
	}
	id := resp.GetOrg().GetId()
	d.SetId(id)
	name := resp.GetOrg().GetName()

	tflog.Debug(ctx, "found org", map[string]interface{}{
		"id":   id,
		"name": name,
	})

	users := make([]*schema.ResourceData, 0)
	respUsers, err := client.ListUsers(ctx, &management2.ListUsersRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of users: %v", err)
	}

	for i := range respUsers.Result {
		user := respUsers.Result[i]

		userdata := &schema.ResourceData{}
		userdata.SetId(user.GetId())
		if errDiag := readUser(ctx, userdata, m, clientinfo); errDiag != nil {
			return errDiag
		}

		users = append(users, userdata)
	}
	if err := d.Set(usersVar, users); err != nil {
		return diag.Errorf("failed to set list of users: %v", err)
	}

	if err := d.Set(nameVar, name); err != nil {
		return diag.Errorf("failed to set org name: %v", err)
	}
	if err := d.Set(orgVar, id); err != nil {
		return diag.Errorf("failed to set org: %v", err)
	}
	d.SetId(id)
	return nil
}
