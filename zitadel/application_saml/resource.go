package application_saml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a SAML application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			MetadataXMLVar: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Metadata as XML file",
				Sensitive:     true,
				ConflictsWith: []string{MetadataURLVar},
				AtLeastOneOf:  []string{MetadataXMLVar, MetadataURLVar},
			},
			MetadataURLVar: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Metadata URL to fetch the SAML metadata from",
				ConflictsWith: []string{MetadataXMLVar},
				AtLeastOneOf:  []string{MetadataXMLVar, MetadataURLVar},
			},
			LoginVersionVar: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Specify the preferred login UI, where the user is redirected to for authentication. If unset, the login UI is chosen by the instance default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						LoginV1Var: {
							Type:          schema.TypeBool,
							Optional:      true,
							Description:   "Login V1",
							ConflictsWith: []string{LoginVersionVar + ".0." + LoginV2Var},
						},
						LoginV2Var: {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Description:   "Login V2",
							ConflictsWith: []string{LoginVersionVar + ".0." + LoginV1Var},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									BaseURIVar: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optionally specify a base uri of the login UI. If unspecified the default URI will be used.",
									},
								},
							},
						},
					},
				},
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer: helper.ImportWithIDAndOptionalOrg(
			AppIDVar,
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertID, false),
		),
	}
}
