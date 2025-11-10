resource "zitadel_action_execution" "default" {
  execution_type = "events"

  event = "user.human.added"
  group = "user"

  targets = [
    data.zitadel_action_target.default.id
  ]
}
