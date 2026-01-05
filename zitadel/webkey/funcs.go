package webkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &webkey.CreateWebKeyRequest{}

	if rsa, ok := d.Get(RSABlock).([]interface{}); ok && len(rsa) > 0 && rsa[0] != nil {
		attrs := rsa[0].(map[string]interface{})
		req.Key = &webkey.CreateWebKeyRequest_Rsa{
			Rsa: &webkey.RSA{
				Bits:   webkey.RSABits(webkey.RSABits_value[attrs[BitsVar].(string)]),
				Hasher: webkey.RSAHasher(webkey.RSAHasher_value[attrs[HasherVar].(string)]),
			},
		}
	} else if ecdsa, ok := d.Get(ECDSABlock).([]interface{}); ok && len(ecdsa) > 0 && ecdsa[0] != nil {
		attrs := ecdsa[0].(map[string]interface{})
		req.Key = &webkey.CreateWebKeyRequest_Ecdsa{
			Ecdsa: &webkey.ECDSA{
				Curve: webkey.ECDSACurve(webkey.ECDSACurve_value[attrs[CurveVar].(string)]),
			},
		}
	} else if ed25519, ok := d.Get(ED25519Block).([]interface{}); ok && len(ed25519) > 0 {
		req.Key = &webkey.CreateWebKeyRequest_Ed25519{
			Ed25519: &webkey.ED25519{},
		}
	}

	if req.Key == nil {
		return diag.Errorf("one of rsa, ecdsa, or ed25519 blocks must be specified")
	}

	resp, err := client.CreateWebKey(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.GetId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	webkeyID := helper.GetID(d, WebKeyIDVar)

	resp, err := client.ListWebKeys(helper.CtxWithOrgID(ctx, d), &webkey.ListWebKeysRequest{})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	for _, key := range resp.GetWebKeys() {
		if key.Id == webkeyID {
			set := map[string]interface{}{
				StateVar: key.GetState().String(),
			}

			var keyType string
			if key.GetRsa() != nil {
				keyType = "RSA"
			} else if key.GetEcdsa() != nil {
				keyType = "ECDSA"
			} else if key.GetEd25519() != nil {
				keyType = "ED25519"
			}
			set[KeyTypeVar] = keyType

			for k, v := range set {
				if err := d.Set(k, v); err != nil {
					return diag.Errorf("failed to set %s of webkey: %v", k, err)
				}
			}
			d.SetId(key.Id)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetWebKeyClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListWebKeys(helper.CtxWithOrgID(ctx, d), &webkey.ListWebKeysRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, key := range resp.GetWebKeys() {
		if key.Id == d.Id() && key.State == webkey.State_STATE_ACTIVE {
			return diag.Errorf("cannot delete an active webkey. Please activate a different key first to deactivate this one")
		}
	}

	_, err = client.DeleteWebKey(helper.CtxWithOrgID(ctx, d), &webkey.DeleteWebKeyRequest{Id: d.Id()})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}
