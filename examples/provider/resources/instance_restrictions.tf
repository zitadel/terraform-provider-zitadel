resource "zitadel_instance_restrictions" "default" {
  disallow_public_org_registration = true
  allowed_languages                = ["en", "de"]
}
