package hosted_login_translation

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Organization-level translations for the **hosted login v2** (`/ui/v2/login`) in one language, using the settings/v2 API. **Requires ZITADEL 4.x.** " +
			"ZITADEL checks the instance-level `iam.policy.read` and `iam.policy.write` permissions for these translations, so the service account needs an instance role such as `IAM_OWNER`; `ORG_OWNER` is not sufficient. " +
			"Keys that are not set fall back to the instance-level translations and then to the built-in translations of the hosted login. " +
			"ZITADEL has no API to remove translations, so destroying this resource only removes it from the Terraform state and leaves the last applied translations in place. " +
			"Instance-level defaults are managed by `zitadel_default_hosted_login_translation`. " +
			"The legacy login UI (v1) is customized with `zitadel_login_texts`.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDComputedResourceField,
			LanguageVar: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: helper.NonEmptyString(LanguageVar),
				Description:      "BCP 47 language tag of the translations, e.g. `en`, `de` or `fr-CH`",
			},
			translationsVar: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: helper.NonEmptyJSONObject(translationsVar),
				DiffSuppressFunc: structure.SuppressJsonDiff,
				Description:      "Translations as a JSON object with at least one key, for example built with `jsonencode`. The keys follow the structure of the [hosted login locale files](https://github.com/zitadel/zitadel/tree/main/apps/login/locales), e.g. `loginname.title`. Replaces all translations previously set for this language on the organization.",
			},
		},
		ReadContext:   read,
		CreateContext: update,
		DeleteContext: delete,
		UpdateContext: update,
		Importer: helper.ImportWithEmptyID(
			helper.NewImportAttribute(LanguageVar, helper.ConvertNonEmpty, false),
			helper.ImportOptionalOrgAttribute,
		),
	}
}
