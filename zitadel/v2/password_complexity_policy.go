package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	passwordCompPolicyOrgIdVar     = "org_id"
	passwordCompPolicyMinLength    = "min_length"
	passwordCompPolicyHasUppercase = "has_uppercase"
	passwordCompPolicyHasLowercase = "has_lowercase"
	passwordCompPolicyHasNumber    = "has_number"
	passwordCompPolicyHasSymbol    = "has_symbol"
	passwordCompPolicyIsDefault    = "is_default"
)

func GetPasswordComplexityPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			passwordCompPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			passwordCompPolicyMinLength: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimal length for the password",
			},
			passwordCompPolicyHasUppercase: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain an upper case letter",
			},
			passwordCompPolicyHasLowercase: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a lower case letter",
			},
			passwordCompPolicyHasNumber: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a number",
			},
			passwordCompPolicyHasSymbol: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a symbol. E.g. \"$\"",
			},
			passwordCompPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
		},
		DeleteContext: deletePasswordComplexityPolicy,
		ReadContext:   readPasswordComplexityPolicy,
		CreateContext: createPasswordComplexityPolicy,
		UpdateContext: updatePasswordComplexityPolicy,
	}
}

func deletePasswordComplexityPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(passwordCompPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetPasswordComplexityPolicyToDefault(ctx, &management2.ResetPasswordComplexityPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset password complexity policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func updatePasswordComplexityPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(passwordCompPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomPasswordComplexityPolicy(ctx, &management2.UpdateCustomPasswordComplexityPolicyRequest{
		MinLength:    uint64(d.Get(passwordCompPolicyMinLength).(int)),
		HasUppercase: d.Get(passwordCompPolicyHasUppercase).(bool),
		HasLowercase: d.Get(passwordCompPolicyHasLowercase).(bool),
		HasNumber:    d.Get(passwordCompPolicyHasNumber).(bool),
		HasSymbol:    d.Get(passwordCompPolicyHasSymbol).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update password complexity policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createPasswordComplexityPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(passwordCompPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomPasswordComplexityPolicy(ctx, &management2.AddCustomPasswordComplexityPolicyRequest{
		MinLength:    uint64(d.Get(passwordCompPolicyMinLength).(int)),
		HasUppercase: d.Get(passwordCompPolicyHasUppercase).(bool),
		HasLowercase: d.Get(passwordCompPolicyHasLowercase).(bool),
		HasNumber:    d.Get(passwordCompPolicyHasNumber).(bool),
		HasSymbol:    d.Get(passwordCompPolicyHasSymbol).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create password complexity policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readPasswordComplexityPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resp, err := client.GetPasswordComplexityPolicy(ctx, &management2.GetPasswordComplexityPolicyRequest{})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get password complexity policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		passwordCompPolicyOrgIdVar:     policy.GetDetails().GetResourceOwner(),
		passwordCompPolicyIsDefault:    policy.GetIsDefault(),
		passwordCompPolicyMinLength:    policy.GetMinLength(),
		passwordCompPolicyHasUppercase: policy.GetHasUppercase(),
		passwordCompPolicyHasLowercase: policy.GetHasLowercase(),
		passwordCompPolicyHasNumber:    policy.GetHasNumber(),
		passwordCompPolicyHasSymbol:    policy.GetHasSymbol(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of password complexity policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
