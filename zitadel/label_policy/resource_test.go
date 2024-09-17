package label_policy_test

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
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/label_policy"
)

func TestAccLabelPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_label_policy")
	testSVGFile := writeFile(t, strings.NewReader(testSVG))
	defer testSVGFile.Close()
	defer os.Remove(testSVGFile.Name())
	testFontFile := writeFile(t, base64.NewDecoder(base64.StdEncoding, strings.NewReader(testFontBase64)))
	defer testFontFile.Close()
	defer os.Remove(testFontFile.Name())
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	resourceExample = regexp.MustCompile("/path/to/[a-zA-Z_]+\\.jpg").ReplaceAllString(resourceExample, testSVGFile.Name())
	resourceExample = regexp.MustCompile("/path/to/[a-zA-Z_]+\\.tff").ReplaceAllString(resourceExample, testFontFile.Name())
	exampleProperty := test_utils.AttributeValue(t, "primary_color", exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "#5469d3",
		"", "", "",
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(*frame)(exampleProperty),
		test_utils.ImportOrgId(frame),
		label_policy.SetActiveVar,
		label_policy.LogoHashVar,
		label_policy.LogoPathVar,
		label_policy.LogoDarkHashVar,
		label_policy.LogoDarkPathVar,
		label_policy.IconHashVar,
		label_policy.IconPathVar,
		label_policy.IconDarkHashVar,
		label_policy.IconDarkPathVar,
		label_policy.FontHashVar,
		label_policy.FontPathVar,
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLabelPolicy(frame, &management.GetLabelPolicyRequest{})
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
	testSVG = `
<svg height="100" width="100">
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
