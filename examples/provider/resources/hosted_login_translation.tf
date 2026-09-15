resource "zitadel_hosted_login_translation" "default" {
  org_id   = data.zitadel_org.default.id
  language = "en"
  translations = jsonencode({
    loginname = {
      title       = "Welcome to ACME"
      description = "Sign in with your ACME account."
    }
  })
}
