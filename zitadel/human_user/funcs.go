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
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	firstName := d.Get(firstNameVar).(string)
	lastName := d.Get(lastNameVar).(string)
	importUser := &management.ImportHumanUserRequest{
		UserName: d.Get(UserNameVar).(string),
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName:         firstName,
			LastName:          lastName,
			Gender:            user.Gender(user.Gender_value[d.Get(genderVar).(string)]),
			PreferredLanguage: d.Get(preferredLanguageVar).(string),
			NickName:          d.Get(nickNameVar).(string),
		},
		Password:               d.Get(InitialPasswordVar).(string),
		PasswordChangeRequired: !d.Get(initialSkipPasswordChange).(bool),
	}

	if hashedPassword, ok := d.GetOk(initialHashedPasswordVar); ok {
		importUser.HashedPassword = &management.ImportHumanUserRequest_HashedPassword{
			Value: hashedPassword.(string),
		}
	}

	if displayname, ok := d.GetOk(DisplayNameVar); ok {
		importUser.Profile.DisplayName = displayname.(string)
	} else {
		if err := d.Set(DisplayNameVar, defaultDisplayName(firstName, lastName)); err != nil {
			return diag.Errorf("failed to set default display name for human user: %v", err)
		}
	}

	if email, ok := d.GetOk(emailVar); ok {
		isVerified, isVerifiedOk := d.GetOk(isEmailVerifiedVar)
		importUser.Email = &management.ImportHumanUserRequest_Email{
			Email:           email.(string),
			IsEmailVerified: false,
		}
		if isVerifiedOk {
			importUser.Email.IsEmailVerified = isVerified.(bool)
		}
	}

	if phone, ok := d.GetOk(phoneVar); ok {
		isVerified, isVerifiedOk := d.GetOk(isPhoneVerifiedVar)
		importUser.Phone = &management.ImportHumanUserRequest_Phone{
			Phone:           phone.(string),
			IsPhoneVerified: false,
		}
		if isVerifiedOk {
			importUser.Phone.IsPhoneVerified = isVerified.(bool)
		}
	}

	respUser, err := client.ImportHumanUser(helper.CtxWithOrgID(ctx, d), importUser)
	if err != nil {
		return diag.Errorf("failed to create human user: %v", err)
	}
	d.SetId(respUser.UserId)
	// To avoid diffs for terraform plan -refresh=false right after creation, we query and set the computed values.
	// The acceptance tests rely on this, too.
	return readFunc(false)(ctx, d, m)
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

	if d.HasChanges(firstNameVar, lastNameVar, nickNameVar, DisplayNameVar, preferredLanguageVar, genderVar) {
		_, err := client.UpdateHumanProfile(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanProfileRequest{
			UserId:            d.Id(),
			FirstName:         d.Get(firstNameVar).(string),
			LastName:          d.Get(lastNameVar).(string),
			NickName:          d.Get(nickNameVar).(string),
			DisplayName:       d.Get(DisplayNameVar).(string),
			PreferredLanguage: d.Get(preferredLanguageVar).(string),
			Gender:            user.Gender(user.Gender_value[d.Get(genderVar).(string)]),
		})
		if err != nil {
			return diag.Errorf("failed to update human profile: %v", err)
		}
	}

	if d.HasChanges(emailVar, isEmailVerifiedVar) {
		_, err = client.UpdateHumanEmail(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanEmailRequest{
			UserId:          d.Id(),
			Email:           d.Get(emailVar).(string),
			IsEmailVerified: d.Get(isEmailVerifiedVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human email: %v", err)
		}
	}

	if d.HasChanges(phoneVar, isPhoneVerifiedVar) {
		_, err = client.UpdateHumanPhone(helper.CtxWithOrgID(ctx, d), &management.UpdateHumanPhoneRequest{
			UserId:          d.Id(),
			Phone:           d.Get(phoneVar).(string),
			IsPhoneVerified: d.Get(isPhoneVerifiedVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human phone: %v", err)
		}
	}
	return nil
}

func readFunc(forDatasource bool) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			return diag.Errorf("failed to get user: %v", err)
		}

		user := respUser.GetUser()
		set := map[string]interface{}{
			helper.OrgIDVar:       user.GetDetails().GetResourceOwner(),
			userStateVar:          user.GetState().String(),
			UserNameVar:           user.GetUserName(),
			loginNamesVar:         user.GetLoginNames(),
			preferredLoginNameVar: user.GetPreferredLoginName(),
		}
		if !forDatasource {
			// This will be ignored using the CustomizeDiff function.
			// However, we should explicitly set it to true or false so that importing a user doesn't produce an immediate plan diff.
			set[initialSkipPasswordChange] = false
		}

		if human := user.GetHuman(); human != nil {
			if profile := human.GetProfile(); profile != nil {
				set[firstNameVar] = profile.GetFirstName()
				set[lastNameVar] = profile.GetLastName()
				set[DisplayNameVar] = profile.GetDisplayName()
				set[nickNameVar] = profile.GetNickName()
				set[preferredLanguageVar] = profile.GetPreferredLanguage()
				if gender := profile.GetGender().String(); gender != "" {
					set[genderVar] = gender
				}
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
		d.SetId(user.GetId())
		return nil
	}
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

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
		ids[i] = res.Id
	}

	d.SetId("-")
	return diag.FromErr(d.Set(userIDsVar, ids))
}

func defaultDisplayName(firstName, lastName string) string {
	return firstName + " " + lastName
}
