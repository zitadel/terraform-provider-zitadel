terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.2.0"
    }
  }
}

provider "zitadel" {
  domain           = "localhost"
  insecure         = "true"
  port             = "8080"
  jwt_profile_file = "local-token"
}
