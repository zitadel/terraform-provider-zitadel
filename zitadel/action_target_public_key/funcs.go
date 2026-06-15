package action_target_public_key

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
	filterv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// configBool reports the practitioner-supplied value of a bool attribute from the
// raw configuration, returning false when the attribute is absent or unknown.
// Unlike d.Get / d.GetOkExists it does NOT fall back to state, which makes it
// the right primitive for distinguishing "the user wrote `active = false`" from
// "the user did not write `active` at all" on Optional+Computed attributes.
func configBool(d *schema.ResourceData, attr string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return false
	}
	v := raw.GetAttr(attr)
	if v.IsNull() || !v.IsKnown() {
		return false
	}
	return v.True()
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	targetID := d.Get(targetIDVar).(string)
	keyID := d.Id()

	// ZITADEL requires deactivating a key before it can be removed.
	_, err = client.DeactivatePublicKey(ctx, &actionv2.DeactivatePublicKeyRequest{
		TargetId: targetID,
		KeyId:    keyID,
	})
	if err != nil && helper.IgnorePreconditionError(err) != nil && helper.IgnoreIfNotFoundError(err) != nil {
		return diag.Errorf("failed to deactivate public key before removal: %v", err)
	}

	_, err = client.RemovePublicKey(ctx, &actionv2.RemovePublicKeyRequest{
		TargetId: targetID,
		KeyId:    keyID,
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) != nil {
		return diag.Errorf("failed to delete public key: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &actionv2.AddPublicKeyRequest{
		TargetId:  d.Get(targetIDVar).(string),
		PublicKey: []byte(d.Get(publicKeyVar).(string)),
	}

	if v, ok := d.GetOk(expirationDateVar); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("failed to parse expiration_date: %v", err)
		}
		req.ExpirationDate = timestamppb.New(t)
	}

	resp, err := client.AddPublicKey(ctx, req)
	if err != nil {
		return diag.Errorf("failed to add public key: %v", err)
	}
	d.SetId(resp.GetKeyId())

	if err := d.Set(keyIDVar, resp.GetKeyId()); err != nil {
		return diag.Errorf("failed to set key_id: %v", err)
	}

	// Only call ActivatePublicKey when the user explicitly opted in via active=true.
	// An unset active leaves the key in ZITADEL's default (inactive) state to preserve
	// pre-existing behavior for configs that don't set the field.
	// d.SetId above ensures the key is tracked in state even if activation fails;
	// the next apply will retry via the Update path rather than orphan or duplicate it.
	// FailedPrecondition (key already active) is treated as idempotent success.
	//
	// Use the raw config rather than d.Get/d.GetOkExists so we can distinguish
	// "absent from config" from "present and false" — d.Get-based readers blend
	// state and config, which is ambiguous for Optional+Computed attributes.
	wantActive := configBool(d, activeVar)

	// Persist active=false to state BEFORE attempting activation so that, if Activate
	// returns a non-precondition error, the partial state we leave behind matches the
	// server (key exists, inactive) and the next plan correctly shows false -> true.
	// This is important when users run with -refresh=false, which would otherwise skip
	// the read() that reconciles state.
	if err := d.Set(activeVar, false); err != nil {
		return diag.Errorf("failed to set active: %v", err)
	}

	if wantActive {
		if _, err := client.ActivatePublicKey(ctx, &actionv2.ActivatePublicKeyRequest{
			TargetId: req.TargetId,
			KeyId:    resp.GetKeyId(),
		}); err != nil && helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to activate public key: %v", err)
		}
		if err := d.Set(activeVar, true); err != nil {
			return diag.Errorf("failed to set active: %v", err)
		}
	}

	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	if !d.HasChange(activeVar) {
		return nil
	}

	// Only act when the user has an explicit active value in *config*. If the field
	// has been removed we preserve the current remote state rather than implicitly
	// deactivating a key the user no longer manages. We inspect the raw config (not
	// d.Get/d.GetOkExists, which blend state and config) so this stays correct even
	// when state has a value populated by Create or Read.
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}
	configActive := rawConfig.GetAttr(activeVar)
	if configActive.IsNull() || !configActive.IsKnown() {
		return nil
	}
	wantActive := configActive.True()

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	targetID := d.Get(targetIDVar).(string)
	keyID := d.Id()

	if wantActive {
		if _, err := client.ActivatePublicKey(ctx, &actionv2.ActivatePublicKeyRequest{
			TargetId: targetID,
			KeyId:    keyID,
		}); err != nil && helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to activate public key: %v", err)
		}
	} else {
		if _, err := client.DeactivatePublicKey(ctx, &actionv2.DeactivatePublicKeyRequest{
			TargetId: targetID,
			KeyId:    keyID,
		}); err != nil && helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to deactivate public key: %v", err)
		}
	}

	if err := d.Set(activeVar, wantActive); err != nil {
		return diag.Errorf("failed to set active: %v", err)
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	targetID := d.Get(targetIDVar).(string)
	keyID := helper.GetID(d, keyIDVar)

	resp, err := client.ListPublicKeys(ctx, &actionv2.ListPublicKeysRequest{
		TargetId: targetID,
		Filters: []*actionv2.PublicKeySearchFilter{
			{
				Filter: &actionv2.PublicKeySearchFilter_KeyIdsFilter{
					KeyIdsFilter: &filterv2.InIDsFilter{
						Ids: []string{keyID},
					},
				},
			},
		},
	})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to list public keys: %v", err)
	}

	keys := resp.GetPublicKeys()
	if len(keys) == 0 {
		d.SetId("")
		return nil
	}

	key := keys[0]

	// Only overwrite public_key in state if it is currently unset (e.g., during import).
	// PEM formatting may differ between the configured value and the API response,
	// which would cause perpetual diffs and forced recreation.
	set := map[string]interface{}{
		targetIDVar:    targetID,
		keyIDVar:       key.GetKeyId(),
		activeVar:      key.GetActive(),
		fingerprintVar: key.GetFingerprint(),
	}

	if currentPublicKey, _ := d.Get(publicKeyVar).(string); currentPublicKey == "" {
		set[publicKeyVar] = string(key.GetPublicKey())
	}

	if key.GetExpirationDate() != nil && key.GetExpirationDate().IsValid() {
		set[expirationDateVar] = key.GetExpirationDate().AsTime().Format(time.RFC3339)
	} else {
		set[expirationDateVar] = ""
	}
	if key.GetCreationDate() != nil && key.GetCreationDate().IsValid() {
		set[creationDateVar] = key.GetCreationDate().AsTime().Format(time.RFC3339)
	} else {
		set[creationDateVar] = ""
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of public key: %v", k, err)
		}
	}
	d.SetId(key.GetKeyId())
	return nil
}
