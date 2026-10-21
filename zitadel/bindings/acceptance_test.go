package bindings_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	authz "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authorization/v2"
	filter "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"
	permission "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/internal_permission/v2"
)

// These are Terraform CLI acceptance cases against a real server, not the
// in-process gRPC fault/validation tests in authorization and administrator.
func fixture(f *test_utils.BaseTestFrame) string {
	return f.ProviderSnippet + fmt.Sprintf(`
resource "zitadel_organization" "owner" { name = "%[1]s-owner" }
resource "zitadel_organization" "grantee" { name = "%[1]s-grantee" }
resource "zitadel_project_v2" "product" {
  org_id = zitadel_organization.owner.id
  name = "%[1]s-product"
}
resource "zitadel_human_user" "human" {
  org_id = zitadel_organization.owner.id
  user_name = "%[1]s-human"
  first_name = "Acceptance"
  last_name = "User"
  email = "%[1]s@example.com"
}
resource "zitadel_machine_user" "machine" {
  org_id = zitadel_organization.owner.id
  user_name = "%[1]s-machine"
  name = "Acceptance machine"
}
resource "zitadel_project_role" "reader" {
  org_id = zitadel_organization.owner.id
  project_id = zitadel_project_v2.product.id
  role_key = "reader"
  display_name = "Reader"
}
resource "zitadel_project_role" "writer" {
  org_id = zitadel_organization.owner.id
  project_id = zitadel_project_v2.product.id
  role_key = "writer"
  display_name = "Writer"
}
resource "zitadel_project_grant" "shared" {
  org_id = zitadel_organization.owner.id
  project_id = zitadel_project_v2.product.id
  granted_org_id = zitadel_organization.grantee.id
  role_keys = ["reader", "writer"]
  depends_on = [zitadel_project_role.reader, zitadel_project_role.writer]
}
`, f.UniqueResourcesID)
}

func sameRoles(got, want []string) bool {
	got = append([]string(nil), got...)
	want = append([]string(nil), want...)
	sort.Strings(got)
	sort.Strings(want)
	return reflect.DeepEqual(got, want)
}

func TestAccAuthorizationV2(t *testing.T) {
	for _, tc := range []struct{ name, subject, organization string }{
		{"human", "human", "owner"}, {"machine", "machine", "owner"}, {"delegated", "machine", "grantee"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := test_utils.NewBindingTestFrame(t, "zitadel_authorization")
			c, err := helper.GetAuthorizationV2Client(f, f.ClientInfo)
			if err != nil {
				t.Fatal(err)
			}
			userResource := "zitadel_" + tc.subject + "_user." + tc.subject
			config := func(role string) string {
				return fixture(f) + fmt.Sprintf(`
resource "zitadel_authorization" "default" {
  user_id = %s.id
  project_id = zitadel_project_v2.product.id
  org_id = zitadel_organization.%s.id
  role_keys = [%q]
  depends_on = [zitadel_project_role.reader, zitadel_project_role.writer, zitadel_project_grant.shared]
}
`, userResource, tc.organization, role)
			}
			var originalID string
			list := func(id string) ([]*authz.Authorization, error) {
				r, err := c.ListAuthorizations(f, &authz.ListAuthorizationsRequest{Pagination: &filter.PaginationRequest{Limit: 2}, Filters: []*authz.AuthorizationsSearchFilter{{Filter: &authz.AuthorizationsSearchFilter_AuthorizationIds{AuthorizationIds: &filter.InIDsFilter{Ids: []string{id}}}}}})
				if err != nil {
					return nil, err
				}
				return r.Authorizations, nil
			}
			check := func(role string) resource.TestCheckFunc {
				return func(s *terraform.State) error {
					r := s.RootModule().Resources["zitadel_authorization.default"].Primary
					if originalID == "" {
						originalID = r.ID
					}
					if r.ID != originalID {
						return fmt.Errorf("role update replaced assignment")
					}
					rows, err := list(r.ID)
					if err != nil {
						return err
					}
					if len(rows) != 1 {
						return fmt.Errorf("expected exactly one remote assignment, got %d", len(rows))
					}
					a := rows[0]
					roles := []string{}
					for _, r := range a.Roles {
						roles = append(roles, r.Key)
					}
					if a.User.Id != r.Attributes["user_id"] || a.Project.Id != r.Attributes["project_id"] || a.Organization.Id != r.Attributes["org_id"] || !sameRoles(roles, []string{role}) {
						return fmt.Errorf("remote authorization differs from Terraform state")
					}
					return nil
				}
			}
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: f.V6ProviderFactories(),
				CheckDestroy: func(*terraform.State) error {
					rows, err := list(originalID)
					if err != nil {
						return err
					}
					if len(rows) != 0 {
						return fmt.Errorf("remote assignment survived destroy")
					}
					return nil
				},
				Steps: []resource.TestStep{
					{Config: config("reader"), Check: check("reader")},
					{Config: config("writer"), Check: check("writer")},
					{ResourceName: "zitadel_authorization.default", ImportState: true, ImportStateVerify: true},
					{Config: config("writer"), PreConfig: func() {
						if _, err := c.UpdateAuthorization(f, &authz.UpdateAuthorizationRequest{Id: originalID, RoleKeys: []string{"reader"}}); err != nil {
							t.Fatal(err)
						}
					}, Check: check("writer")},
				},
			})
		})
	}
}

