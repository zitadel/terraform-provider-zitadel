package instance_custom_domain_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	instancev2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceCustomDomain(t *testing.T) {
	frame := test_utils.NewSystemTestFrame(t, "zitadel_instance_custom_domain")

	instanceClient, err := helper.GetInstanceClient(context.Background(), frame.ClientInfo)
	if err != nil {
		t.Fatalf("failed to get instance client: %v", err)
	}

	instanceResp, err := instanceClient.GetInstance(context.Background(), &instancev2.GetInstanceRequest{})
	if err != nil {
		t.Fatalf("failed to get instance: %v", err)
	}
	instanceID := instanceResp.GetInstance().GetId()
	testID := strings.ToLower(frame.UniqueResourcesID)

	resourceConfig := func(domain string, _ string) string {
		return fmt.Sprintf(`
resource "zitadel_instance_custom_domain" "default" {
    instance_id = "%s"
    domain      = "%s"
}
`, instanceID, domain)
	}

	test_utils.RunLifecyleTest(
		t,
		*frame,
		nil,
		resourceConfig,
		fmt.Sprintf("login-%s.example.com", testID),
		fmt.Sprintf("login-%s-2.example.com", testID),
		"",
		"",
		"",
		false,
		func(expect string) resource.TestCheckFunc {
			return test_utils.CheckNothing
		},
		regexp.MustCompile(`^.+$`),
		test_utils.CheckNothing,
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(*frame),
		),
	)
}
