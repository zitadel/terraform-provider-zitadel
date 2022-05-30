package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	privacyPolicyOrgIdVar    = "org_id"
	privacyPolicyTOSLink     = "tos_link"
	privacyPolicyPrivacyLink = "privacy_link"
	privacyPolicyIsDefault   = "is_default"
	privacyPolicyHelpLink    = "help_link"
)

func GetPrivacyPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			privacyPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			privacyPolicyTOSLink: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			privacyPolicyPrivacyLink: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			privacyPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			privacyPolicyHelpLink: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
		CreateContext: createPrivacyPolicy,
		DeleteContext: deletePrivacyPolicy,
		ReadContext:   readPrivacyPolicy,
		UpdateContext: updatePrivacyPolicy,
	}
}

func deletePrivacyPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(privacyPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetPrivacyPolicyToDefault(ctx, &management2.ResetPrivacyPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset privacy policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func updatePrivacyPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(privacyPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomPrivacyPolicy(ctx, &management2.UpdateCustomPrivacyPolicyRequest{
		TosLink:     d.Get(privacyPolicyTOSLink).(string),
		PrivacyLink: d.Get(privacyPolicyPrivacyLink).(string),
		HelpLink:    d.Get(privacyPolicyHelpLink).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update privacy policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createPrivacyPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(privacyPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomPrivacyPolicy(ctx, &management2.AddCustomPrivacyPolicyRequest{
		TosLink:     d.Get(privacyPolicyTOSLink).(string),
		PrivacyLink: d.Get(privacyPolicyPrivacyLink).(string),
		HelpLink:    d.Get(privacyPolicyHelpLink).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create privacy policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readPrivacyPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(privacyPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPrivacyPolicy(ctx, &management2.GetPrivacyPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get privacy policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		privacyPolicyOrgIdVar:    policy.GetDetails().GetResourceOwner(),
		privacyPolicyIsDefault:   policy.GetIsDefault(),
		privacyPolicyTOSLink:     policy.GetTosLink(),
		privacyPolicyPrivacyLink: policy.GetPrivacyLink(),
		privacyPolicyHelpLink:    policy.GetHelpLink(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of privacy policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
