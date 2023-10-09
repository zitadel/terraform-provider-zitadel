package default_label_policy

const (
	PrimaryColorVar        = "primary_color"
	hideLoginNameSuffixVar = "hide_login_name_suffix"
	warnColorVar           = "warn_color"
	backgroundColorVar     = "background_color"
	fontColorVar           = "font_color"
	primaryColorDarkVar    = "primary_color_dark"
	backgroundColorDarkVar = "background_color_dark"
	warnColorDarkVar       = "warn_color_dark"
	fontColorDarkVar       = "font_color_dark"
	disableWatermarkVar    = "disable_watermark"
	LogoPathVar            = "logo_path"
	LogoHashVar            = "logo_hash"
	logoURLVar             = "logo_url"
	IconPathVar            = "icon_path"
	IconHashVar            = "icon_hash"
	iconURLVar             = "icon_url"
	LogoDarkPathVar        = "logo_dark_path"
	LogoDarkHashVar        = "logo_dark_hash"
	logoURLDarkVar         = "logo_url_dark"
	IconDarkPathVar        = "icon_dark_path"
	IconDarkHashVar        = "icon_dark_hash"
	iconURLDarkVar         = "icon_url_dark"
	FontPathVar            = "font_path"
	FontHashVar            = "font_hash"
	fontURLVar             = "font_url"
	SetActiveVar           = "set_active"
)

const (
	assetAPI       = "/assets/v1"
	labelPolicyURL = "/instance/policy/label"
	logoURL        = assetAPI + labelPolicyURL + "/logo"
	logoDarkURL    = logoURL + "/dark"
	iconURL        = assetAPI + labelPolicyURL + "/icon"
	iconDarkURL    = iconURL + "/dark"
	fontURL        = assetAPI + labelPolicyURL + "/font"
)