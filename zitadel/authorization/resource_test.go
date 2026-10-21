package authorization

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/bindingtest"
	authz "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authorization/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorizationServer struct {
	authz.UnimplementedAuthorizationServiceServer
	t           *testing.T
	assignment  *authz.Authorization
	readErr     error
	writeErr    error
	deleteErr   error
	createCount int
	readCount   int
	emptyReads  int
}

func (s *authorizationServer) CreateAuthorization(_ context.Context, r *authz.CreateAuthorizationRequest) (*authz.CreateAuthorizationResponse, error) {
	s.createCount++
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	s.assignment = &authz.Authorization{Id: "existing-authorization", User: &authz.User{Id: r.UserId}, Project: &authz.Project{Id: r.ProjectId}, Organization: &authz.Organization{Id: r.OrganizationId}, State: authz.State_STATE_ACTIVE}
	for _, key := range r.RoleKeys {
		s.assignment.Roles = append(s.assignment.Roles, &authz.Role{Key: key})
	}
	return &authz.CreateAuthorizationResponse{Id: s.assignment.Id}, nil
}

func (s *authorizationServer) UpdateAuthorization(_ context.Context, r *authz.UpdateAuthorizationRequest) (*authz.UpdateAuthorizationResponse, error) {
	if s.writeErr != nil {
		return nil, s.writeErr
	}
	if s.assignment == nil || r.Id != s.assignment.Id {
		s.t.Errorf("update targeted wrong ID: %v", r)
	}
	s.assignment.Roles = nil
	for _, key := range r.RoleKeys {
		s.assignment.Roles = append(s.assignment.Roles, &authz.Role{Key: key})
	}
	return &authz.UpdateAuthorizationResponse{}, nil
}

func (s *authorizationServer) ListAuthorizations(_ context.Context, r *authz.ListAuthorizationsRequest) (*authz.ListAuthorizationsResponse, error) {
	s.readCount++
	if len(r.Filters) != 1 || len(r.Filters[0].GetAuthorizationIds().GetIds()) != 1 || r.Filters[0].GetAuthorizationIds().GetIds()[0] != "existing-authorization" || r.GetPagination().GetLimit() != 2 {
		s.t.Errorf("read must select exact authorization ID: %v", r)
	}
	if s.readErr != nil {
		return nil, s.readErr
	}
	resp := &authz.ListAuthorizationsResponse{}
	if s.readCount <= s.emptyReads {
		return resp, nil
	}
	if s.assignment != nil {
		resp.Authorizations = []*authz.Authorization{s.assignment}
	}
	return resp, nil
}

func (s *authorizationServer) DeleteAuthorization(_ context.Context, r *authz.DeleteAuthorizationRequest) (*authz.DeleteAuthorizationResponse, error) {
	if r.Id != "existing-authorization" {
		s.t.Errorf("delete targeted wrong ID: %v", r)
	}
	if s.deleteErr != nil {
		return nil, s.deleteErr
	}
	s.assignment = nil
	return &authz.DeleteAuthorizationResponse{}, nil
}

func testResource(t *testing.T, s *authorizationServer) *schema.Resource {
	t.Helper()
	s.t = t
	conn := bindingtest.BindingTestConnection(t, func(server *grpc.Server) { authz.RegisterAuthorizationServiceServer(server, s) })
	r := newResource(func(context.Context, interface{}) (authz.AuthorizationServiceClient, error) {
		return authz.NewAuthorizationServiceClient(conn), nil
	})
	return r
}

