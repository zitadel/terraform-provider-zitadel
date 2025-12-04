package machine_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"
	userv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetUserV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteUser(helper.CtxWithOrgID(ctx, d), &userv2.DeleteUserRequest{
		UserId: d.Id(),
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

	client, err := helper.GetUserV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(nameVar).(string)
	username := d.Get(UserNameVar).(string)
	orgID := d.Get(helper.OrgIDVar).(string)

	machineUser := &userv2.CreateUserRequest_Machine{
		Name: name,
	}

	if description, ok := d.GetOk(DescriptionVar); ok && description.(string) != "" {
		desc := description.(string)
		machineUser.Description = &desc
	}

	req := &userv2.CreateUserRequest{
		OrganizationId: orgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Machine_{
			Machine: machineUser,
		},
	}

	if userID, ok := d.GetOk(UserIDVar); ok {
		uid := userID.(string)
		req.UserId = &uid
	}

	respUser, err := client.CreateUser(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create machine user: %v", err)
	}
	d.SetId(respUser.Id)

	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetUserV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &userv2.UpdateUserRequest{
		UserId: d.Id(),
	}

	if d.HasChange(UserNameVar) {
		username := d.Get(UserNameVar).(string)
		req.Username = &username
	}

	if d.HasChanges(nameVar, DescriptionVar) {
		machineUpdate := &userv2.UpdateUserRequest_Machine{}

		if d.HasChange(nameVar) {
			name := d.Get(nameVar).(string)
			machineUpdate.Name = &name
		}
		if d.HasChange(DescriptionVar) {
			desc := d.Get(DescriptionVar).(string)
			machineUpdate.Description = &desc
		}

		req.UserType = &userv2.UpdateUserRequest_Machine_{
			Machine: machineUpdate,
		}
	}

	if req.Username != nil || req.UserType != nil {
		_, err = client.UpdateUser(helper.CtxWithOrgID(ctx, d), req)
		if err != nil {
			return diag.Errorf("failed to update machine user: %v", err)
		}
	}

	if d.HasChange(WithSecretVar) {
		managementClient, err := helper.GetManagementClient(ctx, clientinfo)
		if err != nil {
			return diag.FromErr(err)
		}

		if d.Get(WithSecretVar).(bool) {
			resp, err := managementClient.GenerateMachineSecret(helper.CtxWithOrgID(ctx, d), &management.GenerateMachineSecretRequest{
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
			_, err := managementClient.RemoveMachineSecret(helper.CtxWithOrgID(ctx, d), &management.RemoveMachineSecretRequest{
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

	client, err := helper.GetUserV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.GetUserByID(helper.CtxWithOrgID(ctx, d), &userv2.GetUserByIDRequest{
		UserId: helper.GetID(d, UserIDVar),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user: %v", err)
	}

	user := respUser.GetUser()
	set := map[string]interface{}{
		helper.OrgIDVar:       user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState().String(),
		UserNameVar:           user.GetUsername(),
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
	d.SetId(user.GetUserId())
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
	d.SetId("-")
	return diag.FromErr(d.Set(userIDsVar, ids))
}
