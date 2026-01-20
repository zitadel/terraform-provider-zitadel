package zitadel

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
				"",        // systemAPIKeyFile
				"",        // systemAPIKey
				"",        // systemAPIUser
				"",        // systemAPIAudience
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
				"",            // systemAPIKeyFile
				"",            // systemAPIKey
				"",            // systemAPIUser
				"",            // systemAPIAudience
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

func expectedResourceTemplate(name string, includeImport bool) string {
	content := `---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

## Example Usage

{{ tffile "examples/provider/resources/` + name + `.tf" }}

{{ .SchemaMarkdown | trimspace }}
`
	if includeImport {
		content += `
## Import

{{ codefile "bash" "examples/provider/resources/` + name + `-import.sh" }}
`
	}
	return content
}

func expectedDataSourceTemplate(name string) string {
	return `---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

## Example Usage

{{ tffile "examples/provider/data-sources/` + name + `.tf" }}

{{ .SchemaMarkdown | trimspace }}
`
}

func TestDocumentationTemplates(t *testing.T) {
	sdkProvider := Provider()
	frameworkProvider := NewProviderPV6()

	t.Run("resources", func(t *testing.T) {
		type resourceInfo struct {
			hasImport bool
		}

		ctx := context.Background()
		resources := make(map[string]resourceInfo)

		for name, res := range sdkProvider.ResourcesMap {
			name := strings.TrimPrefix(name, "zitadel_")
			resources[name] = resourceInfo{hasImport: res.Importer != nil}
		}

		for _, factory := range frameworkProvider.Resources(ctx) {
			res := factory()
			var resp resource.MetadataResponse
			res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "zitadel"}, &resp)
			name := strings.TrimPrefix(resp.TypeName, "zitadel_")
			info := resources[name]
			if _, ok := res.(resource.ResourceWithImportState); ok {
				info.hasImport = true
			}
			resources[name] = info
		}

		names := make([]string, 0, len(resources))
		for name := range resources {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			name := name
			info := resources[name]
			t.Run(name, func(t *testing.T) {
				templatePath := filepath.Join("..", "templates", "resources", name+".md.tmpl")
				content, err := os.ReadFile(templatePath)
				if err != nil {
					t.Fatalf("failed to read template %s: %v", templatePath, err)
				}

				expected := expectedResourceTemplate(name, info.hasImport)
				if string(content) != expected {
					t.Fatalf("template %s does not match expected format", templatePath)
				}

				examplePath := filepath.Join("..", "examples", "provider", "resources", name+".tf")
				if _, err := os.Stat(examplePath); err != nil {
					t.Fatalf("example file %s missing: %v", examplePath, err)
				}

				if info.hasImport {
					importPath := filepath.Join("..", "examples", "provider", "resources", name+"-import.sh")
					if _, err := os.Stat(importPath); err != nil {
						t.Fatalf("import example %s missing: %v", importPath, err)
					}
				}
			})
		}
	})

	t.Run("data_sources", func(t *testing.T) {
		ctx := context.Background()
		dataSources := make(map[string]struct{})
		for name := range sdkProvider.DataSourcesMap {
			dataSources[strings.TrimPrefix(name, "zitadel_")] = struct{}{}
		}
		for _, factory := range frameworkProvider.DataSources(ctx) {
			ds := factory()
			var resp datasource.MetadataResponse
			ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "zitadel"}, &resp)
			dataSources[strings.TrimPrefix(resp.TypeName, "zitadel_")] = struct{}{}
		}

		names := make([]string, 0, len(dataSources))
		for name := range dataSources {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			name := name
			t.Run(name, func(t *testing.T) {
				templatePath := filepath.Join("..", "templates", "data-sources", name+".md.tmpl")
				content, err := os.ReadFile(templatePath)
				if err != nil {
					t.Fatalf("failed to read template %s: %v", templatePath, err)
				}

				expected := expectedDataSourceTemplate(name)
				if string(content) != expected {
					t.Fatalf("template %s does not match expected format", templatePath)
				}

				examplePath := filepath.Join("..", "examples", "provider", "data-sources", name+".tf")
				if _, err := os.Stat(examplePath); err != nil {
					t.Fatalf("example file %s missing: %v", examplePath, err)
				}
			})
		}
	})
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
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}
	gocache := filepath.Join(repoRoot, ".gocache")
	gomodcache := filepath.Join(repoRoot, ".go")
	if err := os.MkdirAll(gocache, 0o755); err != nil {
		t.Fatalf("failed to create gocache dir: %v", err)
	}
	if err := os.MkdirAll(gomodcache, 0o755); err != nil {
		t.Fatalf("failed to create gomodcache dir: %v", err)
	}

	generate := exec.Command("go", "run",
		"github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.14.1",
		"generate")
	generate.Dir = ".."
	generate.Env = append(os.Environ(),
		"GOCACHE="+gocache,
		"GOMODCACHE="+gomodcache,
	)
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
