package org_test_dep

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/org"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, resourceName string) (string, string, *test_utils.OrgTestFrame) {
	otherFrame := frame.AnotherOrg(t, fmt.Sprintf("%s_%s", resourceName, frame.UniqueResourcesID))
	cfg, id := test_utils.CreateDefaultDependency(t, "zitadel_org", org.OrgIDVar, func() (string, error) {
		return otherFrame.OrgID, nil
	})
	return strings.Replace(cfg, "default", resourceName, 1), id, otherFrame
}
