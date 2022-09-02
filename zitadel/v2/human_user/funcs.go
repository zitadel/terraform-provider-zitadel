package human_user

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
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	firstName := d.Get(firstNameVar).(string)
	lastName := d.Get(lastNameVar).(string)
	addUser := &management.AddHumanUserRequest{
		UserName: d.Get(userNameVar).(string),
		Profile: &management.AddHumanUserRequest_Profile{
			FirstName: firstName,
			LastName:  lastName,
		},
	}

	nickname := d.Get(nickNameVar).(string)
	if nickname != "" {
		addUser.Profile.NickName = nickname
	}

	displayname := d.Get(displayNameVar).(string)
	if displayname != "" {
		addUser.Profile.DisplayName = displayname
	} else {
		if err := d.Set(displayNameVar, defaultDisplayName(firstName, lastName)); err != nil {
			return diag.Errorf("failed to set default display name for human user: %v", err)
		}
	}

	prefLang := d.Get(preferredLanguageVar).(string)
	if prefLang != "" {
		addUser.Profile.PreferredLanguage = prefLang
	} else {
		if err := d.Set(preferredLanguageVar, defaultPreferredLanguage); err != nil {
			return diag.Errorf("failed to set default preferred language for human user: %v", err)
		}
	}

	gender := d.Get(genderVar).(string)
	if gender != "" {
		addUser.Profile.Gender = user.Gender(user.Gender_value[gender])
	} else {
		if err := d.Set(genderVar, defaultGenderString); err != nil {
			return diag.Errorf("failed to set default gender for human user: %v", err)
		}
	}

	pwd := d.Get(initialPasswordVar).(string)
	if pwd != "" {
		addUser.InitialPassword = pwd
	}

	email := d.Get(emailVar).(string)
	if email != "" {
		isVerified := d.Get(isEmailVerifiedVar)
		addUser.Email = &management.AddHumanUserRequest_Email{
			Email:           email,
			IsEmailVerified: false,
		}
		if isVerified != nil {
			addUser.Email.IsEmailVerified = isVerified.(bool)
		}
	}

	phone := d.Get(phoneVar).(string)
	if phone != "" {
		isVerified := d.Get(isPhoneVerifiedVar)
		addUser.Phone = &management.AddHumanUserRequest_Phone{
			Phone:           phone,
			IsPhoneVerified: false,
		}
		if isVerified != nil {
			addUser.Phone.IsPhoneVerified = isVerified.(bool)
		}
	}

	respUser, err := client.AddHumanUser(ctx, addUser)
	if err != nil {
		return diag.Errorf("failed to create human user: %v", err)
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

	currentUser, err := client.GetUserByID(ctx, &management.GetUserByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	username := d.Get(userNameVar).(string)
	if currentUser.GetUser().GetUserName() != username {
		_, err = client.UpdateUserName(ctx, &management.UpdateUserNameRequest{
			UserId:   d.Id(),
			UserName: username,
		})
		if err != nil {
			return diag.Errorf("failed to update username: %v", err)
		}
	}

	nickname := d.Get(nickNameVar)
	displayname := d.Get(displayNameVar)
	prefLang := d.Get(preferredLanguageVar)
	gender := d.Get(genderVar)
	email := d.Get(emailVar)
	emailVerfied := d.Get(isEmailVerifiedVar)
	phone := d.Get(phoneVar)
	phoneVerified := d.Get(isPhoneVerifiedVar)

	currentHuman := currentUser.GetUser().GetHuman()
	if currentHuman.GetProfile().GetFirstName() != d.Get(firstNameVar).(string) ||
		currentHuman.GetProfile().GetLastName() != d.Get(lastNameVar).(string) ||
		(nickname != nil && currentHuman.GetProfile().GetNickName() != nickname.(string)) ||
		(displayname != nil && currentHuman.GetProfile().GetDisplayName() != displayname.(string)) ||
		(prefLang != nil && currentHuman.GetProfile().GetPreferredLanguage() != prefLang.(string)) ||
		(gender != nil && currentHuman.GetProfile().GetGender().String() != gender.(string)) {

		_, err := client.UpdateHumanProfile(ctx, &management.UpdateHumanProfileRequest{
			UserId:            d.Id(),
			FirstName:         d.Get(firstNameVar).(string),
			LastName:          d.Get(lastNameVar).(string),
			NickName:          d.Get(nickNameVar).(string),
			DisplayName:       d.Get(displayNameVar).(string),
			PreferredLanguage: d.Get(preferredLanguageVar).(string),
			Gender:            user.Gender(user.Gender_value[gender.(string)]),
		})
		if err != nil {
			return diag.Errorf("failed to update human profile: %v", err)
		}
	}

	if currentHuman.GetEmail().GetEmail() != email.(string) || currentHuman.GetEmail().GetIsEmailVerified() != emailVerfied.(bool) {
		_, err = client.UpdateHumanEmail(ctx, &management.UpdateHumanEmailRequest{
			UserId:          d.Id(),
			Email:           email.(string),
			IsEmailVerified: emailVerfied.(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human email: %v", err)
		}
	}

	if currentHuman.GetPhone().GetPhone() != phone.(string) || currentHuman.GetPhone().GetIsPhoneVerified() != phoneVerified.(bool) {
		_, err = client.UpdateHumanPhone(ctx, &management.UpdateHumanPhoneRequest{
			UserId:          d.Id(),
			Phone:           phone.(string),
			IsPhoneVerified: phoneVerified.(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human phone: %v", err)
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
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get list of users: %v", err)
	}

	user := respUser.GetUser()
	set := map[string]interface{}{
		orgIDVar:              user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState().String(),
		userNameVar:           user.GetUserName(),
		loginNamesVar:         user.GetLoginNames(),
		preferredLoginNameVar: user.GetPreferredLoginName(),
	}

	if human := user.GetHuman(); human != nil {
		if profile := human.GetProfile(); profile != nil {
			set[firstNameVar] = profile.GetFirstName()
			set[lastNameVar] = profile.GetLastName()
			set[displayNameVar] = profile.GetDisplayName()
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

func defaultDisplayName(firstName, lastName string) string {
	return firstName + " " + lastName
}
