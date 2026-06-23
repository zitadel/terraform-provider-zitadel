package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

// EphemeralSecretBase factors out the boilerplate every ephemeral secret
// resource in this provider shares: capturing the configured *ClientInfo that
// the provider hands down through Configure.
//
// The plugin-framework delivers provider-level data to an ephemeral resource
// via provider.ConfigureResponse.EphemeralResourceData (a *different* field
// from the ResourceData/DataSourceData used by managed resources and data
// sources). The provider's Configure sets all three to the same *ClientInfo,
// so embedding this base is all an ephemeral resource needs to obtain an
// authenticated client.
//
// Ephemeral resources are how this provider lets operators mint or rotate a
// generated credential (client secret, machine secret, key, PAT, domain
// validation token) during a single `terraform apply` *without* persisting the
// secret to Terraform state — unlike the equivalent managed resources, which
// expose the generated secret as a Computed attribute and therefore write it
// to the state file. See issue #413.
type EphemeralSecretBase struct {
	// ClientInfo is populated in Configure and read in each resource's Open.
	// It is nil until Configure runs, which the framework guarantees happens
	// before Open during a real apply.
	ClientInfo *ClientInfo
}

// Configure stores the provider-supplied *ClientInfo on the embedding resource.
//
// req.ProviderData is nil during early validation walks (before the provider
// itself is configured); we must tolerate that and simply return, exactly as
// the managed framework resources do.
func (b *EphemeralSecretBase) Configure(_ context.Context, req ephemeral.ConfigureRequest, _ *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	b.ClientInfo = req.ProviderData.(*ClientInfo)
}
