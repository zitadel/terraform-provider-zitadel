package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	lockoutPolicyOrgIdVar            = "org_id"
	lockoutPolicyMaxPasswordAttempts = "max_password_attempts"
)

func GetLockoutPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom lockout policy of an organization.",
		Schema: map[string]*schema.Schema{
			lockoutPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			lockoutPolicyMaxPasswordAttempts: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correct or the password is reset.",
			},
		},
		DeleteContext: deleteLockoutPolicy,
		CreateContext: createLockoutPolicy,
		UpdateContext: updateLockoutPolicy,
		ReadContext:   readLockoutPolicy,
	}
}

func deleteLockoutPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(lockoutPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetLockoutPolicyToDefault(ctx, &management2.ResetLockoutPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset lockout policy: %v", err)
	}
	return nil
}

func updateLockoutPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(lockoutPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomLockoutPolicy(ctx, &management2.UpdateCustomLockoutPolicyRequest{
		MaxPasswordAttempts: uint32(d.Get(lockoutPolicyMaxPasswordAttempts).(int)),
	})
	if err != nil {
		return diag.Errorf("failed to update lockout policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createLockoutPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(lockoutPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomLockoutPolicy(ctx, &management2.AddCustomLockoutPolicyRequest{
		MaxPasswordAttempts: uint32(d.Get(lockoutPolicyMaxPasswordAttempts).(int)),
	})
	if err != nil {
		return diag.Errorf("failed to create lockout policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readLockoutPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(domainPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLockoutPolicy(ctx, &management2.GetLockoutPolicyRequest{})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get lockout policy: %v", err)
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		lockoutPolicyOrgIdVar:            policy.GetDetails().GetResourceOwner(),
		lockoutPolicyMaxPasswordAttempts: policy.GetMaxPasswordAttempts(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of lockout policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
