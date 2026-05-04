package machine_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func legacyCreateMachineUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for create")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(nameVar).(string)
	username := d.Get(UserNameVar).(string)

	req := &management.AddMachineUserRequest{
		UserName:    username,
		Name:        name,
		Description: d.Get(DescriptionVar).(string),
	}

	if desiredType := d.Get(accessTokenTypeVar).(string); desiredType != defaultAccessTokenType {
		req.AccessTokenType = user.AccessTokenType(user.AccessTokenType_value[desiredType])
	}

	resp, err := client.AddMachineUser(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create machine user: %v", err)
	}
	d.SetId(resp.GetUserId())

	if d.Get(WithSecretVar).(bool) {
		secretResp, err := client.GenerateMachineSecret(helper.CtxWithOrgID(ctx, d), &management.GenerateMachineSecretRequest{
			UserId: resp.GetUserId(),
		})
		if err != nil {
			return diag.Errorf("failed to generate machine user secret: %v", err)
		}
		if err := d.Set(clientIDVar, secretResp.GetClientId()); err != nil {
			return diag.Errorf("failed to set %s of user: %v", clientIDVar, err)
		}
		if err := d.Set(clientSecretVar, secretResp.GetClientSecret()); err != nil {
			return diag.Errorf("failed to set %s of user: %v", clientSecretVar, err)
		}
	}

	return nil
}

func legacyDeleteUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for delete")

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

func legacyReadMachineUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for read")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetUserByID(helper.CtxWithOrgID(ctx, d), &management.GetUserByIDRequest{
		Id: helper.GetID(d, UserIDVar),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user: %v", err)
	}

	u := resp.GetUser()
	set := map[string]interface{}{
		helper.OrgIDVar:       u.GetDetails().GetResourceOwner(),
		userStateVar:          u.GetState().String(),
		UserNameVar:           u.GetUserName(),
		loginNamesVar:         u.GetLoginNames(),
		preferredLoginNameVar: u.GetPreferredLoginName(),
	}
	if machine := u.GetMachine(); machine != nil {
		set[nameVar] = machine.GetName()
		set[DescriptionVar] = machine.GetDescription()
		set[accessTokenTypeVar] = machine.GetAccessTokenType().String()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of user: %v", k, err)
		}
	}
	d.SetId(u.GetId())
	return nil
}

func legacyUpdateMachineUsername(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for update username")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserName(helper.CtxWithOrgID(ctx, d), &management.UpdateUserNameRequest{
		UserId:   d.Id(),
		UserName: d.Get(UserNameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update machine user: %v", err)
	}
	return nil
}

func legacyListMachineUsers(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for list")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &management.ListUsersRequest{}
	var queries []*user.SearchQuery

	queries = append(queries, &user.SearchQuery{
		Query: &user.SearchQuery_TypeQuery{
			TypeQuery: &user.TypeQuery{
				Type: user.Type_TYPE_MACHINE,
			},
		},
	})

	if userName, ok := d.GetOk(UserNameVar); ok {
		userNameMethod := d.Get(userNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_UserNameQuery{
				UserNameQuery: &user.UserNameQuery{
					UserName: userName.(string),
					Method:   object.TextQueryMethod(object.TextQueryMethod_value[userNameMethod]),
				},
			},
		})
	}

	if loginName, ok := d.GetOk(loginNameVar); ok {
		loginNameMethod := d.Get(loginNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_LoginNameQuery{
				LoginNameQuery: &user.LoginNameQuery{
					LoginName: loginName.(string),
					Method:    object.TextQueryMethod(object.TextQueryMethod_value[loginNameMethod]),
				},
			},
		})
	}

	req.Queries = queries

	resp, err := client.ListUsers(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while listing users: %v", err)
	}

	ids := make([]string, len(resp.Result))
	for i, res := range resp.Result {
		ids[i] = res.GetId()
	}

	d.SetId("-")
	return diag.FromErr(d.Set(userIDsVar, ids))
}
