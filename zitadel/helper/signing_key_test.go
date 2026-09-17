package helper

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSigningKeyRotationDiff(t *testing.T) {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"signing_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_signing_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DurationDiffSuppress,
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			return SigningKeyRotationDiff(d, "expiration_signing_key", "signing_key")
		},
	}
	tests := []struct {
		name          string
		oldExpiration string
		newExpiration string
		wantUnknown   bool
	}{
		{
			name:          "setting an expiration marks the key unknown",
			oldExpiration: "",
			newExpiration: "0s",
			wantUnknown:   true,
		},
		{
			name:          "changing the expiration marks the key unknown",
			oldExpiration: "1h",
			newExpiration: "0s",
			wantUnknown:   true,
		},
		{
			name:          "an unchanged expiration keeps the key",
			oldExpiration: "0s",
			newExpiration: "0s",
			wantUnknown:   false,
		},
		{
			name:          "an equivalent expiration keeps the key",
			oldExpiration: "0s",
			newExpiration: "0",
			wantUnknown:   false,
		},
		{
			name:          "removing the expiration keeps the key",
			oldExpiration: "0s",
			newExpiration: "",
			wantUnknown:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID: "id",
				Attributes: map[string]string{
					"id":          "id",
					"signing_key": "key",
				},
			}
			if tt.oldExpiration != "" {
				state.Attributes["expiration_signing_key"] = tt.oldExpiration
			}
			config := map[string]interface{}{}
			if tt.newExpiration != "" {
				config["expiration_signing_key"] = tt.newExpiration
			}
			diff, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("Diff() error = %v", err)
			}
			gotUnknown := diff != nil && diff.Attributes["signing_key"] != nil && diff.Attributes["signing_key"].NewComputed
			if gotUnknown != tt.wantUnknown {
				t.Errorf("signing_key unknown = %v, want %v", gotUnknown, tt.wantUnknown)
			}
		})
	}
}
