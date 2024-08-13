package machine_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUser(helper.CtxWithOrgID(ctx, d), &management.RemoveUserRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete user: %v", err)
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

	respUser, err := client.AddMachineUser(helper.CtxWithOrgID(ctx, d), &management.AddMachineUserRequest{
		UserName:        d.Get(UserNameVar).(string),
		Name:            d.Get(nameVar).(string),
		Description:     d.Get(DescriptionVar).(string),
		AccessTokenType: user.AccessTokenType(user.AccessTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
	})
	if err != nil {
		return diag.Errorf("failed to create machine user: %v", err)
	}
	d.SetId(respUser.UserId)

	if d.Get(WithSecretVar).(bool) {
		resp, err := client.GenerateMachineSecret(helper.CtxWithOrgID(ctx, d), &management.GenerateMachineSecretRequest{
			UserId: respUser.UserId,
		})
		if err != nil {
			return diag.Errorf("failed to generate machine user secret: %v", err)
		}
		if err := d.Set(clientIDVar, resp.GetClientId()); err != nil {
			return diag.Errorf("failed to set %s of user: %v", clientIDVar, err)
		}
		if err := d.Set(clientSecretVar, resp.GetClientSecret()); err != nil {
			return diag.Errorf("failed to set %s of user: %v", clientSecretVar, err)
		}
	}

	// To avoid diffs for terraform plan -refresh=false right after creation, we query and set the computed values.
	// The acceptance tests rely on this, too.
	return read(ctx, d, m)
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(UserNameVar) {
		_, err = client.UpdateUserName(helper.CtxWithOrgID(ctx, d), &management.UpdateUserNameRequest{
			UserId:   d.Id(),
			UserName: d.Get(UserNameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update username: %v", err)
		}
	}

	if d.HasChanges(nameVar, DescriptionVar, accessTokenTypeVar) {
		_, err := client.UpdateMachine(helper.CtxWithOrgID(ctx, d), &management.UpdateMachineRequest{
			UserId:          d.Id(),
			Name:            d.Get(nameVar).(string),
			Description:     d.Get(DescriptionVar).(string),
			AccessTokenType: user.AccessTokenType(user.AccessTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
		})
		if err != nil {
			return diag.Errorf("failed to update machine user: %v", err)
		}
	}

	if d.HasChange(WithSecretVar) {
		if d.Get(WithSecretVar).(bool) {
			resp, err := client.GenerateMachineSecret(helper.CtxWithOrgID(ctx, d), &management.GenerateMachineSecretRequest{
				UserId: d.Id(),
			})
			if err != nil {
				return diag.Errorf("failed to generate machine user secret: %v", err)
			}
			if err := d.Set(clientIDVar, resp.GetClientId()); err != nil {
				return diag.Errorf("failed to set %s of user: %v", clientIDVar, err)
			}
			if err := d.Set(clientSecretVar, resp.GetClientSecret()); err != nil {
				return diag.Errorf("failed to set %s of user: %v", clientSecretVar, err)
			}
		} else {
			_, err := client.RemoveMachineSecret(helper.CtxWithOrgID(ctx, d), &management.RemoveMachineSecretRequest{
				UserId: d.Id(),
			})
			if err != nil {
				return diag.Errorf("failed to remove machine user secret: %v", err)
			}
			if err := d.Set(clientIDVar, ""); err != nil {
				return diag.Errorf("failed to set %s of user: %v", clientIDVar, err)
			}
			if err := d.Set(clientSecretVar, ""); err != nil {
				return diag.Errorf("failed to set %s of user: %v", clientSecretVar, err)
			}
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.GetUserByID(helper.CtxWithOrgID(ctx, d), &management.GetUserByIDRequest{Id: helper.GetID(d, UserIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user")
	}

	user := respUser.GetUser()
	set := map[string]interface{}{
		helper.OrgIDVar:       user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState().String(),
		UserNameVar:           user.GetUserName(),
		loginNamesVar:         user.GetLoginNames(),
		preferredLoginNameVar: user.GetPreferredLoginName(),
	}
	if machine := user.GetMachine(); machine != nil {
		set[nameVar] = machine.GetName()
		set[DescriptionVar] = machine.GetDescription()
		set[accessTokenTypeVar] = machine.GetAccessTokenType().String()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of user: %v", k, err)
		}
	}
	d.SetId(user.GetId())
	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	userName := d.Get(UserNameVar).(string)
	userNameMethod := d.Get(userNameMethodVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListUsersRequest{}
	if userName != "" {
		req.Queries = append(req.Queries,
			&user.SearchQuery{
				Query: &user.SearchQuery_UserNameQuery{
					UserNameQuery: &user.UserNameQuery{
						UserName: userName,
						Method:   object.TextQueryMethod(object.TextQueryMethod_value[userNameMethod]),
					},
				},
			})
	}
	resp, err := client.ListUsers(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while getting user by username %s: %v", userName, err)
	}
	ids := make([]string, len(resp.Result))
	for i, res := range resp.Result {
		ids[i] = res.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(userIDsVar, ids))
}
