package pat

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemovePersonalAccessToken(ctx, &management.RemovePersonalAccessTokenRequest{
		UserId:  d.Get(userIDVar).(string),
		TokenId: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete PAT: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &management.AddPersonalAccessTokenRequest{
		UserId: d.Get(userIDVar).(string),
	}
	if expiration, ok := d.GetOk(expirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, expiration.(string))
		if err != nil {
			return diag.Errorf("failed to parse time: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddPersonalAccessToken(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tokenVar, resp.GetToken()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetTokenId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(userIDVar).(string)
	resp, err := client.GetPersonalAccessTokenByIDs(ctx, &management.GetPersonalAccessTokenByIDsRequest{
		UserId:  userID,
		TokenId: d.Id(),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get pat")
	}

	set := map[string]interface{}{
		expirationDateVar: resp.GetToken().GetExpirationDate().AsTime().Format(time.RFC3339),
		userIDVar:         userID,
		orgIDVar:          orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of pat: %v", k, err)
		}
	}
	d.SetId(resp.GetToken().GetId())
	return nil
}
