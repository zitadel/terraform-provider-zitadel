package administrator

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/bindingtest"
	permission "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/internal_permission/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type administratorServer struct {
	permission.UnimplementedInternalPermissionServiceServer
	t           *testing.T
	target      target
	roles       []string
	present     bool
	readErr     error
	writeErr    error
	deleteErr   error
	createCount int
	wrongResult bool
}

func (s *administratorServer) CreateAdministrator(_ context.Context, r *permission.CreateAdministratorRequest) (*permission.CreateAdministratorResponse, error) {
	s.createCount++
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	if r.UserId != s.target.user || !proto.Equal(r.Resource, s.target.resource()) {
		s.t.Errorf("create targeted wrong tuple: %v", r)
	}
	s.roles, s.present = r.Roles, true
	return &permission.CreateAdministratorResponse{}, nil
}

func (s *administratorServer) UpdateAdministrator(_ context.Context, r *permission.UpdateAdministratorRequest) (*permission.UpdateAdministratorResponse, error) {
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	if r.UserId != s.target.user || !proto.Equal(r.Resource, s.target.resource()) {
		s.t.Errorf("update targeted wrong tuple: %v", r)
	}
	s.roles = r.Roles
	return &permission.UpdateAdministratorResponse{}, nil
}

func (s *administratorServer) DeleteAdministrator(_ context.Context, r *permission.DeleteAdministratorRequest) (*permission.DeleteAdministratorResponse, error) {
	if r.UserId != s.target.user || !proto.Equal(r.Resource, s.target.resource()) {
		s.t.Errorf("delete targeted wrong tuple: %v", r)
	}
	if s.deleteErr != nil {
		return nil, s.deleteErr
	}
	s.present = false
	return &permission.DeleteAdministratorResponse{}, nil
}

func (s *administratorServer) ListAdministrators(_ context.Context, r *permission.ListAdministratorsRequest) (*permission.ListAdministratorsResponse, error) {
	if len(r.Filters) != 2 || len(r.Filters[0].GetInUserIdsFilter().GetIds()) != 1 || r.Filters[0].GetInUserIdsFilter().GetIds()[0] != s.target.user || !proto.Equal(r.Filters[1].GetResource(), s.target.filter()) || r.GetPagination().GetLimit() != 2 {
		s.t.Errorf("read must select exact user/target tuple: %v", r)
	}
	if s.readErr != nil {
		return nil, s.readErr
	}
	resp := &permission.ListAdministratorsResponse{}
	if !s.present {
		return resp, nil
	}
	a := &permission.Administrator{User: &permission.User{Id: s.target.user}, Roles: s.roles}
	switch s.target.scope {
	case "instance":
		a.Resource = &permission.Administrator_Instance{Instance: true}
	case "organization":
		a.Resource = &permission.Administrator_Organization{Organization: &permission.Organization{Id: s.target.org}}
	case "project":
		a.Resource = &permission.Administrator_Project{Project: &permission.Project{Id: s.target.project}}
	case "project_grant":
		a.Resource = &permission.Administrator_ProjectGrant{ProjectGrant: &permission.ProjectGrant{ProjectId: s.target.project, OrganizationId: "project-owner", GrantedOrganizationId: s.target.org}}
	}
	if s.wrongResult {
		a.User.Id = "another-user"
	}
	resp.Administrators = []*permission.Administrator{a}
	return resp, nil
}

func testResource(t *testing.T, s *administratorServer) *schema.Resource {
	t.Helper()
	s.t = t
	conn := bindingtest.BindingTestConnection(t, func(server *grpc.Server) { permission.RegisterInternalPermissionServiceServer(server, s) })
	return newResource(func(context.Context, interface{}) (permission.InternalPermissionServiceClient, error) {
		return permission.NewInternalPermissionServiceClient(conn), nil
	})
}

