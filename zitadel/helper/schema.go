package helper

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	OrgIDVar                         = "org_id"
	InsecureSkipVerifyTLSVar         = "insecure_skip_verify_tls"
	InsecureSkipVerifyTLSDescription = "Disable TLS certificate verification. Only use in development/testing environments with self-signed certificates."
	TransportHeadersVar              = "transport_headers"
	TransportHeadersDescription      = "Custom headers to add to both HTTP (authentication) and gRPC (API) requests. Useful for proxy authentication (e.g., GCP IAP with Proxy-Authorization header)."
)

var (
	// ZitadelGeneratedIdPattern matches IDs like 123456789012345678
	// ZITADEL IDs have 18 digits
	ZitadelGeneratedIdPattern   = `\d{18}`
	ZitadelGeneratedIdOnlyRegex = regexp.MustCompile(fmt.Sprintf(`^%s$`, ZitadelGeneratedIdPattern))

	OrgIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "ID of the organization. If not provided, the organization of the authenticated user/service account is used.",
		ForceNew:    true,
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			_, err := ConvertID(i.(string))
			return diag.FromErr(err)
		},
	}

	ResourceIDDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the resource",
	}
	OrgIDDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "ID of the organization",
	}
)
