package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	idpJwtEndpoint  = "jwt_endpoint"
	idpKeysEndpoint = "keys_endpoint"
	idpHeaderName   = "header_name"
)

func GetOrgJWTIDP() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a domain of the organization.",
		Schema: map[string]*schema.Schema{
			idpOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			idpNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the IDP",
			},
			idpStylingTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Some identity providers specify the styling of the button to their login",
			},
			idpJwtEndpoint: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the endpoint where the jwt can be extracted",
			},
			idpKeysEndpoint: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the endpoint to the key (JWK) which are used to sign the JWT with",
			},
			idpIssuerVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the issuer of the jwt (for validation)",
			},
			idpHeaderName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of the header where the JWT is sent in, default is authorization",
			},
			idpAutoRegister: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "auto register for users from this idp",
			},
		},
		ReadContext:   readOrgJWTIDP,
		CreateContext: createOrgJWTIDP,
		UpdateContext: updateOrgJWTIDP,
		DeleteContext: deleteOrgIDP,
	}
}
func createOrgJWTIDP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(domainOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	stylingType := d.Get(idpStylingTypeVar)
	resp, err := client.AddOrgJWTIDP(ctx, &management2.AddOrgJWTIDPRequest{
		Name:         d.Get(idpNameVar).(string),
		StylingType:  idp.IDPStylingType(idp.IDPStylingType_value[stylingType.(string)]),
		JwtEndpoint:  d.Get(idpJwtEndpoint).(string),
		Issuer:       d.Get(idpIssuerVar).(string),
		KeysEndpoint: d.Get(idpKeysEndpoint).(string),
		HeaderName:   d.Get(idpHeaderName).(string),
		AutoRegister: d.Get(idpAutoRegister).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create jwt idp: %v", err)
	}
	d.SetId(resp.IdpId)
	return nil
}

func updateOrgJWTIDP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(idpOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetOrgIDPByID(ctx, &management2.GetOrgIDPByIDRequest{Id: d.Get("id").(string)})
	if err != nil {
		return diag.Errorf("failed to read jwt idp: %v", err)
	}

	idpID := d.Id()
	name := d.Get(idpNameVar).(string)
	stylingType := d.Get(idpStylingTypeVar).(string)
	autoRegister := d.Get(idpAutoRegister).(bool)
	if resp.GetIdp().GetName() != name ||
		resp.GetIdp().GetStylingType().String() != stylingType ||
		resp.GetIdp().GetAutoRegister() != autoRegister {
		_, err := client.UpdateOrgIDP(ctx, &management2.UpdateOrgIDPRequest{
			IdpId:        idpID,
			Name:         name,
			StylingType:  idp.IDPStylingType(idp.IDPStylingType_value[stylingType]),
			AutoRegister: autoRegister,
		})
		if err != nil {
			return diag.Errorf("failed to update jwt idp: %v", err)
		}
	}

	jwt := resp.GetIdp().GetJwtConfig()
	jwtEndpoint := d.Get(idpJwtEndpoint).(string)
	issuer := d.Get(idpIssuerVar).(string)
	keysEndpoint := d.Get(idpKeysEndpoint).(string)
	headerName := d.Get(idpHeaderName).(string)

	//either nothing changed on the IDP or something besides the secret changed
	if jwt.GetJwtEndpoint() != jwtEndpoint ||
		jwt.GetIssuer() != issuer ||
		jwt.GetKeysEndpoint() != keysEndpoint ||
		jwt.GetHeaderName() != headerName {

		_, err = client.UpdateOrgIDPJWTConfig(ctx, &management2.UpdateOrgIDPJWTConfigRequest{
			IdpId:        idpID,
			JwtEndpoint:  jwtEndpoint,
			Issuer:       issuer,
			KeysEndpoint: keysEndpoint,
			HeaderName:   headerName,
		})
		if err != nil {
			return diag.Errorf("failed to update jwt idp config: %v", err)
		}
	}
	return nil
}

func readOrgJWTIDP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(idpOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetOrgIDPByID(ctx, &management2.GetOrgIDPByIDRequest{Id: d.Id()})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read jwt idp: %v", err)
	}

	idp := resp.GetIdp()
	jwt := idp.GetJwtConfig()
	set := map[string]interface{}{
		idpOrgIDVar:       idp.GetDetails().ResourceOwner,
		idpNameVar:        idp.GetName(),
		idpStylingTypeVar: idp.GetStylingType().String(),
		idpJwtEndpoint:    jwt.GetJwtEndpoint(),
		idpIssuerVar:      jwt.GetIssuer(),
		idpKeysEndpoint:   jwt.GetKeysEndpoint(),
		idpHeaderName:     jwt.GetHeaderName(),
		idpAutoRegister:   idp.GetAutoRegister(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of jwt idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)

	return nil
}
