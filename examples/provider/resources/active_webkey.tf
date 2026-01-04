resource "zitadel_webkey" "key_v1" {
  org_id = data.zitadel_org.default.id
  rsa {}
}

resource "zitadel_webkey" "key_v2" {
  org_id = data.zitadel_org.default.id
  ecdsa {}
}

resource "zitadel_active_webkey" "default" {
  org_id = data.zitadel_org.default.id
  key_id = zitadel_webkey.key_v1.id
}
