# The resource can be imported using the ID format `<id:user_id[:org_id][:key_details][:public_key]>`, e.g.
# When importing with a public key, make sure to base64 encode it
# terraform import zitadel_machine_key.imported '123456789012345678:123456789012345678:123456789012345678::Ii0tLS0tQkVHSU4gUF...

terraform import zitadel_machine_key.imported '123456789012345678:123456789012345678:123456789012345678:{"type":"serviceaccount","keyId":"123456789012345678","key":"-----BEGIN RSA PRIVATE KEY-----\nMIIEpQ...-----END RSA PRIVATE KEY-----\n","userId":"123456789012345678"}'
