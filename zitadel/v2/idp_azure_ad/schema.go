package idp_azure_ad

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	TenantTypeVar    = "tenant_type"
	TenantIDVar      = "tenant_id"
	EmailVerifiedVar = "email_verified"
)

var (
	TenantTypeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "AZURE_AD_TENANT_TYPE_COMMON",
		Description: "the azure ad tenant type",
	}
	TenantTypeDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the azure ad tenant type",
	}
	TenantIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: fmt.Sprintf("if %s is not set, the %s is used", TenantIDVar, TenantTypeVar),
	}
	TenantIDDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the azure ad tenant id",
	}
	EmailVerifiedResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "automatically mark emails as verified",
	}
	EmailVerifiedDataSourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "automatically mark emails as verified",
	}
)
