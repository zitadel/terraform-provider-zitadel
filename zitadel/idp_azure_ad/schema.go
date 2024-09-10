package idp_azure_ad

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
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
		Description: "the azure ad tenant type",
		ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
			return helper.EnumValueValidation(TenantTypeVar, value, idp.AzureADTenantType_value)
		},
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
