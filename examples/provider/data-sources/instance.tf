data "zitadel_instance" "default" {
  instance_id = "123456789012345678"
}

output "instance" {
  value = data.zitadel_instance.default
}
