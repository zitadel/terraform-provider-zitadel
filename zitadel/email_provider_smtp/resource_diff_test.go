package email_provider_smtp_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_smtp"
)

func TestCustomizeDiffForcesReplacementWhenUserIsRemoved(t *testing.T) {
	tests := []struct {
		name        string
		oldUser     string
		newUser     string
		requiresNew bool
	}{
		{name: "empty to non-empty", oldUser: "", newUser: "user"},
		{name: "non-empty to empty", oldUser: "user", newUser: "", requiresNew: true},
		{name: "non-empty to non-empty", oldUser: "old-user", newUser: "new-user"},
		{name: "empty to empty", oldUser: "", newUser: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := smtpDiff(t, map[string]string{"user": tt.oldUser}, map[string]interface{}{"user": tt.newUser})
			got := diff != nil && diff.Attributes["user"] != nil && diff.Attributes["user"].RequiresNew
			if got != tt.requiresNew {
				t.Fatalf("user RequiresNew = %v, want %v", got, tt.requiresNew)
			}
		})
	}
}

func TestCustomizeDiffForcesReplacementWhenPasswordIsRemoved(t *testing.T) {
	tests := []struct {
		name        string
		oldHash     string
		newPassword string
		requiresNew bool
	}{
		{name: "non-empty to empty", oldHash: "existing-hash", requiresNew: true},
		// covers state from older provider versions and imports, where password_hash is not set yet
		{name: "empty to non-empty", newPassword: "password"},
		{name: "non-empty to non-empty", oldHash: "existing-hash", newPassword: "new-password"},
		{name: "empty to empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := map[string]interface{}{}
			if tt.newPassword != "" {
				config["password"] = tt.newPassword
			}
			diff := smtpDiff(t, map[string]string{"password_hash": tt.oldHash}, config)
			got := diff != nil && diff.RequiresNew()
			if got != tt.requiresNew {
				t.Fatalf("RequiresNew = %v, want %v", got, tt.requiresNew)
			}
		})
	}
}

func smtpDiff(t *testing.T, stateAttributes map[string]string, config map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()

	state := &terraform.InstanceState{
		ID:         "id",
		Attributes: stateAttributes,
	}
	for key, value := range map[string]string{
		"sender_address": "sender@example.com",
		"sender_name":    "sender",
		"host":           "localhost:25",
	} {
		state.Attributes[key] = value
		config[key] = value
	}

	resource := email_provider_smtp.GetResource()
	state.RawConfig = rawConfig(resource.CoreConfigSchema().ImpliedType(), config)

	diff, err := resource.Diff(
		context.Background(),
		state,
		terraform.NewResourceConfigRaw(config),
		nil,
	)
	if err != nil {
		t.Fatalf("Diff() error = %v", err)
	}
	return diff
}

// rawConfig builds the cty config that the SDK exposes through GetRawConfig,
// which write-only attributes are read from.
func rawConfig(ty cty.Type, config map[string]interface{}) cty.Value {
	attrs := map[string]cty.Value{}
	for name, attrType := range ty.AttributeTypes() {
		value, ok := config[name].(string)
		if !ok {
			attrs[name] = cty.NullVal(attrType)
			continue
		}
		attrs[name] = cty.StringVal(value)
	}
	return cty.ObjectVal(attrs)
}
