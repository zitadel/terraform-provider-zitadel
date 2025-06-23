package action_target_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2beta"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccTarget(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_action_target")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)

	nameAttribute := test_utils.AttributeValue(t, action_target.NameVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, nameAttribute, frame.UniqueResourcesID, 1)

	exampleProperty := test_utils.AttributeValue(t, action_target.EndpointVar, exampleAttributes).AsString()
	updatedProperty := exampleProperty + "-updated"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty,
		updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame),
		test_utils.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expectedEndpoint string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			ctx := metadata.AppendToOutgoingContext(context.Background(), "x-zitadel-orgid", frame.OrgID)
			remoteResource, err := client.GetTarget(
				ctx,
				&action.GetTargetRequest{Id: frame.State(state).ID},
			)
			if err != nil {
				return err
			}
			actualEndpoint := remoteResource.GetTarget().GetEndpoint()
			if actualEndpoint != expectedEndpoint {
				return fmt.Errorf("expected endpoint %q, but got %q", expectedEndpoint, actualEndpoint)
			}
			return nil
		}
	}
}
