package email_provider_smtp_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_smtp"
)

func TestCustomizeDiffForcesReplacementWhenUserConfigurationChanges(t *testing.T) {
	tests := []struct {
		name        string
		oldUser     string
		newUser     string
		requiresNew bool
	}{
		{name: "empty to non-empty", oldUser: "", newUser: "user", requiresNew: true},
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

func TestCustomizeDiffForcesReplacementWhenPasswordConfigurationChanges(t *testing.T) {
	tests := []struct {
		name        string
		oldHash     string
		newPassword string
		requiresNew bool
	}{
		{name: "empty to non-empty", newPassword: "password", requiresNew: true},
		{name: "non-empty to non-empty", oldHash: "existing-hash", newPassword: "new-password"},
		{name: "empty to empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := smtpDiff(
				t,
				map[string]string{"password_hash": tt.oldHash},
				map[string]interface{}{"password": tt.newPassword},
			)
			got := diff != nil && diff.Attributes["password"] != nil && diff.Attributes["password"].RequiresNew
			if got != tt.requiresNew {
				t.Fatalf("password RequiresNew = %v, want %v", got, tt.requiresNew)
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

	diff, err := email_provider_smtp.GetResource().Diff(
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
