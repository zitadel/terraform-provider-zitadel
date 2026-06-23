// Package ephtest provides a minimal acceptance-test harness for the ephemeral
// secret resources.
//
// It deliberately does NOT depend on zitadel/helper/test_utils. That package
// pulls in github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource, whose
// init() registers a global `-sweep` flag. The ephemeral tests must use the
// newer github.com/hashicorp/terraform-plugin-testing/helper/resource (only it
// understands ephemeral resources and TerraformVersionChecks), and that package
// registers the SAME `-sweep` flag. Importing both into one test binary panics
// with "flag redefined: sweep". Keeping this harness free of the SDKv2 testing
// package lets the ephemeral tests link cleanly.
package ephtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	mgmt "github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	insecure = true
	port     = "8080"
)

// OrgFrame bundles everything an ephemeral resource test needs: the provider
// HCL snippet, a data "zitadel_org" "default" dependency, the resolved default
// org id, a unique suffix for resource names, a management client for remote
// assertions, and the proto-v6 provider factories (zitadel + echo).
type OrgFrame struct {
	Context           context.Context
	ProviderSnippet   string
	OrgDependency     string
	OrgID             string
	UniqueID          string
	Client            *mgmt.Client
	ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
}

// NewOrgFrame configures a provider against the org-level acceptance instance
// and resolves the default organization. Use it for ephemeral resources scoped
// to the default org (client secrets, machine secrets, keys, PATs).
func NewOrgFrame(t *testing.T) *OrgFrame {
	t.Helper()
	return newFrame(t, acceptance.GetConfig().OrgLevel, true)
}

// NewInstanceFrame configures a provider against the instance-level acceptance
// instance and does NOT resolve a default org. Use it for tests that create a
// fresh organization and need org-domain validation to actually require a
// challenge: the instance-level instance's domain policy keeps added domains
// unverified, whereas the org-level instance auto-verifies them.
func NewInstanceFrame(t *testing.T) *OrgFrame {
	t.Helper()
	return newFrame(t, acceptance.GetConfig().InstanceLevel, false)
}

// newFrame builds a frame for the given isolated instance. When resolveOrg is
// true it looks up the default organization and exposes a data "zitadel_org"
// "default" dependency; otherwise OrgID/OrgDependency are left empty.
func newFrame(t *testing.T, cfg acceptance.IsolatedInstance, resolveOrg bool) *OrgFrame {
	t.Helper()
	ctx := context.Background()

	// Configure the SDKv2 provider once so we obtain a *helper.ClientInfo (its
	// Meta) and can build a management client for remote assertions.
	sdkProvider := zitadel.Provider()
	if d := sdkProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":           cfg.Domain,
		"insecure":         insecure,
		"port":             port,
		"jwt_profile_json": string(cfg.AdminSAJSON),
	})); d.HasError() {
		t.Fatalf("failed to configure test provider: %v", d)
	}
	clientInfo := sdkProvider.Meta().(*helper.ClientInfo)

	client, err := helper.GetManagementClient(ctx, clientInfo)
	if err != nil {
		t.Fatalf("failed to build management client: %v", err)
	}

	var orgID, orgDependency string
	if resolveOrg {
		org, oerr := client.GetOrgByDomainGlobal(ctx, &management.GetOrgByDomainGlobalRequest{Domain: "zitadel." + cfg.Domain})
		if oerr != nil {
			t.Fatalf("failed to look up default org: %v", oerr)
		}
		orgID = org.GetOrg().GetId()
		orgDependency = fmt.Sprintf(`
data "zitadel_org" "default" {
  id = "%s"
}
`, orgID)
	}

	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
  domain           = "%s"
  insecure         = %t
  port             = "%s"
  jwt_profile_json = <<KEY
%s
KEY
}
`, cfg.Domain, insecure, port, string(cfg.AdminSAJSON))

	// The mux factory serves both halves of the provider: the plugin-framework
	// provider (which hosts the ephemeral resources under test) and the SDKv2
	// provider upgraded to protocol v6 (which hosts the managed dependency
	// resources such as zitadel_project / zitadel_application_v2).
	factory := func() (tfprotov6.ProviderServer, error) {
		// Upgrade the SDKv2 provider to protocol v6 up front so a failure
		// surfaces as a real error instead of being swallowed into a nil
		// provider inside the mux (which would later nil-dereference).
		upgraded, err := tf5to6server.UpgradeServer(ctx, func() tfprotov5.ProviderServer {
			return zitadel.Provider().GRPCProvider()
		})
		if err != nil {
			return nil, err
		}
		muxServer, err := tf6muxserver.NewMuxServer(ctx,
			providerserver.NewProtocol6(zitadel.NewProviderPV6()),
			func() tfprotov6.ProviderServer { return upgraded },
		)
		if err != nil {
			return nil, err
		}
		return muxServer.ProviderServer(), nil
	}

	return &OrgFrame{
		Context:         ctx,
		ProviderSnippet: providerSnippet,
		OrgDependency:   orgDependency,
		OrgID:           orgID,
		UniqueID:        acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
		Client:          client,
		ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"zitadel": factory,
			"echo":    echoprovider.NewProviderServer(),
		},
	}
}
