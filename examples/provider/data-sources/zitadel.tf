data "zitadel" "default" {
  //
}

output "token" {
  value     = data.zitadel.default.token
  sensitive = true
}
