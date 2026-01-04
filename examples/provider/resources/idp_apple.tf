resource "zitadel_idp_apple" "default" {
  name                 = "Apple"
  client_id            = "com.example.app"
  team_id              = "ABCDE12345"
  key_id               = "FGHIJ67890"
  private_key          = "*****abc123"
  scopes               = ["name", "email"]
  is_linking_allowed   = false
  is_creation_allowed  = true
  is_auto_creation     = false
  is_auto_update       = true
  auto_linking         = "AUTO_LINKING_OPTION_USERNAME"
}
