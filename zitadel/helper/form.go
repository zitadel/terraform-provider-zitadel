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
	for k, v := range additionalHeaders {
		r.Header.Add(k, v)
	}

	if clientInfo.KeyPath != "" {
		client, err = NewClientWithInterceptorFromKeyFile(ctx, clientInfo.Issuer, clientInfo.KeyPath, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	} else if len(clientInfo.Data) > 0 {
		client, err = NewClientWithInterceptorFromKeyFileData(ctx, clientInfo.Issuer, clientInfo.Data, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
		if err != nil {
			return diag.Errorf("failed to create client: %v", err)
		}
	} else {
		return diag.Errorf("either 'jwt_profile_file' or 'jwt_profile_json' is required")
	}

	resp, err := client.Do(r)
	if err != nil || resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to do asset request: %v", err)
	}
	return nil
}

type Interceptor struct {
	tokenSource oauth2.TokenSource
	core        http.RoundTripper
}

func NewClientWithInterceptorFromKeyFile(ctx context.Context, issuer, keyPath string, scopes []string) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFile(ctx, issuer, keyPath, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: http.DefaultTransport, tokenSource: ts},
	}, nil
}

func NewClientWithInterceptorFromKeyFileData(ctx context.Context, issuer string, data []byte, scopes []string) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFileData(ctx, issuer, data, scopes)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: Interceptor{core: http.DefaultTransport, tokenSource: ts},
	}, nil
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		_ = r.Body.Close()
	}()

	ts := oauth2.ReuseTokenSource(nil, i.tokenSource)

	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	r.Header.Set("authorization", token.TokenType+" "+token.AccessToken)
	return i.core.RoundTrip(r)
}
