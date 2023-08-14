package idp_test_utils

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func RunInstanceIDPLifecyleTest(
	t *testing.T,
	frame test_utils.InstanceTestFrame,
	resourceFunc func(string, string) string,
	secretAttribute string,
) {
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		resourceFunc,
		"an initial provider name", "an updated provider name",
		secretAttribute, "an_initial_secret", "an_updated_secret",
		false,
		CheckProviderName(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		CheckDestroy(frame),
		test_utils.ImportStateId(frame.BaseTestFrame, secretAttribute),
	)
}
