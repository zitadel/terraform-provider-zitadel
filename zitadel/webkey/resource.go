package webkey

import (
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						BitsVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "RSA_BITS_2048",
							Description: "Bit size of the RSA key.",
						},
						HasherVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "RSA_HASHER_SHA256",
							Description: "Signing algorithm used.",
						},
					},
				},
				ExactlyOneOf: []string{RSABlock, ECDSABlock, ED25519Block},
				ForceNew:     true,
			},
			ECDSABlock: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						CurveVar: {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "ECDSA_CURVE_P256",
							Description: "Curve of the ECDSA key.",
						},
					},
				},
				ForceNew: true,
			},
			ED25519Block: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Create a ED25519 key pair.",
				Elem:        &schema.Resource{Schema: map[string]*schema.Schema{}},
				ForceNew:    true,
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
		CreateContext: createWebKey,
		ReadContext:   readWebKey,
		DeleteContext: deleteWebKey,
		Importer:      helper.ImportWithIDAndOptionalOrg(WebKeyIDVar),
	}
}
