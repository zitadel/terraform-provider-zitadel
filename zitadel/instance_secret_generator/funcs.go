package instance_secret_generator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "secret generators cannot be deleted")
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	generatorType := d.Get(generatorTypeVar).(string)
	if _, ok := generatorTypeMap[generatorType]; !ok {
		return diag.Errorf("invalid generator_type %q", generatorType)
	}

	d.SetId(generatorType)

	diags := update(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	return read(ctx, d, m)
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	if !d.HasChanges(lengthVar, expiryVar, includeLowerLettersVar, includeUpperLettersVar, includeDigitsVar, includeSymbolsVar) {
		return nil
	}

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	generatorType := d.Get(generatorTypeVar).(string)
	genType, ok := generatorTypeMap[generatorType]
	if !ok {
		return diag.Errorf("invalid generator_type %q", generatorType)
	}

	// The proto uses value types for all fields (except expiry), so every
	// field must be sent. Seed from the current server state and only
	// override fields that actually changed.
	current, err := client.GetSecretGenerator(ctx, &admin.GetSecretGeneratorRequest{
		GeneratorType: genType,
	})
	if err != nil {
		return diag.Errorf("failed to get current secret generator state: %v", err)
	}

	sg := current.GetSecretGenerator()
	req := &admin.UpdateSecretGeneratorRequest{
		GeneratorType:       genType,
		Length:              sg.GetLength(),
		Expiry:              sg.GetExpiry(),
		IncludeLowerLetters: sg.GetIncludeLowerLetters(),
		IncludeUpperLetters: sg.GetIncludeUpperLetters(),
		IncludeDigits:       sg.GetIncludeDigits(),
		IncludeSymbols:      sg.GetIncludeSymbols(),
	}

	if d.HasChange(lengthVar) {
		req.Length = uint32(d.Get(lengthVar).(int))
	}
	if d.HasChange(expiryVar) {
		expiry, err := time.ParseDuration(d.Get(expiryVar).(string))
		if err != nil {
			return diag.Errorf("failed to parse expiry: %v", err)
		}
		req.Expiry = durationpb.New(expiry)
	}
	if d.HasChange(includeLowerLettersVar) {
		req.IncludeLowerLetters = d.Get(includeLowerLettersVar).(bool)
	}
	if d.HasChange(includeUpperLettersVar) {
		req.IncludeUpperLetters = d.Get(includeUpperLettersVar).(bool)
	}
	if d.HasChange(includeDigitsVar) {
		req.IncludeDigits = d.Get(includeDigitsVar).(bool)
	}
	if d.HasChange(includeSymbolsVar) {
		req.IncludeSymbols = d.Get(includeSymbolsVar).(bool)
	}

	_, err = client.UpdateSecretGenerator(ctx, req)
	if helper.IgnorePreconditionError(err) != nil {
		return diag.Errorf("failed to update secret generator: %v", err)
	}

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	generatorType := d.Get(generatorTypeVar).(string)
	genType, ok := generatorTypeMap[generatorType]
	if !ok {
		return diag.Errorf("invalid generator_type %q", generatorType)
	}

	resp, err := client.GetSecretGenerator(ctx, &admin.GetSecretGeneratorRequest{
		GeneratorType: genType,
	})
	if err != nil {
		if helper.IgnoreIfNotFoundError(err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get secret generator: %v", err)
	}

	gen := resp.GetSecretGenerator()

	var expiryStr string
	if gen.GetExpiry() != nil {
		expiryStr = gen.GetExpiry().AsDuration().String()
	}

	set := map[string]interface{}{
		generatorTypeVar:       generatorType,
		lengthVar:              int(gen.GetLength()),
		expiryVar:              expiryStr,
		includeLowerLettersVar: gen.GetIncludeLowerLetters(),
		includeUpperLettersVar: gen.GetIncludeUpperLetters(),
		includeDigitsVar:       gen.GetIncludeDigits(),
		includeSymbolsVar:      gen.GetIncludeSymbols(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of secret generator: %v", k, err)
		}
	}

	d.SetId(generatorType)
	return nil
}
