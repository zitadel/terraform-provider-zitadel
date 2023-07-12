package default_oidc_settings_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultOIDCSettings(t *testing.T) {
	resourceName := "zitadel_default_oidc_settings"
	initialAccessTokenLifetime := "123h0m0s"
	updatedAccessTokenLifetime := "456h0m0s"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(accessTokenLifetime, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	access_token_lifetime = "%s"
  	id_token_lifetime = "777h0m0s"
  	refresh_token_idle_expiration = "888h0m0s"
  	refresh_token_expiration = "999h0m0s"
}`, resourceName, frame.UniqueResourcesID, accessTokenLifetime)
		},
		initialAccessTokenLifetime, updatedAccessTokenLifetime,
		"", "",
		checkAccessTokenLifetime(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		func(state *terraform.State) error { return nil },
		nil, nil, "", "",
	)
}

func checkAccessTokenLifetime(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expectAccessTokenLifetime string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetOIDCSettings(frame, &admin.GetOIDCSettingsRequest{})
			if err != nil {
				return fmt.Errorf("getting oidc settings failed: %w", err)
			}
			actual := resp.GetSettings().GetAccessTokenLifetime().AsDuration().String()
			if actual != expectAccessTokenLifetime {
				return fmt.Errorf("expected access token lifetime %s, but got %s", expectAccessTokenLifetime, actual)
			}
			return nil
		}
	}
}
