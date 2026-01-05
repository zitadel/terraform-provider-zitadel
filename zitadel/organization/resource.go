package organization

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the org",
			},
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom unique identifier for the organization",
			},
			adminsVar: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Description: "List of users to be granted organization admin roles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						adminUserIDVar: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the user to be added as org admin",
						},
						adminRolesVar: {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of roles to grant to the user. If empty, ORG_OWNER is granted by default",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			primaryDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the org",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the org",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   get,
		UpdateContext: update,
		Importer:      helper.ImportWithID(OrgIDVar),
	}
}
