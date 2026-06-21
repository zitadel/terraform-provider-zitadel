package helper_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/zitadel/oidc/v3/pkg/client/profile"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// TestAccPATAssetUpload is a regression test for
// https://github.com/zitadel/terraform-provider-zitadel/issues/411: when the
// provider authenticates with a Personal Access Token (access_token), asset
// uploads (logo/icon/font) must succeed. They previously failed because the
// asset-upload HTTP client only understood the JWT-profile file modes.
func TestAccPATAssetUpload(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("set TF_ACC=1 and run the acceptance container to run this test")
	}

	ctx := context.Background()
	cfg := acceptance.GetConfig().InstanceLevel
	const port = "8080"
	issuer := "http://" + cfg.Domain + ":" + port

	// A PAT is just a bearer token; mint an equivalent one from the instance's
	// admin service-account key so the test needs no extra fixtures.
	ts, err := profile.NewJWTProfileTokenSourceFromKeyFileData(ctx, issuer, cfg.AdminSAJSON, []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()})
	if err != nil {
		t.Fatalf("build token source: %v", err)
	}
	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("fetch bearer token: %v", err)
	}

	// Configure the provider exactly as a PAT user would: access_token only.
	clientInfo, err := helper.GetClientInfo(ctx, true, cfg.Domain, tok.AccessToken, "", "", "", "", "", "", "", "", "", "", port, false, nil)
	if err != nil {
		t.Fatalf("GetClientInfo: %v", err)
	}

	diags := helper.InstanceFormFilePost(ctx, clientInfo, "/assets/v1/instance/policy/label/logo", writePNG(t))
	if diags.HasError() {
		t.Fatalf("asset upload with access_token failed: %v", diags)
	}
}

func writePNG(t *testing.T) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "logo-*.png")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		t.Fatalf("write png: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close png: %v", err)
	}
	return f.Name()
}
