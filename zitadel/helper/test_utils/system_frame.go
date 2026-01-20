package test_utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel"
)

func NewSystemTestFrame(t *testing.T, resourceType string) *BaseTestFrame {
	ctx := context.Background()
	cfg := acceptance.GetConfig().SystemAPI

	zitadelProvider := zitadel.Provider()
	diag := zitadelProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":   cfg.Domain,
		"insecure": insecure,
		"port":     port,
		"system_api": []interface{}{
			map[string]interface{}{
				"user": cfg.User,
				"key":  string(cfg.KeyPEM),
			},
		},
	}))
	if diag.HasError() {
		t.Fatalf("setting up system test context failed: %v", diag)
	}

	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
  domain            = "%s"
  insecure          = %t
  port              = "%s"
  system_api {
    user = "%s"
    key  = <<KEY
%s
KEY
  }
}
`, cfg.Domain, insecure, port, cfg.User, string(cfg.KeyPEM))

	return buildBaseTestFrame(ctx, resourceType, cfg.Domain, providerSnippet, zitadelProvider)
}
