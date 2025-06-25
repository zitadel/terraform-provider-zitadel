package helper_test

import (
	"context"
	"testing"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	settingsclient "github.com/zitadel/zitadel-go/v3/pkg/client/settings/v2"
	settingspb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func TestHelper_DirectConnection(t *testing.T) {
	cfg := acceptance.GetConfig().OrgLevel
	jwtProfileJSON := string(cfg.AdminSAJSON)

	clientInfo, err := helper.GetClientInfo(
		context.Background(),
		true,
		cfg.Domain,
		"",
		"",
		"",
		jwtProfileJSON,
		"8080",
		"",
		"",
	)
	if err != nil {
		t.Fatalf("GetClientInfo() for direct connection failed: %v", err)
	}

	settingsClient, err := settingsclient.NewClient(
		context.Background(),
		clientInfo.Issuer,
		clientInfo.Domain,
		[]string{oidc.ScopeOpenID, client.ScopeZitadelAPI()},
		clientInfo.Options...,
	)
	if err != nil {
		t.Fatalf("settings.NewClient() for direct connection failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = settingsClient.GetGeneralSettings(ctx, &settingspb.GetGeneralSettingsRequest{})
	if err != nil {
		t.Fatalf("GetGeneralSettings() call for direct connection failed: %v", err)
	}
}

func TestHelper_ProxyConnection(t *testing.T) {
	cfg := acceptance.GetConfig().OrgLevel
	jwtProfileJSON := string(cfg.AdminSAJSON)
	proxyURL := "socks5://testuser:testpassword@localhost:1080"

	clientInfo, err := helper.GetClientInfo(
		context.Background(),
		true,
		cfg.Domain,
		"",
		"",
		"",
		jwtProfileJSON,
		"8080",
		proxyURL,
		"",
	)
	if err != nil {
		t.Fatalf("GetClientInfo() for proxy connection failed: %v", err)
	}

	settingsClient, err := settingsclient.NewClient(
		context.Background(),
		clientInfo.Issuer,
		clientInfo.Domain,
		[]string{oidc.ScopeOpenID, client.ScopeZitadelAPI()},
		clientInfo.Options...,
	)
	if err != nil {
		t.Fatalf("settings.NewClient() for proxy connection failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = settingsClient.GetGeneralSettings(ctx, &settingspb.GetGeneralSettingsRequest{})
	if err != nil {
		t.Fatalf("GetGeneralSettings() call through proxy failed: %v", err)
	}
}
