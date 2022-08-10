package v2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	patOrgIDVar          = "org_id"
	patUserIDVar         = "user_id"
	patTokenVar          = "token"
	patExpirationDateVar = "expiration_date"
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

	t, err := time.Parse(time.RFC3339, d.Get(patExpirationDateVar).(string))
	if err != nil {
		return diag.Errorf("failed to parse time: %v", err)
	}

	resp, err := client.AddPersonalAccessToken(ctx, &management2.AddPersonalAccessTokenRequest{
		UserId:         d.Get(patUserIDVar).(string),
		ExpirationDate: timestamppb.New(t),
	})

	if err := d.Set(patTokenVar, resp.GetToken()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetTokenId())
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
	if err != nil {
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		patExpirationDateVar: resp.GetToken().GetExpirationDate().AsTime().Format(time.RFC3339),
		patUserIDVar:         userID,
		patOrgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	d.SetId(resp.GetToken().GetId())
	return nil
}
