package webkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2beta"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func createWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create webkey")
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
	} else if _, ok := d.Get(ED25519Block).([]interface{}); ok {
		req.Key = &webkey.CreateWebKeyRequest_Ed25519{Ed25519: &webkey.ED25519{}}
	}
	resp, err := client.CreateWebKey(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetId())
	return nil
}

func readWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read webkey")
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
		if key.Id == d.Id() {
			if err := d.Set(StateVar, key.GetState().String()); err != nil {
				return diag.Errorf("failed to set state: %v", err)
			}
			var keyType string
			if key.GetRsa() != nil {
				keyType = "RSA"
			} else if key.GetEcdsa() != nil {
				keyType = "ECDSA"
			} else if key.GetEd25519() != nil {
				keyType = "ED25519"
			}
			if err := d.Set(KeyTypeVar, keyType); err != nil {
				return diag.Errorf("failed to set key_type: %v", err)
			}
			return nil
		}
	}
	d.SetId("")
	return nil
}

func updateWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("zitadel_webkey resource does not support in-place updates")
}

func deleteWebKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete webkey")
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
		return diag.FromErr(err)
	}
	return nil
}
