package test_utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

const (
	insecure           = true
	port               = "8080"
	ExamplesResourceID = "123456789012345678"
)

type BaseTestFrame struct {
	context.Context
	ClientInfo                         *helper.ClientInfo
	ProviderSnippet, UniqueResourcesID string
	ResourceType                       string
	InstanceDomain                     string
	TerraformName                      string
	v6ProviderFactories                map[string]func() (tfprotov6.ProviderServer, error)
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
	terraformName := fmt.Sprintf("%s.default", resourceType)
	frame := &BaseTestFrame{
		Context:           ctx,
		ProviderSnippet:   providerSnippet,
		ClientInfo:        clientInfo,
		UniqueResourcesID: uniqueID,
		TerraformName:     terraformName,
		ResourceType:      resourceType,
		InstanceDomain:    domain,
	}
	frame.v6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"zitadel": func() (tfprotov6.ProviderServer, error) {
			muxServer, err := tf6muxserver.NewMuxServer(frame,
				providerserver.NewProtocol6(zitadel.NewProviderPV6()),
				func() tfprotov6.ProviderServer {
					upgraded, err := tf5to6server.UpgradeServer(frame, func() tfprotov5.ProviderServer {
						return zitadelProvider.GRPCProvider()
					})
					if err != nil {
						return nil
					}
					return upgraded
				},
			)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}
	return frame, nil
}

func (b *BaseTestFrame) State(state *terraform.State) *terraform.InstanceState {
	resources := state.RootModule().Resources
	resource := resources[b.TerraformName]
	if resource != nil {
		return resource.Primary
	}
	resource = resources["data."+b.TerraformName]
	if resource != nil {
		return resource.Primary
	}
	return &terraform.InstanceState{}
}
