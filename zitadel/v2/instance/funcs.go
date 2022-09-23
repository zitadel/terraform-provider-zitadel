package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/system"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSystemClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveInstance(ctx, &system.RemoveInstanceRequest{InstanceId: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSystemClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	var email *system.AddInstanceRequest_Email
	if value := helper.GetIfSetString(d, ownerEmailVar); value != "" {
		email = &system.AddInstanceRequest_Email{
			Email:           value,
			IsEmailVerified: helper.GetIfSetBool(d, ownerIsEmailVerifiedVar),
		}
	}

	var profile *system.AddInstanceRequest_Profile
	valueFN := helper.GetIfSetString(d, ownerFirstNameVar)
	valueLN := helper.GetIfSetString(d, ownerLastNameVar)
	valuePL := helper.GetIfSetString(d, ownerPreferredLanguageVar)
	if valueFN != "" || valueLN != "" || valuePL != "" {
		profile = &system.AddInstanceRequest_Profile{
			FirstName:         valueFN,
			LastName:          valueLN,
			PreferredLanguage: valuePL,
		}
	}

	var pw *system.AddInstanceRequest_Password
	if value := helper.GetIfSetString(d, ownerPasswordVar); value != "" {
		pw = &system.AddInstanceRequest_Password{
			Password:               value,
			PasswordChangeRequired: helper.GetIfSetBool(d, ownerPasswordChangeRequiredVar),
		}
	}

	firstOrgName := helper.GetIfSetString(d, firstOrgNameVar)
	ownerUserName := helper.GetIfSetString(d, ownerUserNameVar)
	resp, err := client.AddInstance(ctx, &system.AddInstanceRequest{
		InstanceName:    d.Get(instanceNameVar).(string),
		FirstOrgName:    firstOrgName,
		CustomDomain:    helper.GetIfSetString(d, customDomainVar),
		OwnerUserName:   ownerUserName,
		OwnerEmail:      email,
		OwnerProfile:    profile,
		OwnerPassword:   pw,
		DefaultLanguage: helper.GetIfSetString(d, defaultLanguageVar),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.InstanceId)
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
	//TODO
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	//TODO
	return nil
}
