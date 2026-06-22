package helper

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/zitadel/oidc/v3/pkg/client/profile"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	"golang.org/x/oauth2"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createMultipartRequest(issuer, endpoint, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset: %v", err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes("file"), escapeQuotes(filepath.Base(file.Name()))))
	h.Set("Content-Type", mimetype.Detect(data).String())
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset part: %v", err)
	}
	io.Copy(part, bytes.NewBuffer(data))
	writer.Close()

	r, err := http.NewRequest(http.MethodPost, issuer+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset request: %v", err)
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	return r, nil
}

// InstanceFormFilePost uploads an asset (logo, icon or font) to an instance-level
// endpoint. Instance-scoped assets need no organization header.
func InstanceFormFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path string) diag.Diagnostics {
	return formFilePost(ctx, clientInfo, endpoint, path, map[string]string{})
}

// OrgFormFilePost uploads an asset to an org-level endpoint. The x-zitadel-orgid
// header selects which organization the asset belongs to, mirroring the
// per-request org context the gRPC client sets via helper.CtxSetOrgID.
func OrgFormFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path, orgID string) diag.Diagnostics {
	return formFilePost(ctx, clientInfo, endpoint, path, map[string]string{"x-zitadel-orgid": orgID})
}

// formFilePost performs the multipart asset upload. ZITADEL serves asset uploads
// over a plain HTTP endpoint that sits outside the gRPC API, so this path builds
// its own authenticated *http.Client instead of reusing the gRPC connection. It
// must therefore reproduce, by hand, the two things the gRPC stack does
// automatically: attach the provider's transport_headers and authenticate the
// request with the configured credential.
func formFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path string, additionalHeaders map[string]string) diag.Diagnostics {
	var client *http.Client

	// Build the multipart body (the asset file) and the base request.
	r, err := createMultipartRequest(clientInfo.Issuer, endpoint, path)
	if err != nil {
		return diag.Errorf("failed to create asset request: %v", err)
	}
	// Bind the caller's context so a Terraform cancellation or timeout actually
	// interrupts the in-flight upload.
	r = r.WithContext(ctx)

	// createMultipartRequest set Content-Type to the multipart boundary; remember
	// it so the header loops below cannot clobber it (a stray Content-Type in
	// transport_headers would otherwise corrupt the request body framing).
	multipartContentType := r.Header.Get("Content-Type")

	// Apply the provider's transport_headers first, then the per-request headers,
	// so a caller-supplied header such as x-zitadel-orgid can never be overridden
	// by a transport header of the same name. Set (not Add) keeps single-valued
	// headers deterministic.
	for k, v := range clientInfo.TransportHeaders {
		r.Header.Set(k, v)
	}
	for k, v := range additionalHeaders {
		r.Header.Set(k, v)
	}
	r.Header.Set("Content-Type", multipartContentType)

	// The asset-upload transport must honor the provider's insecure_skip_verify_tls
	// the same way the gRPC client does, since this path does not go through gRPC.
	core := assetTransport(clientInfo.InsecureSkipVerifyTLS)

	// Pick an authenticated HTTP client that matches the provider's auth mode.
	// The order mirrors how the credential is stored on ClientInfo.
	switch {
	case clientInfo.TokenSource != nil:
		// access_token (PAT), jwt_file and system_api: a ready bearer token source.
		client = NewClientWithInterceptor(clientInfo.TokenSource, core)
	case clientInfo.KeyPath != "":
		// jwt_profile_file / token: a JWT-profile key on disk, exchanged for a token.
		client, err = NewClientWithInterceptorFromKeyFile(ctx, clientInfo.Issuer, clientInfo.KeyPath, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, core)
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	case len(clientInfo.Data) > 0:
		// jwt_profile_json: the same key provided inline as JSON.
		client, err = NewClientWithInterceptorFromKeyFileData(ctx, clientInfo.Issuer, clientInfo.Data, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, core)
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	default:
		// Unreachable in practice (GetClientInfo rejects an unauthenticated
		// provider), but kept as a guard with an actionable message.
		return diag.Errorf("no authentication method available for asset upload; configure one of 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api'")
	}

	resp, err := client.Do(r)
	if err != nil {
		return diag.Errorf("failed to do asset request: %v", err)
	}
	// Always drain the body before closing so the underlying connection can be
	// reused by keep-alive, then close it.
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	// ZITADEL signals upload failures with a non-200 status; surface the status
	// and response body so the cause (e.g. unsupported file type) is visible. Cap
	// the read so a large or hostile body cannot blow up memory or diagnostics.
	if resp.StatusCode != http.StatusOK {
		const maxErrBody = 4 << 10 // 4 KiB is plenty for an API error message.
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrBody))
		if readErr != nil {
			return diag.Errorf("asset request returned %s (failed to read response body: %v)", resp.Status, readErr)
		}
		msg := strings.TrimSpace(string(body))
		if len(body) == maxErrBody {
			msg += " (truncated)"
		}
		return diag.Errorf("asset request returned %s: %s", resp.Status, msg)
	}
	return nil
}

