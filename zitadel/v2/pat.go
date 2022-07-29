package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	patOrgIDVar          = "org_id"
	patUserIDVar         = "user_id"
	patTokenVar          = "token"
	patExpirationDateVar = "expiration_date"
	timeFormat           = "2519-04-01T08:45:00.000000Z"
)

func GetPAT() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a personal access token of a user",
		Schema: map[string]*schema.Schema{
			patOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			patUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			patTokenVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the token",
				Sensitive:   true,
			},
			patExpirationDateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Expiration date of the token",
				ForceNew:    true,
			},
		},
		DeleteContext: deletePAT,
		CreateContext: createPAT,
		ReadContext:   readPAT,
	}
}

func deletePAT(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(patOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemovePersonalAccessToken(ctx, &management2.RemovePersonalAccessTokenRequest{
		UserId:  d.Get(patUserIDVar).(string),
		TokenId: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete PAT: %v", err)
	}
	return nil
}

func createPAT(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(patOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	t, err := time.Parse(timeFormat, d.Get(patExpirationDateVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddPersonalAccessToken(ctx, &management2.AddPersonalAccessTokenRequest{
		UserId:         d.Get(patUserIDVar).(string),
		ExpirationDate: timestamppb.New(t),
	})
	d.SetId(resp.GetTokenId())
	if err := d.Set(patTokenVar, resp.GetToken()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readPAT(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(patOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(patUserIDVar).(string)
	resp, err := client.GetPersonalAccessTokenByIDs(ctx, &management2.GetPersonalAccessTokenByIDsRequest{
		UserId:  userID,
		TokenId: d.Id(),
	})
	d.SetId(resp.GetToken().GetId())
	set := map[string]interface{}{
		patExpirationDateVar: resp.GetToken().GetExpirationDate().String(),
		patUserIDVar:         userID,
		patOrgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	if err := d.Set(patTokenVar, resp.GetToken()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
