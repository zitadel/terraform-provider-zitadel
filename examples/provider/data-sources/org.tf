data "zitadel_org" "org" {
  id = "177073608051458051"
}

output "org" {
  value = data.zitadel_org.org
}
