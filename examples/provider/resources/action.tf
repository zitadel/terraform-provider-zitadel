resource "zitadel_action" "action" {
  org_id          = data.zitadel_org.org.id
  name            = "actionname"
  script          = "testscript"
  timeout         = "10s"
  allowed_to_fail = true
}
