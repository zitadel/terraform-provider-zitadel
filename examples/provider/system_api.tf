provider "zitadel" {
  domain   = "org-level-tests.default.127.0.0.1.sslip.io"
  insecure = true
  port     = "8080"

  system_api {
    user = "system-api-sa"
    key  = file("keys/system-api-sa.pem")
  }
}

# Alternative split-key configuration:
# provider "zitadel" {
#   domain   = "org-level-tests.default.127.0.0.1.sslip.io"
#   insecure = true
#   port     = "8080"
# 
#   system_api {
#     user        = "system-api-sa"
#     private_key = file("keys/system-api-sa-private.pem")
#     public_key  = file("keys/system-api-sa-public.pem")
#   }
# }
