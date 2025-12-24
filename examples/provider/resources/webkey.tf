resource "zitadel_webkey" "default" {
	org_id = data.zitadel_org.default.id
	rsa {
		bits = "RSA_BITS_2048"
	}
}
