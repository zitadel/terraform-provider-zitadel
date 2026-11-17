package zitadel_test

import (
	"context"
	"net"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/oidctest"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestSessionTokenDatasourceOIDCRequests(t *testing.T) {
	const headerName, headerValue = "X-Forwarded-Host", "zitadel.example.com"
	for _, useData := range []bool{false, true} {
		name := "key file"
		if useData {
			name = "key data"
		}
		t.Run(name, func(t *testing.T) {
			idp := oidctest.New(t, true)
			u, err := url.Parse(idp.URL)
			if err != nil {
				t.Fatal(err)
			}
			host, port, err := net.SplitHostPort(u.Host)
			if err != nil {
				t.Fatal(err)
			}
			key := oidctest.KeyFileJSON(t)
			keyFile, keyJSON := "", ""
			if useData {
				keyJSON = string(key)
			} else {
				keyFile = t.TempDir() + "/key.json"
				if err := os.WriteFile(keyFile, key, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			info, err := helper.GetClientInfo(context.Background(), false, host, "", "", "", keyFile, keyJSON, "", "", "", "", "", "", port, true, map[string]string{headerName: headerValue})
			if err != nil {
				t.Fatal(err)
			}

			ds := zitadel.GetZitadelDatasource()
			d := ds.TestResourceData()
			if diags := ds.ReadContext(context.Background(), d, info); diags.HasError() {
				t.Fatalf("ReadContext() = %v", diags)
			}
			if got := d.Get("token").(string); got != "token" {
				t.Errorf("token = %q, want %q", got, "token")
			}
			for _, path := range []string{oidctest.DiscoveryPath, oidctest.TokenPath} {
				if got := idp.Header(path, headerName); got != headerValue {
					t.Errorf("%s header %s = %q, want %q", path, headerName, got, headerValue)
				}
			}
		})
	}
}

// TestSessionTokenDatasourceHostHeader checks the instance that is reached on an internal
// address with its external host as Host header.
func TestSessionTokenDatasourceHostHeader(t *testing.T) {
	const external = "zitadel.example.com"
	idp := oidctest.New(t, false)
	u, err := url.Parse(idp.URL)
	if err != nil {
		t.Fatal(err)
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatal(err)
	}
	info, err := helper.GetClientInfo(context.Background(), true, host, "", "", "", "", string(oidctest.KeyFileJSON(t)), "", "", "", "", "", "", port, false, map[string]string{"Host": external})
	if err != nil {
		t.Fatal(err)
	}

	ds := zitadel.GetZitadelDatasource()
	d := ds.TestResourceData()
	if diags := ds.ReadContext(context.Background(), d, info); diags.HasError() {
		t.Fatalf("ReadContext() = %v", diags)
	}

	for _, path := range []string{oidctest.DiscoveryPath, oidctest.TokenPath} {
		if got := idp.Host(path); got != external {
			t.Errorf("%s host = %q, want %q", path, got, external)
		}
	}
}

func TestAccSessionTokenDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel")
	config := `data "zitadel_zitadel" "token" {}`
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		resource.TestCheckResourceAttrSet("data.zitadel_zitadel.token", "token"),
		nil,
	)
}
