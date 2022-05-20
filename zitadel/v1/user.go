package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v2 "github.com/zitadel/terraform-provider-zitadel/zitadel/v2"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
	"strconv"
)

const (
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

	isEmailVerifiedVar = "is_email_verified"
	emailVar           = "email"

	isPhoneVerifiedVar = "is_phone_verified"
	phoneVar           = "phone"

	machineNameVar = "name"
	descriptionVar = "description"
)

func GetUserDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			resourceOwnerVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			userStateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "State of the user",
			},
			userNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username",
			},
			loginNamesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Loginnames",
			},
			preferredLoginNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Preferred login name",
			},
			typeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the user",
			},

			firstNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "First name of the user",
			},
			lastNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last name of the user",
			},
			nickNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Nick name of the user",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DIsplay name of the user",
			},
			preferredLanguageVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Preferred language of the user",
			},
			genderVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Gender of the user",
			},

			emailVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the user",
			},
			isEmailVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the email verified of the user",
			},

			phoneVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Phone of the user",
			},
			isPhoneVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the phone verified of the user",
			},

			machineNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the machine user",
			},
			descriptionVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the user",
			},
		},
	}
}

func readUser(ctx context.Context, d *schema.ResourceData, m interface{}, info *ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	client, err := getManagementClient(info, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	respUser, err := client.GetUserByID(ctx, &management2.GetUserByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to get list of users: %v", err)
	}

	user := respUser.GetUser()
	userType := ""
	loginNames := []string{}
	for _, v := range user.GetLoginNames() {
		loginNames = append(loginNames, v)
	}
	set := map[string]interface{}{
		resourceOwnerVar:      user.GetDetails().GetResourceOwner(),
		userStateVar:          user.GetState(),
		userNameVar:           user.GetUserName(),
		loginNamesVar:         loginNames,
		preferredLoginNameVar: user.GetPreferredLoginName(),
		typeVar:               userType,
	}
	if human := user.GetHuman(); human != nil {
		set[typeVar] = v2.HumanUser
		if profile := human.GetProfile(); profile != nil {
			set[firstNameVar] = profile.GetFirstName()
			set[lastNameVar] = profile.GetLastName()
			set[displayNameVar] = profile.GetDisplayName()
			set[nickNameVar] = profile.GetNickName()
			set[preferredLanguageVar] = profile.GetPreferredLanguage()
			if gender := profile.GetGender().String(); gender != "" {
				genderInt, err := strconv.Atoi(gender)
				if err != nil {
					return diag.Errorf("failed to parse gender: %v", err)
				}
				set[genderVar] = genderInt
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
		set[typeVar] = v2.MachineUser
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
