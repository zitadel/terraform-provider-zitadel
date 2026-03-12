package instance_secret_generator

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	generatorTypes := make([]string, 0, len(generatorTypeMap))
	for k := range generatorTypeMap {
		generatorTypes = append(generatorTypes, k)
	}
	sort.Strings(generatorTypes)

	return &schema.Resource{
		Description: "Resource representing a secret generator configuration.",
		Schema: map[string]*schema.Schema{
			generatorTypeVar: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      fmt.Sprintf("Type of the secret generator, supported values: %s", strings.Join(generatorTypes, ", ")),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(generatorTypes, false)),
			},
			lengthVar: {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				Description:      "Length of the generated secret",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			expiryVar: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Expiry duration of the generated secret, e.g. 1h, 15m, 24h",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			includeLowerLettersVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include lowercase letters in the generated secret",
			},
			includeUpperLettersVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include uppercase letters in the generated secret",
			},
			includeDigitsVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include digits in the generated secret",
			},
			includeSymbolsVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include symbols in the generated secret",
			},
		},
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				if _, ok := generatorTypeMap[id]; !ok {
					return nil, fmt.Errorf("invalid generator_type %q for import", id)
				}
				if err := d.Set(generatorTypeVar, id); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