func TestAdministratorV2LifecycleAndImport(t *testing.T) {
	for _, scope := range []target{
		{user: "user", scope: "instance"},
		{user: "user", scope: "organization", org: "org"},
		{user: "user", scope: "project", project: "project"},
		{user: "user", scope: "project_grant", project: "project", org: "grantee"},
	} {
		t.Run(scope.scope, func(t *testing.T) {
			s := &administratorServer{target: scope}
			r := testResource(t, s)
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"user_id": scope.user, "resource_type": scope.scope, "project_id": scope.project, "org_id": scope.org, "roles": []interface{}{"ROLE_A", "ROLE_B"}})
			ctx := context.Background()
			if diags := r.CreateContext(ctx, d, nil); diags.HasError() {
				t.Fatal(diags)
			}
			if d.Id() != scope.identity() {
				t.Fatalf("unexpected tuple ID: %s", d.Id())
			}
			if err := d.Set(RolesVar, []string{"ROLE_A"}); err != nil {
				t.Fatal(err)
			}
			if diags := r.UpdateContext(ctx, d, nil); diags.HasError() {
				t.Fatal(diags)
			}
			if !reflect.DeepEqual(s.roles, []string{"ROLE_A"}) {
				t.Fatal("update did not replace complete role set")
			}
			imported := schema.TestResourceDataRaw(t, r.Schema, nil)
			imported.SetId(d.Id())
			if _, err := r.Importer.StateContext(ctx, imported, nil); err != nil {
				t.Fatal(err)
			}
			if diags := r.ReadContext(ctx, imported, nil); diags.HasError() {
				t.Fatal(diags)
			}
			if got := targetFrom(imported); got != scope || s.createCount != 1 {
				t.Fatalf("import changed identity or recreated assignment: %v", got)
			}
			// Read restores drifted roles and never changes the target to whichever
			// values happen to be supplied alongside an imported external ID.
			s.roles = []string{"EXTERNAL_ROLE"}
			if err := d.Set(UserIDVar, "new-spec-user"); err != nil {
				t.Fatal(err)
			}
			if diags := r.ReadContext(ctx, d, nil); diags.HasError() {
				t.Fatal(diags)
			}
			if targetFrom(d) != scope || !reflect.DeepEqual(d.Get(RolesVar).(*schema.Set).List(), []interface{}{"EXTERNAL_ROLE"}) {
				t.Fatal("read did not restore remote tuple/roles")
			}
			for _, code := range []codes.Code{codes.PermissionDenied, codes.Unavailable, codes.NotFound} {
				s.readErr = status.Error(code, "read failed")
				if diags := r.ReadContext(ctx, d, nil); !diags.HasError() || d.Id() == "" {
					t.Fatalf("read error %s lost identity", code)
				}
			}
			s.readErr = nil
			s.wrongResult = true
			if diags := r.ReadContext(ctx, d, nil); !diags.HasError() || d.Id() == "" {
				t.Fatal("wrong tuple accepted")
			}
			s.wrongResult = false
			s.deleteErr = status.Error(codes.PermissionDenied, "denied")
			if diags := r.DeleteContext(ctx, d, nil); !diags.HasError() || d.Id() == "" {
				t.Fatal("delete error lost identity")
			}
			s.deleteErr = nil
			if diags := r.DeleteContext(ctx, d, nil); diags.HasError() || d.Id() != "" {
				t.Fatalf("delete failed: %v", diags)
			}
			if diags := r.ReadContext(ctx, imported, nil); diags.HasError() || imported.Id() != "" {
				t.Fatalf("successful empty list did not detect deletion: %v", diags)
			}
		})
	}
}

func TestAdministratorIdentityValidation(t *testing.T) {
	for _, id := range []string{"", "user", "user:unknown", "user:instance:org", "user:organization", "user:organization:", "user:project:project:org", "user:project_grant:project", ":instance", "%XX:instance"} {
		if _, err := parseIdentity(id); err == nil {
			t.Errorf("invalid identity accepted: %q", id)
		}
	}
	want := target{user: "custom:user+name", scope: "project_grant", project: "custom:project", org: "custom/org"}
	got, err := parseIdentity(want.identity())
	if err != nil || got != want {
		t.Fatalf("escaped identity did not round-trip: %v, %v", got, err)
	}
	for _, invalid := range []target{
		{user: "u", scope: "instance", org: "o"}, {user: "u", scope: "instance", project: "p"},
		{user: "u", scope: "organization"}, {user: "u", scope: "organization", org: "o", project: "p"},
		{user: "u", scope: "project"}, {user: "u", scope: "project", project: "p", org: "o"},
		{user: "u", scope: "project_grant", project: "p"}, {user: "u", scope: "project_grant", org: "o"},
		{scope: "instance"}, {user: "u", scope: "unknown"},
	} {
		if err := invalid.validate(); err == nil {
			t.Errorf("invalid target accepted: %v", invalid)
		}
	}
}

func TestAdministratorCreateDoesNotImplicitlyAdopt(t *testing.T) {
	s := &administratorServer{target: target{user: "user", scope: "instance"}, writeErr: status.Error(codes.AlreadyExists, "duplicate")}
	r := testResource(t, s)
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"user_id": "user", "resource_type": "instance", "roles": []interface{}{"IAM_OWNER"}})
	if diags := r.CreateContext(context.Background(), d, nil); !diags.HasError() || d.Id() != "" {
		t.Fatal("duplicate implicitly adopted")
	}
	if err := d.Set("org_id", "unexpected"); err != nil {
		t.Fatal(err)
	}
	if diags := r.CreateContext(context.Background(), d, nil); !diags.HasError() || s.createCount != 1 {
		t.Fatal("invalid scope reached the API")
	}
}
