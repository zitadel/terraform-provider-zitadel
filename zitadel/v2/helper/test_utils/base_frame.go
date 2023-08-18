package test_utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zclconf/go-cty/cty"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	insecure   = true
	port       = "8080"
	ResourceID = "123456789012345678"
)

type BaseTestFrame struct {
	context.Context
	ClientInfo                         *helper.ClientInfo
	ProviderSnippet, UniqueResourcesID string
	ResourceType                       string
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
	terraformName := fmt.Sprintf("%s.default", resourceType)
	frame := &BaseTestFrame{
		Context:           ctx,
		ProviderSnippet:   providerSnippet,
		ClientInfo:        clientInfo,
		UniqueResourcesID: uniqueID,
		TerraformName:     terraformName,
		ResourceType:      resourceType,
	}
	_, v5Resource := zitadelProvider.ResourcesMap[resourceType]
	_, v5Datasource := zitadelProvider.DataSourcesMap[resourceType]
	if v5Resource || v5Datasource {
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

func (b *BaseTestFrame) State(state *terraform.State) *terraform.InstanceState {
	resources := state.RootModule().Resources
	resource := resources[b.TerraformName]
	if resource != nil {
		return resource.Primary
	}
	return resources["data."+b.TerraformName].Primary
}

type examplesFolder string

const (
	Datasources examplesFolder = "data-sources"
	Resources   examplesFolder = "resources"
)

func (b *BaseTestFrame) ReadExample(t *testing.T, folder examplesFolder, exampleType string) (string, hcl.Attributes) {
	fileName := strings.Replace(exampleType, "zitadel_", "", 1) + ".tf"
	filePath := path.Join("..", "..", "..", "examples", "provider", string(folder), fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("error reading example file: %v", err)
	}
	hclFile, diags := hclparse.NewParser().ParseHCL(content, filePath)
	if diags.HasErrors() {
		t.Fatalf("error parsing example file: %s", diags.Error())
	}
	blocks := hclFile.BlocksAtPos(hcl.Pos{
		Line:   1,
		Column: 1,
		Byte:   1,
	})
	if len(blocks) != 1 {
		t.Fatalf("error parsing example file: %s", "unexpected number of blocks")
	}
	attr, diag := blocks[0].Body.JustAttributes()
	if diag.HasErrors() {
		t.Fatalf("error parsing example file: %s", diag.Error())
	}
	return string(content), attr
}

func AttributeValue(t *testing.T, key string, attributes hcl.Attributes) cty.Value {
	val, diag := attributes[key].Expr.Value(&hcl.EvalContext{})
	if diag.HasErrors() {
		t.Fatalf("error parsing example file: %s", diag.Error())
	}
	return val
}
