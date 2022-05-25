package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/pkg/client/zitadel/user"
)

const (
	idVar    = "id"
	oldIDVar = "old_id"

	resourceOwnerVar      = "resource_owner"
	userStateVar          = "state"
	userNameVar           = "user_name"
	loginNamesVar         = "login_names"
	preferredLoginNameVar = "preferred_login_name"
	typeVar               = "type"

	firstNameVar         = "first_name"
	lastNameVar          = "last_name"
	nickNameVar          = "nick_name"
	displayNameVar       = "display_name"
	preferredLanguageVar = "preferred_language"
	genderVar            = "gender"

	isPhoneVerifiedVar = "is_phone_verified"
	emailVar           = "email"

	isEmailVerifiedVar = "is_email_verified"
	phoneVar           = "phone"

	machineNameVar = "name"
	descriptionVar = "description"

	HumanUser   = "human"
	MachineUser = "maschine"
)

func GetUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			idVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			oldIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Old ID of the user",
				ForceNew:    true,
			},
			resourceOwnerVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			userStateVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "State of the user",
			},
			userNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username",
			},
			loginNamesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Loginnames",
				ForceNew:    true,
			},
			preferredLoginNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Preferred login name",
				ForceNew:    true,
			},
			typeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the user",
				ForceNew:    true,
			},

			firstNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name of the user",
			},
			lastNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Last name of the user",
			},
			nickNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Nick name of the user",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DIsplay name of the user",
			},
			preferredLanguageVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred language of the user",
			},
			genderVar: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Gender of the user",
			},

			emailVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email of the user",
			},
			isEmailVerifiedVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the email verified of the user",
			},

			phoneVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone of the user",
			},
			isPhoneVerifiedVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the phone verified of the user",
			},

			machineNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the machine user",
			},
			descriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the user",
			},
		},
		ReadContext:   readUser,
		CreateContext: createUser,
		DeleteContext: deleteUser,
		UpdateContext: updateUser,
	}
}

func deleteUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(resourceOwnerVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUser(ctx, &management2.RemoveUserRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func createUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	userType := d.Get(typeVar).(string)
	switch userType {
	case HumanUser:
		gender := d.Get(genderVar).(user.Gender)
		respUser, err := client.AddHumanUser(ctx, &management2.AddHumanUserRequest{
			UserName: d.Get(userNameVar).(string),
			Profile: &management2.AddHumanUserRequest_Profile{
				FirstName:         d.Get(firstNameVar).(string),
				LastName:          d.Get(lastNameVar).(string),
				NickName:          d.Get(nickNameVar).(string),
				DisplayName:       d.Get(displayNameVar).(string),
				PreferredLanguage: d.Get(preferredLanguageVar).(string),
				Gender:            gender,
			},
			Email: &management2.AddHumanUserRequest_Email{
				Email:           d.Get(emailVar).(string),
				IsEmailVerified: d.Get(isEmailVerifiedVar).(bool),
			},
			Phone: &management2.AddHumanUserRequest_Phone{
				Phone:           d.Get(phoneVar).(string),
				IsPhoneVerified: d.Get(isPhoneVerifiedVar).(bool),
			},
		})
		if err != nil {
			return diag.Errorf("failed to create human user: %v", err)
		}
		d.SetId(respUser.UserId)
	case MachineUser:
		respUser, err := client.AddMachineUser(ctx, &management2.AddMachineUserRequest{
			UserName:    d.Get(userNameVar).(string),
			Name:        d.Get(machineNameVar).(string),
			Description: d.Get(descriptionVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to create machine user: %v", err)
		}
		d.SetId(respUser.UserId)
	}
	return nil
}

func updateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserName(ctx, &management2.UpdateUserNameRequest{
		UserId:   d.Id(),
		UserName: d.Get(userNameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update username: %v", err)
	}

	userType := d.Get(typeVar).(string)
	switch userType {
	case HumanUser:
		gender := d.Get(genderVar).(user.Gender)
		_, err := client.UpdateHumanProfile(ctx, &management2.UpdateHumanProfileRequest{
			UserId:            d.Id(),
			FirstName:         d.Get(firstNameVar).(string),
			LastName:          d.Get(lastNameVar).(string),
			NickName:          d.Get(nickNameVar).(string),
			DisplayName:       d.Get(displayNameVar).(string),
			PreferredLanguage: d.Get(preferredLanguageVar).(string),
			Gender:            gender,
		})
		if err != nil {
			return diag.Errorf("failed to update human profile: %v", err)
		}
		_, err = client.UpdateHumanEmail(ctx, &management2.UpdateHumanEmailRequest{
			UserId:          d.Id(),
			Email:           d.Get(emailVar).(string),
			IsEmailVerified: d.Get(isEmailVerifiedVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human email: %v", err)
		}

		_, err = client.UpdateHumanPhone(ctx, &management2.UpdateHumanPhoneRequest{
			UserId:          d.Id(),
			Phone:           d.Get(phoneVar).(string),
			IsPhoneVerified: d.Get(isPhoneVerifiedVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update human phone: %v", err)
		}
	case MachineUser:
		_, err := client.UpdateMachine(ctx, &management2.UpdateMachineRequest{
			UserId:      d.Id(),
			Name:        d.Get(machineNameVar).(string),
			Description: d.Get(descriptionVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update machine user: %v", err)
		}
	}
	return nil
}

func readUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.GetUserByID(ctx, &management2.GetUserByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to get list of users: %v", err)
	}

	user := respUser.GetUser()
	set := map[string]interface{}{
		idVar:                 user.GetId(),
		resourceOwnerVar:      user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState(),
		userNameVar:           user.GetUserName(),
		loginNamesVar:         user.GetLoginNames(),
		preferredLoginNameVar: user.GetPreferredLoginName(),
		typeVar:               user.GetType(),
	}
	if human := user.GetHuman(); human != nil {
		set[typeVar] = HumanUser
		if profile := human.GetProfile(); profile != nil {
			set[firstNameVar] = profile.GetFirstName()
			set[lastNameVar] = profile.GetLastName()
			set[displayNameVar] = profile.GetDisplayName()
			set[nickNameVar] = profile.GetNickName()
			set[preferredLanguageVar] = profile.GetPreferredLanguage()
			if gender := profile.GetGender().Number(); gender != 0 {
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
	} else if machine := user.GetMachine(); machine != nil {
		set[typeVar] = MachineUser
		set[machineNameVar] = machine.GetName()
		set[descriptionVar] = machine.GetDescription()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of user: %v", k, err)
		}
	}
	d.SetId(user.GetId())
	return nil
}
