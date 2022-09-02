terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.0-alpha.7"
    }
  }
}

provider zitadel {
  domain = "localhost"
  insecure = "true"
  port = "8080"
  project = "170832731415117995"
  token   = "local-token"
}