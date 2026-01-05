resource "zitadel_organization" "default" {
  name = "terraform-test"
}

resource "zitadel_organization" "with_admins" {
  name = "terraform-test-with-admins"
  admins = [
    {
      user_id = "123456789012345678"
      roles = ["ORG_OWNER"]
    }
  ]
}

resource "zitadel_organization" "with_custom_id" {
  name   = "terraform-test-custom"
  org_id = "custom-org-id-123"
}
