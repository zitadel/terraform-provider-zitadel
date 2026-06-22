package helper

import (
	"bytes"
	"context"
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

func InstanceFormFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path string) diag.Diagnostics {
	return formFilePost(ctx, clientInfo, endpoint, path, map[string]string{})
}

func OrgFormFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path, orgID string) diag.Diagnostics {
	return formFilePost(ctx, clientInfo, endpoint, path, map[string]string{"x-zitadel-orgid": orgID})
}

func formFilePost(ctx context.Context, clientInfo *ClientInfo, endpoint, path string, additionalHeaders map[string]string) diag.Diagnostics {
	var client *http.Client
	r, err := createMultipartRequest(clientInfo.Issuer, endpoint, path)
	if err != nil {
		return diag.Errorf("failed to create asset request: %v", err)
	}
	// Asset uploads bypass the gRPC transport, so the provider's transport_headers
	// have to be applied to the plain HTTP request here as well. Set them first so
	// the per-request additionalHeaders (e.g. x-zitadel-orgid) always win, and use
	// Set rather than Add so single-valued headers stay deterministic.
	for k, v := range clientInfo.TransportHeaders {
		r.Header.Set(k, v)
	}
	for k, v := range additionalHeaders {
		r.Header.Set(k, v)
	}

	switch {
	case clientInfo.TokenSource != nil:
		// access_token, jwt_file and system_api carry a ready bearer token source.
		client = NewClientWithInterceptor(clientInfo.TokenSource)
	case clientInfo.KeyPath != "":
		client, err = NewClientWithInterceptorFromKeyFile(ctx, clientInfo.Issuer, clientInfo.KeyPath, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	case len(clientInfo.Data) > 0:
		client, err = NewClientWithInterceptorFromKeyFileData(ctx, clientInfo.Issuer, clientInfo.Data, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	default:
		return diag.Errorf("no authentication method available for asset upload; configure one of 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api'")
	}

	resp, err := client.Do(r)
	if err != nil {
		return diag.Errorf("failed to do asset request: %v", err)
	}
	defer func() {
		// Drain so the connection can be reused, then close.
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return diag.Errorf("asset request returned %s (failed to read response body: %v)", resp.Status, readErr)
		}
		return diag.Errorf("asset request returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

type Interceptor struct {
	tokenSource oauth2.TokenSource
	core        http.RoundTripper
}

// NewClientWithInterceptor returns an HTTP client that authenticates requests
// with the given token source, used for the access_token, jwt_file and
// system_api auth modes where the bearer token is already available.
func NewClientWithInterceptor(tokenSource oauth2.TokenSource) *http.Client {
	return &http.Client{
		Transport: Interceptor{core: http.DefaultTransport, tokenSource: oauth2.ReuseTokenSource(nil, tokenSource)},
	}
}

func NewClientWithInterceptorFromKeyFile(ctx context.Context, issuer, keyPath string, scopes []string) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFile(ctx, issuer, keyPath, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: http.DefaultTransport, tokenSource: oauth2.ReuseTokenSource(nil, ts)},
	}, nil
}

func NewClientWithInterceptorFromKeyFileData(ctx context.Context, issuer string, data []byte, scopes []string) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFileData(ctx, issuer, data, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: http.DefaultTransport, tokenSource: oauth2.ReuseTokenSource(nil, ts)},
	}, nil
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		_ = r.Body.Close()
	}()

	// tokenSource is already wrapped in oauth2.ReuseTokenSource at construction,
	// so tokens are cached and reused across requests for this client.
	token, err := i.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	r.Header.Set("authorization", token.TokenType+" "+token.AccessToken)
	return i.core.RoundTrip(r)
}
