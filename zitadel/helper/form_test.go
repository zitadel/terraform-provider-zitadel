package helper

import (
	"context"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/oidctest"
)

func writeAsset(t *testing.T) string {
	t.Helper()
	asset := filepath.Join(t.TempDir(), "logo.png")
	if err := os.WriteFile(asset, []byte("\x89PNG\r\n\x1a\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return asset
}

// TestInstanceFormFilePostOIDCRequests checks that the transport_headers and
// insecure_skip_verify_tls reach the OIDC discovery and token requests made for an
// asset upload, for each JWT-profile mode.
func TestInstanceFormFilePostOIDCRequests(t *testing.T) {
	const headerName, headerValue = "X-Forwarded-Host", "zitadel.example.com"
	asset := writeAsset(t)

	servers := []struct {
		name    string
		useTLS  bool
		skipTLS bool
		headers map[string]string
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

				if diags := InstanceFormFilePost(context.Background(), info, oidctest.UploadPath, asset); diags.HasError() {
					t.Fatalf("InstanceFormFilePost() = %v", diags)
				}

				want := srv.headers[headerName]
				for _, path := range []string{oidctest.DiscoveryPath, oidctest.TokenPath, oidctest.UploadPath} {
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

// TestInstanceFormFilePostHostHeader checks the instance that is reached on an internal
// address with its external host as Host header: the OIDC requests and the upload are
// made for that host, and the issuer is the one the instance announces for it.
func TestInstanceFormFilePostHostHeader(t *testing.T) {
	const external = "zitadel.example.com"
	asset := writeAsset(t)

	for _, mode := range jwtProfileModes() {
		for _, useTLS := range []bool{false, true} {
			name := mode.name + "/plain"
			if useTLS {
				name = mode.name + "/tls"
			}
			t.Run(name, func(t *testing.T) {
				idp := oidctest.New(t, useTLS)
				info := clientInfoFor(t, idp, mode, useTLS, map[string]string{"Host": external})

				if diags := InstanceFormFilePost(context.Background(), info, oidctest.UploadPath, asset); diags.HasError() {
					t.Fatalf("InstanceFormFilePost() = %v", diags)
				}

				for _, path := range []string{oidctest.DiscoveryPath, oidctest.TokenPath, oidctest.UploadPath} {
					if got := idp.Host(path); got != external {
						t.Errorf("%s host = %q, want %q", path, got, external)
					}
				}
			})
		}
	}
}

// hostPort splits the host and port of the URL of a test server.
func hostPort(t *testing.T, rawURL string) (string, string) {
	t.Helper()
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatal(err)
	}
	return host, port
}
