package default_hosted_login_translation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default hosted login translations cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSecuritySettingsClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	language := d.Get(LanguageVar).(string)
	if d.IsNewResource() || d.HasChange(translationsVar) {
		translations, err := structure.ExpandJsonFromString(d.Get(translationsVar).(string))
		if err != nil {
			return diag.Errorf("failed to parse translations: %v", err)
		}
		translationsStruct, err := structpb.NewStruct(translations)
		if err != nil {
			return diag.Errorf("failed to convert translations: %v", err)
		}

		_, err = client.SetHostedLoginTranslation(ctx, &settingsv2.SetHostedLoginTranslationRequest{
			Level:        &settingsv2.SetHostedLoginTranslationRequest_Instance{Instance: true},
			Locale:       language,
			Translations: translationsStruct,
		})
		if helper.IsUnimplemented(err) {
			return diag.Errorf("hosted login translations require Zitadel server v4+")
		}
		if err != nil {
			return diag.Errorf("failed to update default hosted login translation: %v", err)
		}
	}

	d.SetId(language)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetSecuritySettingsClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	language := d.Id()
	resp, err := client.GetHostedLoginTranslation(ctx, &settingsv2.GetHostedLoginTranslationRequest{
		Level:             &settingsv2.GetHostedLoginTranslationRequest_Instance{Instance: true},
		Locale:            language,
		IgnoreInheritance: true,
	})
	if err != nil {
		if helper.IsUnimplemented(err) {
			return diag.Errorf("hosted login translations require Zitadel server v4+")
		}
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get default hosted login translation: %v", err)
	}

	translations, err := structure.FlattenJsonToString(resp.GetTranslations().AsMap())
	if err != nil {
		return diag.Errorf("failed to convert translations: %v", err)
	}

	set := map[string]interface{}{
		LanguageVar:     language,
		translationsVar: translations,
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default hosted login translation: %v", k, err)
		}
	}

	return nil
}
