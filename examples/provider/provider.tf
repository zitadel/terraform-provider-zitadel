terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.0-alpha.1"
    }
  }
}

provider zitadel {
  domain = "localhost:8080"
  insecure = "true"
  project = "170832731415117995"
  token   = "local-token"
}