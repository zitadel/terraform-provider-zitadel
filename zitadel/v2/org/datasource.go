package org

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/org"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An organizations resource ID.",
			},
			nameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the org.",
			},
			domainVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A domain of the org.",
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
		ReadContext: queryDatasource,
	}
}
