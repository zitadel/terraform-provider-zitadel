
resource zitadel_action action {
  depends_on = [zitadel_org.org]

  org_id          = zitadel_org.org.id
  name            = "actionname"
  script          = "testscript"
  timeout         = "10s"
  allowed_to_fail = true
}