package human_user

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

func legacyCreateHumanUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for create")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	firstName := d.Get(firstNameVar).(string)
	lastName := d.Get(lastNameVar).(string)
	email := d.Get(emailVar).(string)
	username := d.Get(UserNameVar).(string)

	profile := &management.ImportHumanUserRequest_Profile{
		FirstName:         firstName,
		LastName:          lastName,
		PreferredLanguage: d.Get(preferredLanguageVar).(string),
	}

	if nickName, ok := d.GetOk(nickNameVar); ok && nickName.(string) != "" {
		profile.NickName = nickName.(string)
	}

	if dn, ok := d.GetOk(DisplayNameVar); ok {
		profile.DisplayName = dn.(string)
	} else {
		profile.DisplayName = defaultDisplayName(firstName, lastName)
	}

	genderVal := user.Gender(user.Gender_value[d.Get(genderVar).(string)])
	profile.Gender = genderVal

	humanEmail := &management.ImportHumanUserRequest_Email{
		Email: email,
	}
	if isVerified, ok := d.GetOk(isEmailVerifiedVar); ok && isVerified.(bool) {
		humanEmail.IsEmailVerified = true
	}

	req := &management.ImportHumanUserRequest{
		UserName: username,
		Profile:  profile,
		Email:    humanEmail,
	}

	if phone, ok := d.GetOk(phoneVar); ok {
		phoneReq := &management.ImportHumanUserRequest_Phone{
			Phone: phone.(string),
		}
		if isVerified, ok := d.GetOk(isPhoneVerifiedVar); ok && isVerified.(bool) {
			phoneReq.IsPhoneVerified = true
		}
		req.Phone = phoneReq
	}

	if password, ok := d.GetOk(InitialPasswordVar); ok {
		req.Password = password.(string)
		req.PasswordChangeRequired = !d.Get(initialSkipPasswordChange).(bool)
	} else if hashedPassword, ok := d.GetOk(initialHashedPasswordVar); ok {
		req.HashedPassword = &management.ImportHumanUserRequest_HashedPassword{
			Value: hashedPassword.(string),
		}
	}

	resp, err := client.ImportHumanUser(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create human user: %v", err)
	}
	d.SetId(resp.GetUserId())

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

func legacyReadHumanUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo, forDatasource bool) diag.Diagnostics {
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
	if !forDatasource {
		set[initialSkipPasswordChange] = false
	}

	if human := u.GetHuman(); human != nil {
		if profile := human.GetProfile(); profile != nil {
			set[firstNameVar] = profile.GetFirstName()
			set[lastNameVar] = profile.GetLastName()
			set[DisplayNameVar] = profile.GetDisplayName()
			set[nickNameVar] = profile.GetNickName()
			set[preferredLanguageVar] = profile.GetPreferredLanguage()
			set[genderVar] = profile.GetGender().String()
		}
		if email := human.GetEmail(); email != nil {
			set[emailVar] = email.GetEmail()
			set[isEmailVerifiedVar] = email.GetIsEmailVerified()
		}
		if phone := human.GetPhone(); phone != nil {
			set[phoneVar] = phone.GetPhone()
			set[isPhoneVerifiedVar] = phone.GetIsPhoneVerified()
		}
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of user: %v", k, err)
		}
	}
	d.SetId(u.GetId())
	return nil
}

func legacyUpdateHumanUser(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for update")

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
			return diag.Errorf("failed to update human user username: %v", err)
		}
	}

	if d.HasChanges(firstNameVar, lastNameVar, nickNameVar, DisplayNameVar, preferredLanguageVar, genderVar) {
		_, err = client.UpdateHumanProfile(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanProfileRequest{
			UserId:            d.Id(),
			FirstName:         d.Get(firstNameVar).(string),
			LastName:          d.Get(lastNameVar).(string),
			NickName:          d.Get(nickNameVar).(string),
			DisplayName:       d.Get(DisplayNameVar).(string),
			PreferredLanguage: d.Get(preferredLanguageVar).(string),
			Gender:            user.Gender(user.Gender_value[d.Get(genderVar).(string)]),
		})
		if err != nil {
			return diag.Errorf("failed to update human user profile: %v", err)
		}
	}

	if d.HasChanges(emailVar, isEmailVerifiedVar) {
		oldEmail, newEmail := d.GetChange(emailVar)
		_, isVerifiedInConfig := d.GetOk(isEmailVerifiedVar)

		if oldEmail != newEmail || isVerifiedInConfig {
			_, err = client.UpdateHumanEmail(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanEmailRequest{
				UserId:          d.Id(),
				Email:           d.Get(emailVar).(string),
				IsEmailVerified: d.Get(isEmailVerifiedVar).(bool),
			})
			if err != nil {
				return diag.Errorf("failed to update human user email: %v", err)
			}
		}
	}

	if d.HasChanges(phoneVar, isPhoneVerifiedVar) {
		oldPhone, newPhone := d.GetChange(phoneVar)
		_, isVerifiedInConfig := d.GetOk(isPhoneVerifiedVar)

		if oldPhone != newPhone || isVerifiedInConfig {
			_, err = client.UpdateHumanPhone(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanPhoneRequest{
				UserId:          d.Id(),
				Phone:           d.Get(phoneVar).(string),
				IsPhoneVerified: d.Get(isPhoneVerifiedVar).(bool),
			})
			if err != nil {
				return diag.Errorf("failed to update human user phone: %v", err)
			}
		}
	}

	return nil
}

func legacyListHumanUsers(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
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
				Type: user.Type_TYPE_HUMAN,
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

	if firstName, ok := d.GetOk(firstNameVar); ok {
		firstNameMethod := d.Get(firstNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_FirstNameQuery{
				FirstNameQuery: &user.FirstNameQuery{
					FirstName: firstName.(string),
					Method:    object.TextQueryMethod(object.TextQueryMethod_value[firstNameMethod]),
				},
			},
		})
	}

	if lastName, ok := d.GetOk(lastNameVar); ok {
		lastNameMethod := d.Get(lastNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_LastNameQuery{
				LastNameQuery: &user.LastNameQuery{
					LastName: lastName.(string),
					Method:   object.TextQueryMethod(object.TextQueryMethod_value[lastNameMethod]),
				},
			},
		})
	}

	if nickName, ok := d.GetOk(nickNameVar); ok {
		nickNameMethod := d.Get(nickNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_NickNameQuery{
				NickNameQuery: &user.NickNameQuery{
					NickName: nickName.(string),
					Method:   object.TextQueryMethod(object.TextQueryMethod_value[nickNameMethod]),
				},
			},
		})
	}

	if displayName, ok := d.GetOk(DisplayNameVar); ok {
		displayNameMethod := d.Get(displayNameMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_DisplayNameQuery{
				DisplayNameQuery: &user.DisplayNameQuery{
					DisplayName: displayName.(string),
					Method:      object.TextQueryMethod(object.TextQueryMethod_value[displayNameMethod]),
				},
			},
		})
	}

	if email, ok := d.GetOk(emailVar); ok {
		emailMethod := d.Get(emailMethodVar).(string)
		queries = append(queries, &user.SearchQuery{
			Query: &user.SearchQuery_EmailQuery{
				EmailQuery: &user.EmailQuery{
					EmailAddress: email.(string),
					Method:       object.TextQueryMethod(object.TextQueryMethod_value[emailMethod]),
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
