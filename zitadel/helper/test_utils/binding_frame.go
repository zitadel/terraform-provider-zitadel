package test_utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel"
)

// NewBindingTestFrame uses the normal acceptance instance in CI. Optional
// ZITADEL_TEST_* settings let the same acceptance cases target a disposable
// local instance without replacing the repository's embedded signing keys.
func NewBindingTestFrame(t *testing.T, resourceType string) *BaseTestFrame {
	t.Helper()
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("live acceptance test; set TF_ACC=1")
	}
	domain := os.Getenv("ZITADEL_TEST_DOMAIN")
	if domain == "" {
		return &NewInstanceTestFrame(t, resourceType).BaseTestFrame
	}
	token := os.Getenv("ZITADEL_TEST_ACCESS_TOKEN")
	if token == "" {
		t.Fatal("ZITADEL_TEST_ACCESS_TOKEN is required with ZITADEL_TEST_DOMAIN")
	}
	port := os.Getenv("ZITADEL_TEST_PORT")
	if port == "" {
		port = "443"
	}
	insecure := false
	if raw := os.Getenv("ZITADEL_TEST_INSECURE"); raw != "" {
		var err error
		insecure, err = strconv.ParseBool(raw)
		if err != nil {
			t.Fatal(err)
		}
	}
	headers := map[string]string{}
	if raw := os.Getenv("ZITADEL_TEST_TRANSPORT_HEADERS"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &headers); err != nil {
			t.Fatal(err)
		}
	}
	p := zitadel.Provider()
	ctx := context.Background()
	headerValues := map[string]interface{}{}
	for key, value := range headers {
		headerValues[key] = value
	}
	if d := p.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain": domain, "port": port, "insecure": insecure,
		"access_token": token, "transport_headers": headerValues,
	})); d.HasError() {
		t.Fatal("could not configure live acceptance provider")
	}
	// Keep the PAT out of rendered HCL and assertion/error output.
	t.Setenv("TF_VAR_test_access_token", token)
	headerJSON, err := json.Marshal(headers)
	if err != nil {
		t.Fatal(err)
	}
	snippet := fmt.Sprintf(`
variable "test_access_token" {
  type = string
  sensitive = true
}
provider "zitadel" {
  domain = %q
  port = %q
  insecure = %t
  access_token = var.test_access_token
  transport_headers = jsondecode(%q)
}
`, domain, port, insecure, string(headerJSON))
	return buildBaseTestFrame(ctx, resourceType, domain, snippet, p)
}
