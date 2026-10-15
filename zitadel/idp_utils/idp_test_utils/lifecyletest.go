package idp_test_utils

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

// RunInstanceIDPLifecyleTest runs the standard IDP lifecycle test. secretOverride
// optionally supplies the initial and updated secret values (exactly two entries:
// create, update) instead of the example value and the generic "an_updated_secret".
// Use it for IDPs whose secret is format-validated by the server (e.g. Apple keys).
func RunInstanceIDPLifecyleTest(t *testing.T, resourceName, secretAttribute string, secretOverride ...string) {
	frame := test_utils.NewInstanceTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	nameProperty := test_utils.AttributeValue(t, idp_utils.NameVar, exampleAttributes).AsString()
	// Using a unique name makes the test idempotent on failures
	resourceExample = strings.Replace(resourceExample, nameProperty, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, idp_utils.IsCreationAllowedVar, exampleAttributes).True()
	exampleSecret := ""
	var importStateVerifyIgnore []string
	importParts := []resource.ImportStateIdFunc{
		test_utils.ImportResourceId(frame.BaseTestFrame),
	}
	if secretAttribute != "" {
		exampleSecret = test_utils.AttributeValue(t, secretAttribute, exampleAttributes).AsString()
		// The secret is write-only: it is never persisted to state, so it cannot
		// be sourced from state for the import ID, nor verified after import. Its
		// companion hash attribute is likewise absent immediately after import.
		importStateVerifyIgnore = []string{secretAttribute, secretAttribute + "_hash"}
	}
	createSecret, updatedSecret := exampleSecret, "an_updated_secret"
	if len(secretOverride) == 2 {
		createSecret, updatedSecret = secretOverride[0], secretOverride[1]
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		true, false,
		secretAttribute, createSecret, updatedSecret,
		false,
		CheckCreationAllowed(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		CheckDestroy(*frame),
		test_utils.ChainImportStateIdFuncs(importParts...),
		importStateVerifyIgnore...,
	)
}

// RunInstanceIDPUnchangedSecretTest checks that updating an IDP without
// changing its secret does not send the secret to ZITADEL, while changing the
// secret still does. secretField is the secret's key in the eventType payload.
// secretOverride works as in RunInstanceIDPLifecyleTest.
func RunInstanceIDPUnchangedSecretTest(t *testing.T, resourceName, eventType, secretAttribute, secretField string, secretOverride ...string) {
	frame := test_utils.NewInstanceTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	nameProperty := test_utils.AttributeValue(t, idp_utils.NameVar, exampleAttributes).AsString()
	exampleSecret := test_utils.AttributeValue(t, secretAttribute, exampleAttributes).AsString()
	config := func(name, secret string) string {
		cfg := strings.Replace(resourceExample, nameProperty, name, 1)
		cfg = strings.Replace(cfg, exampleSecret, secret, 1)
		return frame.ProviderSnippet + "\n" + cfg
	}
	createSecret, updatedSecret := exampleSecret, "an_updated_secret"
	if len(secretOverride) == 2 {
		createSecret, updatedSecret = secretOverride[0], secretOverride[1]
	}
	checkSecretSent := CheckLatestChangedEventSecret(*frame, eventType, secretField)
	renamed := frame.UniqueResourcesID + "_renamed"
	resource.ParallelTest(t, resource.TestCase{
		CheckDestroy: test_utils.CheckAMinute(CheckDestroy(*frame)),
		Steps: []resource.TestStep{
			{ // Create the IDP
				Config: config(frame.UniqueResourcesID, createSecret),
			}, { // Renaming must not re-send the unchanged secret
				Config: config(renamed, createSecret),
				Check:  test_utils.CheckAMinute(checkSecretSent(false)),
			}, { // Changing the secret must send it
				Config: config(renamed, updatedSecret),
				Check:  test_utils.CheckAMinute(checkSecretSent(true)),
			},
		},
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
	})
}
