package organization_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/organization"
)

func TestAccOrganization(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_organization")

	userDep := fmt.Sprintf(`
resource "zitadel_human_user" "dep" {
  org_id            = "%s"
  user_name         = "test-admin-%s@zitadel.com"
  first_name        = "Test"
  last_name         = "Admin"
  email             = "test-admin-%s@zitadel.com"
  is_email_verified = true
  initial_password  = "Password1!"
}
`, frame.OrgID, frame.UniqueResourcesID, frame.UniqueResourcesID)

	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, organization.NameVar, exampleAttributes).AsString()
	initialProperty := "initialorgname_" + frame.UniqueResourcesID
	updatedProperty := "updatedorgname_" + frame.UniqueResourcesID

	resourceFunc := func(property string, secret string) string {
		content := test_utils.ReplaceAll(resourceExample, exampleProperty, "")(property, secret)
		return strings.ReplaceAll(content, "\"123456789012345678\"", "zitadel_human_user.dep.id")
	}

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{userDep},
		resourceFunc,
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, idFromState(frame)),
		helper.ZitadelGeneratedIdOnlyRegex,
		orgGone(frame, idFromState(frame)),
		test_utils.ImportResourceId(frame.BaseTestFrame),
	)
}

func TestAccOrganizationDatasource_ID(t *testing.T) {
	datasourceName := "zitadel_organization"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleID := test_utils.AttributeValue(t, organization.OrgIDVar, attributes).AsString()
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame := frame.AnotherOrg(t, orgName)
	config = strings.Replace(config, exampleID, otherFrame.OrgID, 1)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"id":    otherFrame.OrgID,
			"name":  orgName,
			"state": "ORGANIZATION_STATE_ACTIVE",
		},
	)
}

func TestAccOrganizationsDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_organizations"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, organization.NameVar, attributes).AsString()
	exampleDomain := test_utils.AttributeValue(t, organization.DomainVar, attributes).AsString()
	orgName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	config = strings.Join(strings.Split(config, "\n")[0:7], "\n")
	config = strings.Replace(config, exampleName, orgName, 1)
	config = strings.Replace(config, exampleDomain, orgName, 1)
	otherFrame := frame.AnotherOrg(t, orgName)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		config,
		nil,
		checkRemoteProperty(otherFrame, idFromFrame(otherFrame))(orgName),
		map[string]string{
			"ids.0": otherFrame.OrgID,
			"ids.#": "1",
		},
	)
}

func TestAccOrganizationsDatasources_ID_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_organizations"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	orgName := fmt.Sprintf("%s_%s", test_utils.AttributeValue(t, organization.NameVar, attributes), frame.UniqueResourcesID)
	otherFrame := frame.AnotherOrg(t, orgName)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		config,
		nil,
		checkRemoteProperty(otherFrame, idFromFrame(otherFrame))(orgName),
		map[string]string{"ids.#": "0"},
	)
}

func idFromState(frame *test_utils.OrgTestFrame) func(*terraform.State) string {
	return func(state *terraform.State) string {
		return frame.State(state).ID
	}
}

func idFromFrame(frame *test_utils.OrgTestFrame) func(state *terraform.State) string {
	return func(state *terraform.State) string {
		return frame.OrgID
	}
}

// checkRemoteProperty checks if the organization exists and has the expected name using V2 API
func checkRemoteProperty(frame *test_utils.OrgTestFrame, id func(state *terraform.State) string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetOrgClient(frame.Context, frame.ClientInfo)
			if err != nil {
				return err
			}
			orgID := id(state)
			resp, err := client.ListOrganizations(frame.Context, &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{
					{
						Query: &org.SearchQuery_IdQuery{
							IdQuery: &org.OrganizationIDQuery{
								Id: orgID,
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
			if len(resp.Result) == 0 {
				return fmt.Errorf("org %s not found", orgID)
			}
			remoteResource := resp.Result[0]
			actual := remoteResource.GetName()
			// Check if state is removed/deleted
			if remoteResource.GetState() == org.OrganizationState_ORGANIZATION_STATE_REMOVED {
				return fmt.Errorf("org is removed: %w", test_utils.ErrNotFound)
			}
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}

// orgGone verifies the org is effectively gone after destroy.
func orgGone(
	frame *test_utils.OrgTestFrame,
	idFromState func(*terraform.State) string,
) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		orgID := idFromState(state)
		ctx := frame.Context

		client, err := helper.GetOrgClient(ctx, frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("get org client: %w", err)
		}

		resp, err := client.ListOrganizations(ctx, &org.ListOrganizationsRequest{
			Queries: []*org.SearchQuery{
				{
					Query: &org.SearchQuery_IdQuery{
						IdQuery: &org.OrganizationIDQuery{
							Id: orgID,
						},
					},
				},
			},
		})
		if err != nil {
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.NotFound, codes.PermissionDenied, codes.Unauthenticated:
				return nil
			}
			return fmt.Errorf("unexpected error after destroy (id=%s): %v", orgID, err)
		}

		if len(resp.Result) == 0 {
			return nil
		}

		remoteOrg := resp.Result[0]
		// Call succeeded: consider it "gone" if it's REMOVED
		if remoteOrg.State == org.OrganizationState_ORGANIZATION_STATE_REMOVED {
			return nil
		}

		return fmt.Errorf("expected org to be gone, but it's still present/active (id=%s, state=%s)", orgID, remoteOrg.State)
	}
}