func TestAuthorizationV2LifecycleAndImport(t *testing.T) {
	s := &authorizationServer{}
	r := testResource(t, s)
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"user_id": "human-or-machine", "project_id": "project", "org_id": "grantee", "role_keys": []interface{}{"reader", "writer"}})
	ctx := context.Background()
	if diags := r.CreateContext(ctx, d, nil); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Id() != "existing-authorization" {
		t.Fatalf("wrong ID: %s", d.Id())
	}
	if got := d.Get("state"); got != "STATE_ACTIVE" {
		t.Fatalf("state immediately after create = %v, want STATE_ACTIVE", got)
	}
	if s.assignment.GetOrganization().GetId() != "grantee" {
		t.Fatal("authorization organization was not explicit")
	}
	if err := d.Set(RoleKeysVar, []string{"reader"}); err != nil {
		t.Fatal(err)
	}
	if diags := r.UpdateContext(ctx, d, nil); diags.HasError() {
		t.Fatal(diags)
	}
	if len(s.assignment.Roles) != 1 || s.assignment.Roles[0].Key != "reader" {
		t.Fatal("update did not replace complete role set")
	}
	// Import supplies only the real ID; v2 read hydrates every immutable field.
	imported := schema.TestResourceDataRaw(t, r.Schema, nil)
	imported.SetId(d.Id())
	if _, err := r.Importer.StateContext(ctx, imported, nil); err != nil {
		t.Fatal(err)
	}
	if diags := r.ReadContext(ctx, imported, nil); diags.HasError() {
		t.Fatal(diags)
	}
	for key, want := range map[string]string{"user_id": "human-or-machine", "project_id": "project", "org_id": "grantee", "state": "STATE_ACTIVE"} {
		if imported.Get(key) != want {
			t.Errorf("%s = %v, want %v", key, imported.Get(key), want)
		}
	}
	if s.createCount != 1 {
		t.Fatal("import recreated authorization")
	}
	// Detect external role edits, including a role removed outside Terraform.
	s.assignment.Roles = []*authz.Role{{Key: "external-role"}}
	if diags := r.ReadContext(ctx, d, nil); diags.HasError() {
		t.Fatal(diags)
	}
	got := d.Get(RoleKeysVar).(*schema.Set).List()
	sort.Slice(got, func(i, j int) bool { return got[i].(string) < got[j].(string) })
	if !reflect.DeepEqual(got, []interface{}{"external-role"}) {
		t.Fatalf("drift not refreshed: %v", got)
	}
	// A failed read must preserve the ID; even NotFound on the list RPC can
	// mean a missing parent/context rather than this exact assignment.
	for _, code := range []codes.Code{codes.PermissionDenied, codes.Unavailable, codes.NotFound} {
		s.readErr = status.Error(code, "read failed")
		if diags := r.ReadContext(ctx, d, nil); !diags.HasError() || d.Id() == "" {
			t.Fatalf("read error %s lost identity: %v", code, diags)
		}
	}
	s.readErr = nil
	s.deleteErr = status.Error(codes.PermissionDenied, "delete denied")
	if diags := r.DeleteContext(ctx, d, nil); !diags.HasError() || d.Id() == "" {
		t.Fatal("failed deletion lost identity")
	}
	s.deleteErr = nil
	if diags := r.DeleteContext(ctx, d, nil); diags.HasError() || d.Id() != "" {
		t.Fatalf("delete failed: %v", diags)
	}
	imported.SetId("existing-authorization")
	if diags := r.ReadContext(ctx, imported, nil); diags.HasError() || imported.Id() != "" {
		t.Fatalf("successful empty list did not detect deletion: %v", diags)
	}
}

func TestAuthorizationV2AmbiguousReadAndCreateErrors(t *testing.T) {
	s := &authorizationServer{assignment: &authz.Authorization{Id: "other-id"}}
	r := testResource(t, s)
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"user_id": "user", "project_id": "project", "org_id": "org"})
	d.SetId("existing-authorization")
	if diags := r.ReadContext(context.Background(), d, nil); !diags.HasError() || d.Id() != "existing-authorization" {
		t.Fatal("wrong read result accepted")
	}
	d.SetId("")
	s.writeErr = status.Error(codes.AlreadyExists, "duplicate")
	if diags := r.CreateContext(context.Background(), d, nil); !diags.HasError() || d.Id() != "" {
		t.Fatal("duplicate implicitly adopted")
	}
	if _, err := r.Importer.StateContext(context.Background(), d, nil); err == nil {
		t.Fatal("empty import ID accepted")
	}
}

func TestAuthorizationV2CreateRefresh(t *testing.T) {
	for _, tc := range []struct {
		name       string
		emptyReads int
		readErr    error
		timeout    time.Duration
		wantError  bool
	}{
		{name: "delayed visibility", emptyReads: 2},
		{name: "permission denied", readErr: status.Error(codes.PermissionDenied, "read denied"), wantError: true},
		{name: "unavailable", readErr: status.Error(codes.Unavailable, "read unavailable"), wantError: true},
		{name: "visibility deadline", emptyReads: 100, timeout: 100 * time.Millisecond, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &authorizationServer{emptyReads: tc.emptyReads, readErr: tc.readErr}
			r := testResource(t, s)
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"user_id": "user", "project_id": "project", "org_id": "org", "role_keys": []interface{}{"reader"}})
			ctx := context.Background()
			if tc.timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}
			diags := r.CreateContext(ctx, d, nil)
			if diags.HasError() != tc.wantError {
				t.Fatalf("create diagnostics = %v, want error %v", diags, tc.wantError)
			}
			if d.Id() != "existing-authorization" || s.createCount != 1 {
				t.Fatalf("refresh lost or recreated authorization: ID=%q, creates=%d", d.Id(), s.createCount)
			}
			if !tc.wantError {
				if d.Get("state") != "STATE_ACTIVE" || s.readCount != tc.emptyReads+1 {
					t.Fatalf("state=%v after %d reads", d.Get("state"), s.readCount)
				}
				if !reflect.DeepEqual(d.Get(RoleKeysVar).(*schema.Set).List(), []interface{}{"reader"}) {
					t.Fatal("create did not hydrate role keys")
				}
			}
		})
	}
}
