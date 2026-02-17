resource "zitadel_action_target" "default" {
  name               = "test-target-name"
  endpoint           = "https://example.com/target"
  target_type        = "REST_WEBHOOK"
  timeout            = "15s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
