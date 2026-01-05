resource "zitadel_organization_metadata" "default" {
  organization_id = zitadel_organization.default.id
  key             = "example_key"
  value           = "example_value"
}

resource "zitadel_organization_metadata" "binary" {
  organization_id = zitadel_organization.default.id
  key             = "binary_data"
  value = base64encode("binary data here")
}
