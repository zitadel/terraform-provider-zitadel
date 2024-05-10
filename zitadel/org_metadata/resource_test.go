package org_metadata_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/org_metadata"
)

func TestAccOrgMetadata(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_metadata")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	keyProperty := test_utils.AttributeValue(t, org_metadata.KeyVar, exampleAttributes).AsString()
	exampleProperty := test_utils.AttributeValue(t, org_metadata.ValueVar, exampleAttributes).AsString()
	updatedProperty := "YW5vdGhlciB2YWx1ZQ=="
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame),
		regexp.MustCompile(fmt.Sprintf(`^%s$`, keyProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, org_metadata.KeyVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetOrgMetadata(frame, &management.GetOrgMetadataRequest{
				Key: frame.State(state).ID,
			})
			if err != nil {
				return err
			}
			actual := helper.Base64Encode(resp.GetMetadata().GetValue())
			if expect != actual {
				return fmt.Errorf("expected role %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
