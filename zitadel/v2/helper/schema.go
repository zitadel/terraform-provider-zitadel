package helper

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ResourceIDVar = "id"
	OrgIDVar      = "org_id"
)

var (
	// ZitadelGeneratedIdPattern matches IDs like 123456789012345678
	// ZITADEL IDs have 18 digits
	ZitadelGeneratedIdPattern   = `\d{18}`
	ZitadelGeneratedIdOnlyRegex = regexp.MustCompile(fmt.Sprintf(`^%s$`, ZitadelGeneratedIdPattern))

	OrgIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the organization",
		ForceNew:    true,
	}

	ResourceIDDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the resource",
	}
	OrgIDDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the organization",
	}
)
