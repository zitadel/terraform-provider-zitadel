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
}