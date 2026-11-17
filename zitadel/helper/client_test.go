package helper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/settings/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/oidctest"
)

func TestIgnoreIfNotFoundError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "nil error returns nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:    "not found error returns nil",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: true,
		},
		{
			name:    "already exists error is not ignored",
			err:     status.Error(codes.AlreadyExists, "already exists"),
			wantNil: false,
		},
		{
			name:    "internal error is not ignored",
			err:     status.Error(codes.Internal, "internal error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnoreIfNotFoundError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnoreIfNotFoundError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnoreIfNotFoundError() = nil, want non-nil")
			}
		})
	}
}

func TestIgnorePreconditionError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "precondition error returns nil",
			err:     status.Error(codes.FailedPrecondition, "failed precondition"),
			wantNil: true,
		},
		{
			name:    "not found error is not ignored",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: false,
		},
		{
			name:    "internal error is not ignored",
			err:     status.Error(codes.Internal, "internal error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnorePreconditionError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnorePreconditionError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnorePreconditionError() = nil, want non-nil")
			}
		})
	}
}

func TestIgnoreAlreadyExistsError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "already exists error returns nil",
			err:     status.Error(codes.AlreadyExists, "already exists"),
			wantNil: true,
		},
		{
			name:    "not found error is not ignored",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: false,
		},
		{
			name:    "plain error is not ignored",
			err:     fmt.Errorf("some error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnoreAlreadyExistsError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnoreAlreadyExistsError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnoreAlreadyExistsError() = nil, want non-nil")
			}
		})
	}
}

// jwtProfileMode is one of the provider settings that carry a JWT-profile key.
type jwtProfileMode struct {
	name string
	args func(t *testing.T) (token, profileFile, profileJSON string)
}

func jwtProfileModes() []jwtProfileMode {
	return []jwtProfileMode{
		{"token", func(t *testing.T) (string, string, string) {
			return writeKeyFile(t, oidctest.KeyFileJSON(t)), "", ""
		}},
		{"jwt_profile_file", func(t *testing.T) (string, string, string) {
			return "", writeKeyFile(t, oidctest.KeyFileJSON(t)), ""
		}},
		{"jwt_profile_json", func(t *testing.T) (string, string, string) {
			return "", "", string(oidctest.KeyFileJSON(t))
		}},
	}
}

// clientInfoFor returns the client info of a provider configured with the mode and
// pointed at the test server.
func clientInfoFor(t *testing.T, idp *oidctest.Server, mode jwtProfileMode, skipTLS bool, headers map[string]string) *ClientInfo {
	t.Helper()
	host, port := hostPort(t, idp.URL)
	token, profileFile, profileJSON := mode.args(t)
	info, err := GetClientInfo(context.Background(), !skipTLS, host, "", token, "", profileFile, profileJSON, "", "", "", "", "", "", port, skipTLS, headers)
	if err != nil {
		t.Fatalf("GetClientInfo() error = %v", err)
	}
	return info
}

// TestGetClientInfoJWTProfileOIDCRequests checks that the transport_headers and
// insecure_skip_verify_tls reach the OIDC discovery and token requests of every
// JWT-profile mode, not only the gRPC calls.
func TestGetClientInfoJWTProfileOIDCRequests(t *testing.T) {
	const headerName, headerValue = "X-Forwarded-Host", "zitadel.example.com"

	servers := []struct {
		name    string
		useTLS  bool
		headers map[string]string
		skipTLS bool
	}{
		{name: "transport headers", headers: map[string]string{headerName: headerValue}},
		{name: "transport headers over self-signed TLS", useTLS: true, skipTLS: true, headers: map[string]string{headerName: headerValue}},
		{name: "skip verify only over self-signed TLS", useTLS: true, skipTLS: true},
		{name: "no headers"},
	}
	for _, mode := range jwtProfileModes() {
		for _, srv := range servers {
			t.Run(mode.name+"/"+srv.name, func(t *testing.T) {
				idp := oidctest.New(t, srv.useTLS)
				info := clientInfoFor(t, idp, mode, srv.skipTLS, srv.headers)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				client, err := settingsv2.NewClient(ctx, info.Issuer, info.Domain, []string{"openid"}, info.Options...)
				if err != nil {
					t.Fatalf("NewClient() error = %v", err)
				}
				defer client.Connection.Close()
				// The test server is not a gRPC server, so the call itself fails; the
				// auth interceptor fetches the token before sending it.
				_ = client.Connection.Invoke(ctx, "/test.Service/Method", &emptypb.Empty{}, &emptypb.Empty{})

				want := srv.headers[headerName]
				for _, path := range []string{oidctest.DiscoveryPath, oidctest.TokenPath} {
					if !idp.Requested(path) {
						t.Fatalf("%s was never requested", path)
					}
					if got := idp.Header(path, headerName); got != want {
						t.Errorf("%s header %s = %q, want %q", path, headerName, got, want)
					}
				}
			})
		}
	}
}

func TestJWTProfileTokenSource_WithoutAKey_ReturnsErrNoJWTProfile(t *testing.T) {
	_, err := (&ClientInfo{}).JWTProfileTokenSource(context.Background())
	if !errors.Is(err, ErrNoJWTProfile) {
		t.Errorf("JWTProfileTokenSource() error = %v, want %v", err, ErrNoJWTProfile)
	}
}

func writeKeyFile(t *testing.T, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "key.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
