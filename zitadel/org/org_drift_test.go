package org_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org"
)

func TestAccOrg404Handling(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org.NameVar, exampleAttributes).AsString()
	orgName := "test404org_" + frame.UniqueResourcesID
	resourceExample = strings.Replace(resourceExample, exampleProperty, orgName, 1)

	var orgID string

	resource.Test(t, resource.TestCase{
		ProviderFactories: frame.BaseTestFrame.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceExample,
				Check: resource.ComposeTestCheckFunc(
					func(state *terraform.State) error {
						orgID = frame.State(state).ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					ctx := context.Background()
					adminClient, err := helper.GetAdminClient(ctx, frame.ClientInfo)
					if err != nil {
						t.Fatalf("failed to get admin client: %v", err)
					}
					_, err = adminClient.RemoveOrg(ctx, &admin.RemoveOrgRequest{OrgId: orgID})
					if err != nil {
						t.Fatalf("failed to delete org: %v", err)
					}
					t.Logf("Deleted org %s outside Terraform", orgID)
				},
				Config:             resourceExample,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