func TestAccAdministratorV2(t *testing.T) {
	for _, scope := range []string{"instance", "organization", "project", "project_grant"} {
		t.Run(scope, func(t *testing.T) {
			f := test_utils.NewBindingTestFrame(t, "zitadel_administrator")
			c, err := helper.GetInternalPermissionV2Client(f, f.ClientInfo)
			if err != nil {
				t.Fatal(err)
			}
			fields, first, second := "", "IAM_OWNER_VIEWER", "IAM_USER_MANAGER"
			switch scope {
			case "organization":
				fields, first, second = "org_id = zitadel_organization.owner.id", "ORG_USER_MANAGER", "ORG_OWNER_VIEWER"
			case "project":
				fields, first, second = "project_id = zitadel_project_v2.product.id", "PROJECT_OWNER_VIEWER", "PROJECT_OWNER"
			case "project_grant":
				fields, first, second = "project_id = zitadel_project_v2.product.id\norg_id = zitadel_organization.grantee.id", "PROJECT_GRANT_OWNER_VIEWER", "PROJECT_GRANT_OWNER"
			}
			config := func(role string) string {
				return fixture(f) + fmt.Sprintf(`
resource "zitadel_administrator" "default" {
  user_id = zitadel_machine_user.machine.id
  resource_type = %q
  %s
  roles = [%q]
  depends_on = [zitadel_project_grant.shared]
}
`, scope, fields, role)
			}
			var originalID, userID string
			var target *permission.ResourceFilter
			list := func() ([]*permission.Administrator, error) {
				r, err := c.ListAdministrators(f, &permission.ListAdministratorsRequest{Pagination: &filter.PaginationRequest{Limit: 2}, Filters: []*permission.AdministratorSearchFilter{
					{Filter: &permission.AdministratorSearchFilter_InUserIdsFilter{InUserIdsFilter: &filter.InIDsFilter{Ids: []string{userID}}}},
					{Filter: &permission.AdministratorSearchFilter_Resource{Resource: target}},
				}})
				if err != nil {
					return nil, err
				}
				return r.Administrators, nil
			}
			check := func(role string) resource.TestCheckFunc {
				return func(s *terraform.State) error {
					r := s.RootModule().Resources["zitadel_administrator.default"].Primary
					if originalID == "" {
						originalID = r.ID
					}
					if r.ID != originalID {
						return fmt.Errorf("role update replaced administrator identity")
					}
					userID = r.Attributes["user_id"]
					target = &permission.ResourceFilter{}
					switch scope {
					case "instance":
						target.Resource = &permission.ResourceFilter_Instance{Instance: true}
					case "organization":
						target.Resource = &permission.ResourceFilter_OrganizationId{OrganizationId: r.Attributes["org_id"]}
					case "project":
						target.Resource = &permission.ResourceFilter_ProjectId{ProjectId: r.Attributes["project_id"]}
					case "project_grant":
						target.Resource = &permission.ResourceFilter_ProjectGrant_{ProjectGrant: &permission.ResourceFilter_ProjectGrant{ProjectId: r.Attributes["project_id"], OrganizationId: r.Attributes["org_id"]}}
					}
					rows, err := list()
					if err != nil {
						return err
					}
					if len(rows) != 1 || rows[0].User.Id != userID || !sameRoles(rows[0].Roles, []string{role}) {
						return fmt.Errorf("remote administrator differs from Terraform state")
					}
					a := rows[0]
					matches := scope == "instance" && a.GetInstance() ||
						scope == "organization" && a.GetOrganization().GetId() == r.Attributes["org_id"] ||
						scope == "project" && a.GetProject().GetId() == r.Attributes["project_id"] ||
						scope == "project_grant" && a.GetProjectGrant().GetProjectId() == r.Attributes["project_id"] && a.GetProjectGrant().GetGrantedOrganizationId() == r.Attributes["org_id"]
					if !matches {
						return fmt.Errorf("remote administrator target differs from Terraform state")
					}
					return nil
				}
			}
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: f.V6ProviderFactories(),
				CheckDestroy: func(*terraform.State) error {
					rows, err := list()
					if err != nil {
						return err
					}
					if len(rows) != 0 {
						return fmt.Errorf("remote administrator survived destroy")
					}
					return nil
				},
				Steps: []resource.TestStep{
					{Config: config(first), Check: check(first)},
					{Config: config(second), Check: check(second)},
					{ResourceName: "zitadel_administrator.default", ImportState: true, ImportStateVerify: true},
				},
			})
		})
	}
}
