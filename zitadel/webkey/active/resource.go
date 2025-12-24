package active_webkey

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetActiveResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the active web key.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			KeyIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the webkey to be active.",
			},
		},
		CreateContext: createActiveWebKey,
		ReadContext:   readActiveWebKey,
		UpdateContext: updateActiveWebKey,
		DeleteContext: deleteActiveWebKey,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid import id %q, format must be <org_id>:<key_id>", d.Id())
				}
				if err := d.Set(helper.OrgIDVar, parts[0]); err != nil {
					return nil, err
				}
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
