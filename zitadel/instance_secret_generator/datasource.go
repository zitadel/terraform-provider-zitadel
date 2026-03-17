package instance_secret_generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func GetDatasource() *schema.Resource {
	generatorTypes := make([]string, 0, len(generatorTypeMap))
	for k := range generatorTypeMap {
		generatorTypes = append(generatorTypes, k)
	}
	sort.Strings(generatorTypes)

	return &schema.Resource{
		Description: "Datasource representing a secret generator configuration.",
		Schema: map[string]*schema.Schema{
			generatorTypeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      fmt.Sprintf("Type of the secret generator, supported values: %s", strings.Join(generatorTypes, ", ")),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(generatorTypes, false)),
			},
			lengthVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Length of the generated secret",
			},
			expiryVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry duration of the generated secret, e.g. 1h, 15m, 24h",
			},
			includeLowerLettersVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Include lowercase letters in the generated secret",
			},
			includeUpperLettersVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Include uppercase letters in the generated secret",
			},
			includeDigitsVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Include digits in the generated secret",
			},
			includeSymbolsVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Include symbols in the generated secret",
			},
		},
		ReadContext: read,
	}
}
