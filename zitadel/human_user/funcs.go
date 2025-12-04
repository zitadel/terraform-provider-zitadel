package human_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	firstName := d.Get(firstNameVar).(string)
	lastName := d.Get(lastNameVar).(string)
	email := d.Get(emailVar).(string)
	username := d.Get(UserNameVar).(string)
	orgID := d.Get(helper.OrgIDVar).(string)

	profile := &userv2.SetHumanProfile{
		GivenName:  firstName,
		FamilyName: lastName,
	}

	genderVal := userv2.Gender(userv2.Gender_value[d.Get(genderVar).(string)])
	profile.Gender = &genderVal

	prefLang := d.Get(preferredLanguageVar).(string)
	profile.PreferredLanguage = &prefLang

	if nickName, ok := d.GetOk(nickNameVar); ok && nickName.(string) != "" {
		nick := nickName.(string)
		profile.NickName = &nick
	}

	var displayName string
	if dn, ok := d.GetOk(DisplayNameVar); ok {
		displayName = dn.(string)
	} else {
		displayName = defaultDisplayName(firstName, lastName)
		if err := d.Set(DisplayNameVar, displayName); err != nil {
			return diag.Errorf("failed to set default display name for human user: %v", err)
		}
	}
	profile.DisplayName = &displayName

	humanEmail := &userv2.SetHumanEmail{
		Email: email,
	}

	if isVerified, ok := d.GetOk(isEmailVerifiedVar); ok && isVerified.(bool) {
		humanEmail.Verification = &userv2.SetHumanEmail_IsVerified{
			IsVerified: true,
		}
	}

	humanUser := &userv2.CreateUserRequest_Human{
		Profile: profile,
		Email:   humanEmail,
	}

	if phone, ok := d.GetOk(phoneVar); ok {
		phoneReq := &userv2.SetHumanPhone{
			Phone: phone.(string),
		}
		if isVerified, ok := d.GetOk(isPhoneVerifiedVar); ok && isVerified.(bool) {
			phoneReq.Verification = &userv2.SetHumanPhone_IsVerified{
				IsVerified: true,
			}
		}
		humanUser.Phone = phoneReq
	}

	if password, ok := d.GetOk(InitialPasswordVar); ok {
		pwd := password.(string)
		changeReq := !d.Get(initialSkipPasswordChange).(bool)
		humanUser.PasswordType = &userv2.CreateUserRequest_Human_Password{
			Password: &userv2.Password{
				Password:       pwd,
				ChangeRequired: changeReq,
			},
		}
	} else if hashedPassword, ok := d.GetOk(initialHashedPasswordVar); ok {
		hash := hashedPassword.(string)
		humanUser.PasswordType = &userv2.CreateUserRequest_Human_HashedPassword{
			HashedPassword: &userv2.HashedPassword{
				Hash: hash,
			},
		}
	}

	req := &userv2.CreateUserRequest{
		OrganizationId: orgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Human_{
			Human: humanUser,
		},
	}

	if userID, ok := d.GetOk(UserIDVar); ok {
		uid := userID.(string)
		req.UserId = &uid
	}

	respUser, err := client.CreateUser(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create human user: %v", err)
	}
	d.SetId(respUser.Id)

	return readFunc(false)(ctx, d, m)
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

	var humanUpdate *userv2.UpdateUserRequest_Human

	if d.HasChanges(firstNameVar, lastNameVar, nickNameVar, DisplayNameVar, preferredLanguageVar, genderVar) {
		profile := &userv2.UpdateUserRequest_Human_Profile{}

		if d.HasChange(firstNameVar) {
			firstName := d.Get(firstNameVar).(string)
			profile.GivenName = &firstName
		}
		if d.HasChange(lastNameVar) {
			lastName := d.Get(lastNameVar).(string)
			profile.FamilyName = &lastName
		}
		if d.HasChange(nickNameVar) {
			nickName := d.Get(nickNameVar).(string)
			profile.NickName = &nickName
		}
		if d.HasChange(DisplayNameVar) {
			displayName := d.Get(DisplayNameVar).(string)
			profile.DisplayName = &displayName
		}
		if d.HasChange(preferredLanguageVar) {
			prefLang := d.Get(preferredLanguageVar).(string)
			profile.PreferredLanguage = &prefLang
		}
		if d.HasChange(genderVar) {
			gender := userv2.Gender(userv2.Gender_value[d.Get(genderVar).(string)])
			profile.Gender = &gender
		}

		if humanUpdate == nil {
			humanUpdate = &userv2.UpdateUserRequest_Human{}
		}
		humanUpdate.Profile = profile
	}

	if d.HasChanges(emailVar, isEmailVerifiedVar) {
		emailUpdate := &userv2.SetHumanEmail{
			Email: d.Get(emailVar).(string),
		}
		if d.Get(isEmailVerifiedVar).(bool) {
			emailUpdate.Verification = &userv2.SetHumanEmail_IsVerified{
				IsVerified: true,
			}
		}

		if humanUpdate == nil {
			humanUpdate = &userv2.UpdateUserRequest_Human{}
		}
		humanUpdate.Email = emailUpdate
	}

	if d.HasChanges(phoneVar, isPhoneVerifiedVar) {
		phoneUpdate := &userv2.SetHumanPhone{
			Phone: d.Get(phoneVar).(string),
		}
		if d.Get(isPhoneVerifiedVar).(bool) {
			phoneUpdate.Verification = &userv2.SetHumanPhone_IsVerified{
				IsVerified: true,
			}
		}

		if humanUpdate == nil {
			humanUpdate = &userv2.UpdateUserRequest_Human{}
		}
		humanUpdate.Phone = phoneUpdate
	}

	if humanUpdate != nil {
		req.UserType = &userv2.UpdateUserRequest_Human_{
			Human: humanUpdate,
		}
	}

	if req.Username != nil || req.UserType != nil {
		_, err = client.UpdateUser(helper.CtxWithOrgID(ctx, d), req)
		if err != nil {
			return diag.Errorf("failed to update human user: %v", err)
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
		if !forDatasource {
			set[initialSkipPasswordChange] = false
		}

		if human := user.GetHuman(); human != nil {
			if profile := human.GetProfile(); profile != nil {
				set[firstNameVar] = profile.GetGivenName()
				set[lastNameVar] = profile.GetFamilyName()
				if dn := profile.DisplayName; dn != nil {
					set[DisplayNameVar] = *dn
				}
				if nn := profile.NickName; nn != nil {
					set[nickNameVar] = *nn
				}
				if pl := profile.PreferredLanguage; pl != nil {
					set[preferredLanguageVar] = *pl
				}
				if gender := profile.Gender; gender != nil {
					set[genderVar] = gender.String()
				}
			}
			if email := human.GetEmail(); email != nil {
				set[emailVar] = email.GetEmail()
				set[isEmailVerifiedVar] = email.GetIsVerified()
			}
			if phone := human.GetPhone(); phone != nil {
				set[phoneVar] = phone.GetPhone()
				set[isPhoneVerifiedVar] = phone.GetIsVerified()
			}
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of user: %v", k, err)
			}
		}
		d.SetId(user.GetUserId())
		return nil
	}
}

func defaultDisplayName(firstName, lastName string) string {
	return firstName + " " + lastName
}
