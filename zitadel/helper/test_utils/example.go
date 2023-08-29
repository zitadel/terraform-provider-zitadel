package test_utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type examplesFolder string

const (
	Datasources examplesFolder = "data-sources"
	Resources   examplesFolder = "resources"
)

func ReadExample(t *testing.T, folder examplesFolder, exampleType string) (string, hcl.Attributes) {
	fileName := strings.Replace(exampleType, "zitadel_", "", 1) + ".tf"
	filePath := path.Join("..", "..", "examples", "provider", string(folder), fileName)
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

func ReplaceAll[T comparable](resourceExample string, exampleProperty T, exampleSecret string) func(T, string) string {
	return func(configProperty T, secretProperty string) string {
		cfg := strings.ReplaceAll(resourceExample, fmt.Sprintf("%v", exampleProperty), fmt.Sprintf("%v", configProperty))
		if secretProperty != "" {
			cfg = strings.Replace(cfg, exampleSecret, secretProperty, 1)
		}
		return cfg
	}
}
