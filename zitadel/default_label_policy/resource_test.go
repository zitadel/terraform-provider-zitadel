package default_label_policy_test

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_label_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultLabelPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_label_policy")
	testSVGFile := writeFile(t, strings.NewReader(testSVG))
	defer testSVGFile.Close()
	defer os.Remove(testSVGFile.Name())
	testFontFile := writeFile(t, base64.NewDecoder(base64.StdEncoding, strings.NewReader(testFontBase64)))
	defer testFontFile.Close()
	defer os.Remove(testFontFile.Name())
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	resourceExample = regexp.MustCompile("/path/to/[a-zA-Z_]+\\.jpg").ReplaceAllString(resourceExample, testSVGFile.Name())
	resourceExample = regexp.MustCompile("/path/to/[a-zA-Z_]+\\.tff").ReplaceAllString(resourceExample, testFontFile.Name())
	exampleProperty := test_utils.AttributeValue(t, default_label_policy.PrimaryColorVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "#5469d3",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
		default_label_policy.SetActiveVar,
		default_label_policy.LogoHashVar,
		default_label_policy.LogoPathVar,
		default_label_policy.LogoDarkHashVar,
		default_label_policy.LogoDarkPathVar,
		default_label_policy.IconHashVar,
		default_label_policy.IconPathVar,
		default_label_policy.IconDarkHashVar,
		default_label_policy.IconDarkPathVar,
		default_label_policy.FontHashVar,
		default_label_policy.FontPathVar,
	)
}

func TestAccDefaultLabelPolicyActivation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_label_policy")

	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_default_label_policy" "default" {
  primary_color          = "#5469d4"
  hide_login_name_suffix = false
  warn_color             = "#cd3d56"
  background_color       = "#ffffff"
  font_color             = "#000000"
  primary_color_dark     = "#2073c4"
  background_color_dark  = "#111827"
  warn_color_dark        = "#ff3b5b"
  font_color_dark        = "#ffffff"
  disable_watermark      = false
  set_active             = true
}
`, frame.ProviderSnippet)

	updatedConfig := fmt.Sprintf(`
%s
resource "zitadel_default_label_policy" "default" {
  primary_color          = "#22c55e"
  hide_login_name_suffix = false
  warn_color             = "#cd3d56"
  background_color       = "#ffffff"
  font_color             = "#000000"
  primary_color_dark     = "#2073c4"
  background_color_dark  = "#111827"
  warn_color_dark        = "#ff3b5b"
  font_color_dark        = "#ffffff"
  disable_watermark      = false
  set_active             = true
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "primary_color", "#5469d4"),
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "primary_color", "#22c55e"),
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLabelPolicy(frame, &admin.GetLabelPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetPrimaryColor()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}

func writeFile(t *testing.T, content io.Reader) *os.File {
	file, err := os.CreateTemp("", "TestAccDefaultLabelPolicy")
	if err != nil {
		t.Fatalf("creating temp file failed: %v", err)
	}
	if _, err := io.Copy(file, content); err != nil {
		t.Fatalf("writing temp file failed: %v", err)
	}
	return file
}

const (
	testSVG = `<?xml version="1.0" encoding="UTF-8"?>
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" height="100" width="100">
<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />
</svg>
`
	testFontBase64 = `
AAEAAAAHAEAAAgAwY21hcAAJAHYAAAEAAAAALGdseWbxy2aYAAABNAAAAFxoZWFk8jXd+AAAAHwA
AAA2aGhlYQZhAMoAAAC0AAAAJGhtdHgEdABqAAAA+AAAAAhsb2NhAC4AFAAAASwAAAAGbWF4cAAF
AAsAAADYAAAAIAABAAAAAQAA9ZwpRF8PPPUAAgPoAAAAALSS9AAAAAAA3C+mXAAGAAACWAK8AAAA
AwACAAAAAAAAAAEAAAQA/nAAAAJYAAb//wJYAAEAAAAAAAAAAAAAAAAAAAACAAEAAAACAAsAAgAA
AAAAAAAAAAAAAAAAAAAAAAAAAAACWABkAhwABgAAAAEAAAADAAAADAAEACAAAAAEAAQAAQAAAEH/
/wAAAEH////AAAEAAAAAAAAAFAAuAAAAAgBkAAACWAK8AAMABwAAMxEhESUhESFkAfT+NAGk/lwC
vP1EKAJsAAIABgAAAh0CkAACAAoAABMzAwETMxMjJyMHrcRj/vjaYN1ZPu9CAQsBQP21ApD9cMjI
AA==
`
)
