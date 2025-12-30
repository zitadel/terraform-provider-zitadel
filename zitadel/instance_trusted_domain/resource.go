package instance_trusted_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
