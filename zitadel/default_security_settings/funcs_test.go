package default_security_settings

import (
	"testing"

	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"
)

func TestDynamicClientRegistrationSettings(t *testing.T) {
	yes, no := true, false
	tests := []struct {
		name                 string
		enabled              *bool
		allowUnauthenticated *bool
		current              *settingsv2.DynamicClientRegistrationSettings
		want                 *settingsv2.DynamicClientRegistrationSettings
	}{{
		name:    "unset attributes keep the current values",
		current: &settingsv2.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true},
		want:    &settingsv2.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true},
	}, {
		name: "unset attributes without current settings are false",
		want: &settingsv2.DynamicClientRegistrationSettings{},
	}, {
		name:                 "configured values replace the current values",
		enabled:              &no,
		allowUnauthenticated: &no,
		current:              &settingsv2.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true},
		want:                 &settingsv2.DynamicClientRegistrationSettings{},
	}, {
		name:    "a configured enabled keeps the current allow_unauthenticated",
		enabled: &yes,
		current: &settingsv2.DynamicClientRegistrationSettings{AllowUnauthenticated: true},
		want:    &settingsv2.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true},
	}, {
		name:                 "a configured allow_unauthenticated keeps the current enabled",
		allowUnauthenticated: &no,
		current:              &settingsv2.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true},
		want:                 &settingsv2.DynamicClientRegistrationSettings{Enabled: true},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dynamicClientRegistrationSettings(tt.enabled, tt.allowUnauthenticated, tt.current)
			if got.GetEnabled() != tt.want.GetEnabled() || got.GetAllowUnauthenticated() != tt.want.GetAllowUnauthenticated() {
				t.Errorf("got enabled=%t allow_unauthenticated=%t, want enabled=%t allow_unauthenticated=%t",
					got.GetEnabled(), got.GetAllowUnauthenticated(), tt.want.GetEnabled(), tt.want.GetAllowUnauthenticated())
			}
		})
	}
}
