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
	if v, ok := d.GetOkExists(activeVar); ok && v.(bool) {
		if _, err := client.ActivatePublicKey(ctx, &actionv2.ActivatePublicKeyRequest{
			TargetId: req.TargetId,
			KeyId:    resp.GetKeyId(),
		}); err != nil {
			return diag.Errorf("failed to activate public key: %v", err)
		}
	}

	// AddPublicKey only returns the key ID, so call read() to populate the remaining
	// computed attributes (active, fingerprint, creation_date) from the server.
	return read(ctx, d, m)
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	if !d.HasChange(activeVar) {
		return nil
	}

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

	if d.Get(activeVar).(bool) {
		if _, err := client.ActivatePublicKey(ctx, &actionv2.ActivatePublicKeyRequest{
			TargetId: targetID,
			KeyId:    keyID,
		}); err != nil {
			return diag.Errorf("failed to activate public key: %v", err)
		}
		return read(ctx, d, m)
	}

	if _, err := client.DeactivatePublicKey(ctx, &actionv2.DeactivatePublicKeyRequest{
		TargetId: targetID,
		KeyId:    keyID,
	}); err != nil && helper.IgnorePreconditionError(err) != nil {
		return diag.Errorf("failed to deactivate public key: %v", err)
	}
	return read(ctx, d, m)
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
