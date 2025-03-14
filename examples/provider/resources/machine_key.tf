resource "zitadel_machine_key" "default" {
  org_id          = data.zitadel_org.default.id
  user_id         = data.zitadel_machine_user.default.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
  public_key      = <<-EOT
    -----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApj7JHjDLo2TwiJznwMrD
    97ybWoRegSK1rx37+i+Yrmhaee0GuOyj+hWG8/yKazAbZfYB0atO/zHxy1BtFNfX
    uYZS689TvfZVP6TctonH0VTlDDKOjmkGl472DhJvLvwjPXq1e55jS0kToK5lGRW6
    Qrgm7m/KiF96Qmp5kUbF1sThVtKBW9GIAuzWEk3O9opftd/NH3BxvUToWLgG/GFx
    hLeOTrcuPibVHkHbIjt1VHaOD8rKAaRV+KBZUmyS9vdo629wfSx/ylUmwWZ6YUTj
    khnqTi0s7j/oLGJNk+DSjMzkcgls0gzXAwPfiEnjEB+Xxw3LnR6k17HyYxqQs7kz
    ZwIDAQAB
    -----END PUBLIC KEY-----
EOT
}
