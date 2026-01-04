data "zitadel_instance_custom_domains" "default" {
  # instance_id is optional - uses current context if not provided
}

output "custom_domains" {
  value = data.zitadel_instance_custom_domains.default.domains
}
