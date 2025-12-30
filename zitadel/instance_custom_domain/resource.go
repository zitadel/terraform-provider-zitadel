package instance_custom_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a custom domain on a ZITADEL instance. " +
			"Custom domains are used to route requests to the instance and must be unique across all instances. " +
			"This resource requires system-level permissions (system.domain.write).",
		Schema: map[string]*schema.Schema{
			InstanceIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the instance",
				ForceNew:    true,
			},
			DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom domain to add to the instance (max 253 characters)",
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
