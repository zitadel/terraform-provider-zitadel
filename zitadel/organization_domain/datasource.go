package organization_domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a domain of an organization in ZITADEL.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			IsVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the domain has been verified",
			},
			IsPrimaryVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the primary domain of the organization",
			},
			ValidationTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of domain validation" + helper.DescriptionEnumValuesList(org.DomainValidationType_name),
			},
		},
		ReadContext: get,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing domains of an organization in ZITADEL.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			DomainVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter domains by name",
			},
			domainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of organization domains",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						OrganizationIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the organization",
						},
						DomainVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name",
						},
						IsVerifiedVar: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the domain has been verified",
						},
						IsPrimaryVar: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the primary domain of the organization",
						},
						ValidationTypeVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of domain validation",
						},
					},
				},
			},
		},
		ReadContext: list,
	}
}
