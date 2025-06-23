terraform {
  required_providers {
    zitadel = {
      source = "zitadel/zitadel"
    }
  }
}

provider "zitadel" {
  domain = "localhost"
  insecure = true
  port = "8080"
}
