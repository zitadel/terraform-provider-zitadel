resource "zitadel_action" "default" {
  org_id          = data.zitadel_org.default.id
  name            = "actionname"
  script          = "testscript"
  timeout         = "10s"
  allowed_to_fail = true
}
