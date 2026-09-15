resource "zitadel_default_hosted_login_translation" "default" {
  language = "en"
  translations = jsonencode({
    loginname = {
      title       = "Welcome to ACME"
      description = "Sign in with your ACME account."
    }
  })
}
