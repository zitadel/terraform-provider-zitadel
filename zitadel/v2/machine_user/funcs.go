package machine_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUser(ctx, &management.RemoveUserRequest{
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.AddMachineUser(ctx, &management.AddMachineUserRequest{
		UserName:        d.Get(userNameVar).(string),
		Name:            d.Get(nameVar).(string),
		Description:     d.Get(descriptionVar).(string),
		AccessTokenType: user.AccessTokenType(user.AccessTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
	})
	if err != nil {
		return diag.Errorf("failed to create machine user: %v", err)
	}
	d.SetId(respUser.UserId)
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(userNameVar) {
		_, err = client.UpdateUserName(ctx, &management.UpdateUserNameRequest{
			UserId:   d.Id(),
			UserName: d.Get(userNameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update username: %v", err)
		}
	}

	if d.HasChanges(nameVar, descriptionVar, accessTokenTypeVar) {
		_, err := client.UpdateMachine(ctx, &management.UpdateMachineRequest{
			UserId:          d.Id(),
			Name:            d.Get(nameVar).(string),
			Description:     d.Get(descriptionVar).(string),
			AccessTokenType: user.AccessTokenType(user.AccessTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
		})
		if err != nil {
			return diag.Errorf("failed to update machine user: %v", err)
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.GetUserByID(ctx, &management.GetUserByIDRequest{Id: helper.GetID(d, userIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user")
	}

	user := respUser.GetUser()
	set := map[string]interface{}{
		orgIDVar:              user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState().String(),
		userNameVar:           user.GetUserName(),
		loginNamesVar:         user.GetLoginNames(),
		preferredLoginNameVar: user.GetPreferredLoginName(),
	}
	if machine := user.GetMachine(); machine != nil {
		set[nameVar] = machine.GetName()
		set[descriptionVar] = machine.GetDescription()
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
