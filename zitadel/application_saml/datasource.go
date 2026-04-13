package application_saml

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a SAML application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			AppIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the application",
			},
			MetadataXMLVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metadata as XML file",
			},
			LoginVersionVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specify the preferred login UI, where the user is redirected to for authentication. If unset, the login UI is chosen by the instance default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						LoginV1Var: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Login V1",
						},
						LoginV2Var: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Login V2",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									BaseURIVar: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Optionally specify a base uri of the login UI. If unspecified the default URI will be used.",
									},
								},
							},
						},
					},
				},
			},
		},
		ReadContext: read,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing multiple SAML applications belonging to a project.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			appIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			nameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying applications by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
		},
		ReadContext: list,
	}
}
