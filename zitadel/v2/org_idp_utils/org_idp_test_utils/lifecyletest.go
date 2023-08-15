package org_idp_test_utils

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func RunOrgLifecyleTest(
	t *testing.T,
	frame *test_utils.OrgTestFrame,
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
		CheckProviderName(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		CheckDestroy(*frame),
		test_utils.ConcatImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, secretAttribute),
		),
	)
}
