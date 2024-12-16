resource "zitadel_idp_ldap" "default" {
  name                 = "LDAP"
  servers              = ["ldaps://my.primary.server:389", "ldaps://my.secondary.server:389"]
  start_tls            = false
  base_dn              = "dc=example,dc=com"
  bind_dn              = "cn=admin,dc=example,dc=com"
  bind_password        = "Password1!"
  user_base            = "dn"
  user_object_classes  = ["inetOrgPerson"]
  user_filters         = ["uid", "email"]
  timeout              = "10s"
  id_attribute         = "uid"
  first_name_attribute = "firstname"
  last_name_attribute  = "lastname"
  is_linking_allowed   = false
  is_creation_allowed  = true
  is_auto_creation     = false
  is_auto_update       = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}


