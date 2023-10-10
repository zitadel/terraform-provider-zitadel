package gen

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetEmptySMSAttrs(plan types.Object) {
	plan.Attributes()["title"] = types.StringValue("")
	plan.Attributes()["pre_header"] = types.StringValue("")
	plan.Attributes()["subject"] = types.StringValue("")
	plan.Attributes()["greeting"] = types.StringValue("")
	plan.Attributes()["button_text"] = types.StringValue("")
	plan.Attributes()["footer_text"] = types.StringValue("")
}

func DeleteSMSAttributes(s tfsdk.Schema) {
	//only sms message
	delete(s.Attributes, "title")
	delete(s.Attributes, "pre_header")
	delete(s.Attributes, "subject")
	delete(s.Attributes, "greeting")
	delete(s.Attributes, "button_text")
	delete(s.Attributes, "footer_text")
}
