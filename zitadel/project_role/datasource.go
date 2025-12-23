package project_role

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the project roles, which can be given as authorizations to users.",
		Schema: map[string]*schema.Schema{
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			KeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used for project role",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name used for project role",
			},
			groupVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group used for project role",
			},
		},
		ReadContext: read,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all project roles in a project, which can be given as authorizations to users. Note: Group-based filtering is not supported by the ZITADEL API.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the organization",
			},
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			roleKeysVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all project role keys",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			KeyVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key to filter project roles by",
			},
			keyMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying project roles by key" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(keyMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name to filter project roles by",
			},
			displayNameMethod: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying project roles by display name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(displayNameMethod, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
		},
		ReadContext: list,
	}
}
