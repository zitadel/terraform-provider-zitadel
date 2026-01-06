package zitadel

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

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
				false,     // insecureSkipVerifyTLS
				nil,       // transportHeaders
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
				false,         // insecureSkipVerifyTLS
				nil,           // transportHeaders
			)
			if (err != nil) != tt.wantError {
				t.Errorf("GetClientInfo() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestDocumentation verifies that the generated documentation in the docs/
// directory is up-to-date with the current provider schema.
//
// This test runs tfplugindocs to regenerate documentation and then checks
// each resource and datasource individually. If any documentation file has
// changed or is missing, that specific subtest fails.
//
// To fix failing tests, run:
//
//	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.14.1 generate
//
// Then commit the updated documentation files.
func TestDocumentation(t *testing.T) {
	generate := exec.Command("go", "run",
		"github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.14.1",
		"generate")
	generate.Dir = ".."
	if output, err := generate.CombinedOutput(); err != nil {
		t.Fatalf("tfplugindocs generate failed: %v\n%s", err, output)
	}

	status := exec.Command("git", "status", "--porcelain", "docs/")
	status.Dir = ".."
	output, err := status.Output()
	if err != nil {
		t.Fatalf("git status failed: %v", err)
	}

	changedFiles := make(map[string]bool)
	for _, line := range strings.Split(string(output), "\n") {
		if len(line) > 3 {
			changedFiles[strings.TrimSpace(line[2:])] = true
		}
	}

	sdkProvider := Provider()
	frameworkProvider := NewProviderPV6()

	t.Run("resources", func(t *testing.T) {
		for name := range sdkProvider.ResourcesMap {
			name := strings.TrimPrefix(name, "zitadel_")
			t.Run(name, func(t *testing.T) {
				docPath := filepath.Join("docs", "resources", name+".md")
				if changedFiles[docPath] {
					t.Errorf("documentation is out of date")
				}
			})
		}
		for _, factory := range frameworkProvider.Resources(context.Background()) {
			res := factory()
			var resp resource.MetadataResponse
			res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zitadel"}, &resp)
			name := strings.TrimPrefix(resp.TypeName, "zitadel_")
			t.Run(name, func(t *testing.T) {
				docPath := filepath.Join("docs", "resources", name+".md")
				if changedFiles[docPath] {
					t.Errorf("documentation is out of date")
				}
			})
		}
	})

	t.Run("data_sources", func(t *testing.T) {
		for name := range sdkProvider.DataSourcesMap {
			name := strings.TrimPrefix(name, "zitadel_")
			t.Run(name, func(t *testing.T) {
				docPath := filepath.Join("docs", "data-sources", name+".md")
				if changedFiles[docPath] {
					t.Errorf("documentation is out of date")
				}
			})
		}
	})
}
