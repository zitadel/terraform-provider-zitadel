package webkey

import (
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a web key.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			RSABlock: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						BitsVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "RSA_BITS_2048",
							ForceNew:    true,
							Description: "Bit size of the RSA key.",
							ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
								val := value.(string)
								validValues := []string{
									"RSA_BITS_2048",
									"RSA_BITS_3072",
									"RSA_BITS_4096",
								}
								for _, valid := range validValues {
									if val == valid {
										return nil
									}
								}
								return diag.Errorf("%s: invalid value %s, allowed values: %s", path, val, strings.Join(validValues, ", "))
							},
						},
						HasherVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "RSA_HASHER_SHA256",
							ForceNew:    true,
							Description: "Signing algorithm used.",
							ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
								val := value.(string)
								validValues := []string{
									"RSA_HASHER_SHA256",
									"RSA_HASHER_SHA384",
									"RSA_HASHER_SHA512",
								}
								for _, valid := range validValues {
									if val == valid {
										return nil
									}
								}
								return diag.Errorf("%s: invalid value %s, allowed values: %s", path, val, strings.Join(validValues, ", "))
							},
						},
					},
				},
				ExactlyOneOf: []string{RSABlock, ECDSABlock, ED25519Block},
			},
			ECDSABlock: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						CurveVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "ECDSA_CURVE_P256",
							ForceNew:    true,
							Description: "Curve of the ECDSA key.",
							ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
								val := value.(string)
								validValues := []string{
									"ECDSA_CURVE_P256",
									"ECDSA_CURVE_P384",
									"ECDSA_CURVE_P512",
								}
								for _, valid := range validValues {
									if val == valid {
										return nil
									}
								}
								return diag.Errorf("%s: invalid value %s, allowed values: %s", path, val, strings.Join(validValues, ", "))
							},
						},
					},
				},
			},
			ED25519Block: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Create a ED25519 key pair.",
				Elem:        &schema.Resource{Schema: map[string]*schema.Schema{}},
			},
			StateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the key.",
			},
			KeyTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the key.",
			},
		},
		CreateContext: create,
		ReadContext:   read,
		DeleteContext: delete,
		Importer:      helper.ImportWithIDAndOptionalOrg(WebKeyIDVar),
	}
}
