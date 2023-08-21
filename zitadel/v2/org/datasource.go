package org

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					_, err := helper.ConvertID(i.(string))
					return diag.FromErr(err)
				},
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the org.",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the org" + helper.DescriptionEnumValuesList(org.OrgState_name),
			},
			primaryDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the org",
			},
		},
		ReadContext: get,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			orgIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all organization IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			NameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the org.",
			},
			nameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying orgs by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			DomainVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A domain of the org.",
			},
			domainMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying orgs by domain" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(domainMethodVar, value, object.TextQueryMethod_value)
				},
				InputDefault: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			stateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State of the org" + helper.DescriptionEnumValuesList(org.OrgState_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(stateVar, value, org.OrgState_value)
				},
			},
			primaryDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the org",
			},
		},
		ReadContext: list,
	}
}
