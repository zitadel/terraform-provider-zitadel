package org_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org"
)

func TestAccOrg(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org.NameVar, exampleAttributes).AsString()
	initialProperty := "initialorgname_" + frame.UniqueResourcesID
	updatedProperty := "updatedorgname_" + frame.UniqueResourcesID
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, idFromState(frame)),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, idFromState(frame)), updatedProperty),
		test_utils.ImportResourceId(frame.BaseTestFrame),
	)
}

func idFromState(frame *test_utils.OrgTestFrame) func(*terraform.State) string {
	return func(state *terraform.State) string {
		return frame.State(state).ID
	}
}