// assetTransport returns the base RoundTripper for asset-upload clients. When the
// provider sets insecure_skip_verify_tls it clones the default transport and
// disables certificate verification, matching the gRPC client's behavior; an
// untouched http.DefaultTransport is returned otherwise.
func assetTransport(insecureSkipVerifyTLS bool) http.RoundTripper {
	if !insecureSkipVerifyTLS {
		return http.DefaultTransport
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{}
	}
	t.TLSClientConfig.InsecureSkipVerify = true
	return t
}

// Interceptor is an http.RoundTripper that stamps a bearer token onto every
// outgoing request, the HTTP equivalent of the gRPC auth interceptor.
type Interceptor struct {
	tokenSource oauth2.TokenSource
	core        http.RoundTripper
}

// NewClientWithInterceptor builds a client for an already-resolved token source
// (access_token, jwt_file, system_api). The source is wrapped in
// oauth2.ReuseTokenSource once here so tokens are cached across requests rather
// than re-fetched each round trip. core is the base transport (see assetTransport).
func NewClientWithInterceptor(tokenSource oauth2.TokenSource, core http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: Interceptor{core: core, tokenSource: oauth2.ReuseTokenSource(nil, tokenSource)},
	}
}

// NewClientWithInterceptorFromKeyFile builds a client that authenticates with a
// JWT-profile key read from keyPath, exchanging it for access tokens lazily. core
// is the base transport (see assetTransport).
func NewClientWithInterceptorFromKeyFile(ctx context.Context, issuer, keyPath string, scopes []string, core http.RoundTripper) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFile(ctx, issuer, keyPath, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: core, tokenSource: oauth2.ReuseTokenSource(nil, ts)},
	}, nil
}

// NewClientWithInterceptorFromKeyFileData is like NewClientWithInterceptorFromKeyFile
// but takes the JWT-profile key as in-memory JSON instead of a file path.
func NewClientWithInterceptorFromKeyFileData(ctx context.Context, issuer string, data []byte, scopes []string, core http.RoundTripper) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFileData(ctx, issuer, data, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: core, tokenSource: oauth2.ReuseTokenSource(nil, ts)},
	}, nil
}

// RoundTrip fetches a (cached) token and sets it as the Authorization header
// before delegating to the underlying transport.
func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	// tokenSource is wrapped in ReuseTokenSource at construction, so this returns
	// the cached token until it expires; do not re-wrap here.
	token, err := i.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	r.Header.Set("authorization", token.TokenType+" "+token.AccessToken)
	// Delegate to the underlying transport, which owns closing the request body
	// per the http.RoundTripper contract; closing it here would be redundant and
	// would panic if a future caller sent a bodyless request.
	return i.core.RoundTrip(r)
}
