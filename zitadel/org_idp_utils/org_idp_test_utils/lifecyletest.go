package org_idp_test_utils

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func RunOrgLifecyleTest(t *testing.T, resourceName, secretAttribute string) {
	frame := test_utils.NewOrgTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	nameProperty := test_utils.AttributeValue(t, idp_utils.NameVar, exampleAttributes).AsString()
	// Using a unique name makes the test idempotent on failures
	resourceExample = strings.Replace(resourceExample, nameProperty, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, idp_utils.IsCreationAllowedVar, exampleAttributes).True()
	exampleSecret := ""
	var importStateVerifyIgnore []string
	importParts := []resource.ImportStateIdFunc{
		test_utils.ImportResourceId(frame.BaseTestFrame),
		test_utils.ImportOrgId(frame),
	}
	if secretAttribute != "" {
		exampleSecret = test_utils.AttributeValue(t, secretAttribute, exampleAttributes).AsString()
		// The secret is write-only: it is never persisted to state, so it cannot
		// be sourced from state for the import ID, nor verified after import. Its
		// companion hash attribute is likewise absent immediately after import.
		importStateVerifyIgnore = []string{secretAttribute, secretAttribute + "_hash"}
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		true, false,
		secretAttribute, exampleSecret, "an_updated_secret",
		false,
		CheckCreationAllowed(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		CheckDestroy(*frame),
		test_utils.ChainImportStateIdFuncs(importParts...),
		importStateVerifyIgnore...,
	)
}
