# The resource can be imported using the ID format `<id:project_id:app_id[:org_id][:key_details]>`.
# You can use __SEMICOLON__ to escape :, e.g.
terraform import application_key.imported "123456789012345678:123456789012345678:123456789012345678:123456789012345678:$(cat ~/Downloads/123456789012345678.json | sed -e 's/:/__SEMICOLON__/g')"
