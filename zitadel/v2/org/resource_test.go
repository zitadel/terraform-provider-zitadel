package org_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org"
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
		"", "",
		false,
		checkRemoteProperty(frame, idFromState(frame)),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, idFromState(frame)), updatedProperty),
		nil, nil, "", "",
	)
}

func idFromState(frame *test_utils.OrgTestFrame) func(*terraform.State) string {
	return func(state *terraform.State) string {
		return frame.State(state).ID
	}
}
