package project_v2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	projectpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_v2"
)

// TestAccProjectV2 exercises the full create/update/delete/import lifecycle
// of the new zitadel_project_v2 resource (backed by the Zitadel v2 project
// API) against a live instance. Mirrors the v1 TestAccProject pattern.
func TestAccProjectV2(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_project_v2")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, project_v2.NameVar, exampleAttributes).AsString()
	updatedProperty := "updatedproperty"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
	)
}

// checkRemoteProperty verifies that the project's name in Zitadel matches the
// expected value by calling the v2 GetProject endpoint — the same endpoint
// the resource itself uses. This ensures we're exercising the v2 wire format,
// not the legacy v1 management API used in the existing tests.
func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			id := frame.State(state).ID
			// After a destroy the resource is gone from state and the ID
			// is empty. The v2 GetProject RPC rejects an empty ID with
			// InvalidArgument rather than NotFound, which the not-found
			// destroy assertion would not recognise. Treat an empty ID as
			// the project being absent.
			if id == "" {
				return test_utils.ErrNotFound
			}
			// Scope the verification call with the org id, exactly as the
			// resource's CRUD functions do via helper.CtxWithOrgID, so the
			// check exercises the same request context and does not depend
			// on server-side org defaults.
			ctx := helper.CtxSetOrgID(frame.Context, frame.OrgID)
			client, err := helper.GetProjectV2Client(ctx, frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get project v2 client: %w", err)
			}
			resp, err := client.GetProject(ctx, &projectpb.GetProjectRequest{
				ProjectId: id,
			})
			if err != nil {
				return err
			}
			actual := resp.GetProject().GetName()
			if actual != expect {
				return fmt.Errorf("expected project name %q, got %q", expect, actual)
			}
			return nil
		}
	}
}
