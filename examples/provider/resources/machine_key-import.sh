# The resource can be imported using the ID format `<id:user_id[:org_id][:key_details]>`, e.g.
terraform import machine_key.imported '123456789012345678:123456789012345678:123456789012345678:{"type":"serviceaccount","keyId":"123456789012345678","key":"-----BEGIN RSA PRIVATE KEY-----\nMIIEpQ...-----END RSA PRIVATE KEY-----\n","userId":"123456789012345678"}'
