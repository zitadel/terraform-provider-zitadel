# The resource can be imported using the ID format `<id:user_id[:org_id][:key_details][:public_key]>`, e.g.
# When importing with a public key, make sure to base64 encode it
# terraform import zitadel_machine_key.imported '123456789012345678:123456789012345678:123456789012345678::Ii0tLS0tQkVHSU4gUF...

# The key details contain :, which you can escape with __SEMICOLON__, e.g.
terraform import zitadel_machine_key.imported "123456789012345678:123456789012345678:123456789012345678:$(cat ~/Downloads/123456789012345678.json | sed -e 's/:/__SEMICOLON__/g')"
