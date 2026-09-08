data "zitadel_idps" "default" {}

data "zitadel_idps" "gitlab" {
  type = "PROVIDER_TYPE_GITLAB"
}

data "zitadel_idps" "by_name" {
  name        = "Corporate"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

output "gitlab_idp_ids" {
  value = data.zitadel_idps.gitlab.idps[*].id
}
