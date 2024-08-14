package pat

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemovePersonalAccessToken(helper.CtxWithOrgID(ctx, d), &management.RemovePersonalAccessTokenRequest{
		UserId:  d.Get(UserIDVar).(string),
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &management.AddPersonalAccessTokenRequest{
		UserId: d.Get(UserIDVar).(string),
	}
	if expiration, ok := d.GetOk(ExpirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, expiration.(string))
		if err != nil {
			return diag.Errorf("failed to parse time: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddPersonalAccessToken(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(TokenVar, resp.GetToken()); err != nil {
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

	orgID := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(UserIDVar).(string)
	resp, err := client.GetPersonalAccessTokenByIDs(helper.CtxWithOrgID(ctx, d), &management.GetPersonalAccessTokenByIDsRequest{
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
		ExpirationDateVar: resp.GetToken().GetExpirationDate().AsTime().Format(time.RFC3339),
		UserIDVar:         userID,
		helper.OrgIDVar:   orgID,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of pat: %v", k, err)
		}
	}
	d.SetId(resp.GetToken().GetId())
	return nil
}
