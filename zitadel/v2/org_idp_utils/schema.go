package org_idp_utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	OrgIDVar = "org_id"
)

var (
	OrgIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the organization",
		ForceNew:    true,
	}
	OrgIDDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the organization",
	}
)
