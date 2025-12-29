package instance_trusted_domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a trusted domain on a ZITADEL instance. " +
			"Trusted domains can be used in API responses like OIDC discovery and email templates. " +
			"Unlike custom domains, trusted domains are not used for routing and do not need to be unique across instances. " +
			"This resource requires iam.write permission.",
		Schema: map[string]*schema.Schema{
			InstanceIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the instance. If not provided, the instance from the current context will be used.",
				ForceNew:    true,
			},
			DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trusted domain to add to the instance (max 253 characters)",
				ForceNew:    true,
			},
		},
		CreateContext: create,
		ReadContext:   read,
		DeleteContext: delete,
		Importer:      helper.ImportWithID(InstanceIDVar),
	}
}
