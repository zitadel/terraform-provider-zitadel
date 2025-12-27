resource "zitadel_instance_trusted_domain" "default" {
	instance_id = "123456789012345678"  # Optional if in instance context
	domain = "idp.partner.com"
}
