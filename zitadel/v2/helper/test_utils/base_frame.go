package test_utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	insecure = true
	port     = "8080"
)

type BaseTestFrame struct {
	context.Context
	ClientInfo                         *helper.ClientInfo
	ProviderSnippet, UniqueResourcesID string
	TerraformName                      string
	v6ProviderFactories                map[string]func() (tfprotov6.ProviderServer, error)
	v5ProviderFactories                map[string]func() (tfprotov5.ProviderServer, error)
}

func NewBaseTestFrame(ctx context.Context, resourceType, domain string, jwtProfileJson []byte) (*BaseTestFrame, error) {
	zitadelProvider := zitadel.Provider()
	diag := zitadelProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":           domain,
		"insecure":         insecure,
		"port":             port,
		"jwt_profile_json": string(jwtProfileJson),
	}))
	if diag.HasError() {
		return nil, fmt.Errorf("unknown error configuring the test provider: %v", diag)
	}
	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
  domain   			= "%s"
  insecure 			= "%t"
  port     			= "%s"
  jwt_profile_json  = <<KEY
%s
KEY
}
`, domain, insecure, port, string(jwtProfileJson))
	clientInfo := zitadelProvider.Meta().(*helper.ClientInfo)
	uniqueID := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	terraformName := fmt.Sprintf("%s.%s", resourceType, uniqueID)
	frame := &BaseTestFrame{
		Context:           ctx,
		ProviderSnippet:   providerSnippet,
		ClientInfo:        clientInfo,
		UniqueResourcesID: uniqueID,
		TerraformName:     terraformName,
	}
	_, v5 := zitadelProvider.ResourcesMap[resourceType]
	if v5 {
		frame.v5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){"zitadel": func() (tfprotov5.ProviderServer, error) {
			return zitadelProvider.GRPCProvider(), nil
		}}
	} else {
		frame.v6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){"zitadel": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6(zitadel.NewProviderPV6())(), nil
		}}
	}
	return frame, nil
}
