package project

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the project, which can then be granted to different organizations or users directly, containing different applications.",
		Schema: map[string]*schema.Schema{
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the project",
			},
			helper.OrgIDVar: helper.OrgIDResourceField,
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the project",
			},
			roleAssertionVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "describes if roles of user should be added in token",
			},
			roleCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the user has at least one on this project",
			},
			hasProjectCheckVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "ZITADEL checks if the org of the user has permission to this project",
			},
			privateLabelingSettingVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines from where the private labeling should be triggered",
			},
		},
		ReadContext: read,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the project, which can then be granted to different organizations or users directly, containing different applications.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			projectIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all project IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the project",
			},
			nameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying projects by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
		},
		ReadContext: list,
	}
}
