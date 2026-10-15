package helper

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestChangedWriteOnlyStringValue(t *testing.T) {
	tests := []struct {
		name    string
		oldHash string
		newHash string
		want    string
	}{
		{
			name:    "unchanged secret is omitted",
			oldHash: "hash",
			newHash: "hash",
			want:    "",
		},
		{
			name:    "rotated secret is sent",
			oldHash: "hash",
			newHash: "rotated_hash",
			want:    "secret",
		},
		{
			name:    "secret without a recorded hash is sent",
			oldHash: "",
			newHash: "hash",
			want:    "secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			r := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name":        {Type: schema.TypeString, Optional: true},
					"secret":      {Type: schema.TypeString, Optional: true, WriteOnly: true},
					"secret_hash": {Type: schema.TypeString, Computed: true},
				},
				UpdateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
					got = ChangedWriteOnlyStringValue(d, "secret", "secret_hash")
					return nil
				},
			}
			state := &terraform.InstanceState{
				ID:         "id",
				Attributes: map[string]string{"id": "id", "name": "old", "secret_hash": tt.oldHash},
			}
			diff := &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"name": {Old: "old", New: "new"},
				},
				RawConfig: cty.ObjectVal(map[string]cty.Value{
					"name":   cty.StringVal("new"),
					"secret": cty.StringVal("secret"),
				}),
			}
			if tt.oldHash != tt.newHash {
				diff.Attributes["secret_hash"] = &terraform.ResourceAttrDiff{Old: tt.oldHash, New: tt.newHash}
			}
			if _, diags := r.Apply(context.Background(), state, diff, nil); diags.HasError() {
				t.Fatalf("Apply() diagnostics = %v", diags)
			}
			if got != tt.want {
				t.Errorf("ChangedWriteOnlyStringValue() = %q, want %q", got, tt.want)
			}
		})
	}
}
