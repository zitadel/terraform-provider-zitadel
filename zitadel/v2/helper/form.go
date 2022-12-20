package helper

import (
	"bytes"
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
	"github.com/zitadel/oidc/pkg/client/profile"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	"golang.org/x/oauth2"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func FormFilePost(clientInfo *ClientInfo, endpoint, path string) diag.Diagnostics {
	file, err := os.Open(path)
	if err != nil {
		return diag.Errorf("failed to read file: %v", err)
	}
	info, err := file.Stat()
	if err != nil {
		return diag.Errorf("failed to read file info: %v", err)
	}
	if info.Size() > (1 << 19) {
		return diag.Errorf("file to large")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	data, err := io.ReadAll(file)
	if err != nil {
		return diag.Errorf("failed to read asset: %v", err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes("file"), escapeQuotes(filepath.Base(file.Name()))))
	h.Set("Content-Type", mimetype.Detect(data).String())
	part, err := writer.CreatePart(h)
	if err != nil {
		return diag.Errorf("failed to create asset part: %v", err)
	}
	io.Copy(part, bytes.NewBuffer(data))
	writer.Close()

	r, err := http.NewRequest(http.MethodPost, clientInfo.Issuer+endpoint, body)
	if err != nil {
		return diag.Errorf("failed to create asset request: %v", err)
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	client, err := NewClientWithInterceptor(clientInfo.Issuer, clientInfo.KeyPath, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
	if err != nil {
		return diag.Errorf("failed to create client: %v", err)
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

func NewClientWithInterceptor(issuer, keyPath string, scopes []string) (*http.Client, error) {
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFile(issuer, keyPath, scopes)
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
