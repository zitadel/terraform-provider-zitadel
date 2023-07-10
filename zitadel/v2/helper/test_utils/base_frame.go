package test_utils

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	domain   = "localhost"
	insecure = true
	port     = "8080"
)

type BaseTestFrame struct {
	context.Context
	upgradedV5ProviderFactory, v6ProviderFactory func() (tfprotov6.ProviderServer, error)
	ClientInfo                                   *helper.ClientInfo
	ProviderSnippet, UniqueResourcesID           string
	TerraformName                                string
}

func NewBaseTestFrame(resourceType string) (*BaseTestFrame, error) {
	ctx := context.Background()
	tokenPath := os.Getenv("TF_ACC_ZITADEL_TOKEN")
	zitadelProvider := zitadel.Provider()
	diag := zitadelProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":   domain,
		"insecure": insecure,
		"port":     port,
		"token":    tokenPath,
	}))
	if diag.HasError() {
		return nil, fmt.Errorf("unknown error configuring the test provider: %v", diag)
	}

	providerConfigSnippet := fmt.Sprintf(`
  domain   = "%s"
  insecure = "%t"
  port     = "%s"
  token    = "%s"
`, domain, insecure, port, tokenPath)

	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
	%s
}

provider "upgraded-v5" {
	%s
}
`, providerConfigSnippet, providerConfigSnippet)
	clientInfo := zitadelProvider.Meta().(*helper.ClientInfo)
	uniqueID := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	terraformName := fmt.Sprintf("%s.%s", resourceType, uniqueID)

	upgradedV5Provider, err := tf5to6server.UpgradeServer(ctx, zitadel.Provider().GRPCProvider)
	if err != nil {
		return nil, err
	}

	return &BaseTestFrame{
		Context: ctx,
		upgradedV5ProviderFactory: func() (tfprotov6.ProviderServer, error) {
			return upgradedV5Provider, nil
		},
		v6ProviderFactory: func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6(zitadel.NewProviderPV6())(), nil
		},
		ProviderSnippet:   providerSnippet,
		ClientInfo:        clientInfo,
		UniqueResourcesID: uniqueID,
		TerraformName:     terraformName,
	}, nil
}
