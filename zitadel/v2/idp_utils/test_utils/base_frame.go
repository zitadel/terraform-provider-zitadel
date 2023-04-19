package test_utils

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	Domain   = "localhost"
	insecure = true
	port     = "8080"
)

type BaseTestFrame struct {
	context.Context
	ConfiguredProvider                 *schema.Provider
	ClientInfo                         *helper.ClientInfo
	ProviderSnippet, UniqueResourcesID string
	TerraformName                      string
}

func NewBaseTestFrame(resourceType string) (*BaseTestFrame, error) {
	ctx := context.Background()
	tokenPath := os.Getenv("TF_ACC_ZITADEL_TOKEN")
	zitadelProvider := zitadel.Provider()
	diag := zitadelProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":   Domain,
		"insecure": insecure,
		"port":     port,
		"token":    tokenPath,
	}))
	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
  domain   = "%s"
  insecure = "%t"
  port     = "%s"
  token    = "%s"
}
`, Domain, insecure, port, tokenPath)
	if diag.HasError() {
		return nil, fmt.Errorf("unknown error configuring the test provider: %v", diag)
	}
	clientInfo := zitadelProvider.Meta().(*helper.ClientInfo)
	uniqueID := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	terraformName := fmt.Sprintf("%s.%s", resourceType, uniqueID)

	return &BaseTestFrame{
		Context:            ctx,
		ConfiguredProvider: zitadelProvider,
		ProviderSnippet:    providerSnippet,
		ClientInfo:         clientInfo,
		UniqueResourcesID:  uniqueID,
		TerraformName:      terraformName,
	}, nil
}
