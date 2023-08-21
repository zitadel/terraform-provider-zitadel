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
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/default_label_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
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
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
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
