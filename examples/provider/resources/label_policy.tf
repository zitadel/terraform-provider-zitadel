resource zitadel_label_policy label_policy {
  depends_on = [zitadel_org.org]

  org_id                 = zitadel_org.org.id
  primary_color          = "#5469d4"
  hide_login_name_suffix = true
  warn_color             = "#cd3d56"
  background_color       = "#fafafa"
  font_color             = "#000000"
  primary_color_dark     = "#a5b4fc"
  background_color_dark  = "#111827"
  warn_color_dark        = "#ff3b5b"
  font_color_dark        = "#ffffff"
  disable_watermark      = false
  set_active             = true
  logo_hash              = filemd5("/path/to/logo.jpg")
  logo_path              = "/path/to/logo.jpg"
  logo_dark_hash         = filemd5("/path/to/logo_dark.jpg")
  logo_dark_path         = "/path/to/logo_dark.jpg"
  icon_hash              = filemd5("/path/to/icon.jpg")
  icon_path              = "/path/to/icon.jpg"
  icon_dark_hash         = filemd5("/path/to/icon_dark.jpg")
  icon_dark_path         = "/path/to/icon_dark.jpg"
  font_hash              = filemd5("/path/to/font.tff")
  font_path              = "/path/to/font.tff"
}