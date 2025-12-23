package zitadel

import (
	"context"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func TestClientInfo_Schemes(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		insecure   bool
		wantIssuer string
	}{
		{
			name:       "standard_domain",
			domain:     "instance.zitadel.cloud",
			insecure:   false,
			wantIssuer: "https://instance.zitadel.cloud",
		},
		{
			name:       "domain_with_https",
			domain:     "https://instance.zitadel.cloud",
			insecure:   false,
			wantIssuer: "https://instance.zitadel.cloud",
		},
		{
			name:       "domain_with_http_insecure",
			domain:     "http://localhost",
			insecure:   true,
			wantIssuer: "http://localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyJSON := `{"type":"service_account"}`
			info, err := helper.GetClientInfo(
				context.Background(),
				tt.insecure,
				tt.domain,
				"",        // accessToken
				"",        // token
				"",        // jwtFile
				"",        // jwtProfileFile
				dummyJSON, // jwtProfileJSON
				"",        // port
			)
			if err != nil {
				t.Fatalf("GetClientInfo() error = %v", err)
			}
			if info.Issuer != tt.wantIssuer {
				t.Errorf("GetClientInfo() Issuer = %v, want %v", info.Issuer, tt.wantIssuer)
			}
		})
	}
}

func TestClientInfo_Files(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		jwtFile    string
		jwtProfile string
		wantError  bool
	}{
		{
			name:      "missing_token_file",
			token:     "non_existent_token.json",
			wantError: true,
		},
		{
			name:      "missing_jwt_file",
			jwtFile:   "non_existent_jwt.json",
			wantError: true,
		},
		{
			name:       "missing_jwt_profile_file",
			jwtProfile: "non_existent_profile.json",
			wantError:  true,
		},
		{
			name:      "no_credentials",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := helper.GetClientInfo(
				context.Background(),
				false,
				"example.com",
				"",            // accessToken
				tt.token,      // token
				tt.jwtFile,    // jwtFile
				tt.jwtProfile, // jwtProfileFile
				"",            // jwtProfileJSON
				"",            // port
			)
			if (err != nil) != tt.wantError {
				t.Errorf("GetClientInfo() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
