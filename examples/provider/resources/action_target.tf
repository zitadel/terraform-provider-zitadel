resource "zitadel_action_target" "default" {
  name               = "test-target-name"
  endpoint           = "https://httpstat.us/200"
  target_type        = "REST_WEBHOOK"
  timeout            = "15s"
  interrupt_on_error = true
  payload_type       = "PAYLOAD_TYPE_JSON"
}
