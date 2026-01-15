# event:name for specific event, group:name for event group, all for all events
terraform import zitadel_action_execution_event.imported 'event:user.human.added'
terraform import zitadel_action_execution_event.imported 'group:user.human'
terraform import zitadel_action_execution_event.imported 'all'