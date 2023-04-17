package org_idp_github

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

type tfModel struct {
	// ID is never null or empty, so we can use the native type
	ID string `tfsdk:"id"`
	// OrgID is never null or empty, so we can use the native type
	OrgID             string       `tfsdk:"org_id"`
	Name              types.String `tfsdk:"name"`
	ClientID          types.String `tfsdk:"client_id"`
	Scopes            []string     `tfsdk:"scopes"`
	IsLinkingAllowed  types.Bool   `tfsdk:"is_linking_allowed"`
	IsCreationAllowed types.Bool   `tfsdk:"is_creation_allowed"`
	IsAutoCreation    types.Bool   `tfsdk:"is_auto_creation"`
	IsAutoUpdate      types.Bool   `tfsdk:"is_auto_update"`
}

// We have to duplicate the fields until the framework supports embedding
// https://github.com/hashicorp/terraform-plugin-framework/issues/242
type tfModelSensitive struct {
	// ID is never null or empty, so we can use the native type
	ID string `tfsdk:"id"`
	// OrgID is never null or empty, so we can use the native type
	OrgID             string       `tfsdk:"org_id"`
	Name              types.String `tfsdk:"name"`
	ClientID          types.String `tfsdk:"client_id"`
	Scopes            []string     `tfsdk:"scopes"`
	IsLinkingAllowed  types.Bool   `tfsdk:"is_linking_allowed"`
	IsCreationAllowed types.Bool   `tfsdk:"is_creation_allowed"`
	IsAutoCreation    types.Bool   `tfsdk:"is_auto_creation"`
	IsAutoUpdate      types.Bool   `tfsdk:"is_auto_update"`
	ClientSecret      types.String `tfsdk:"client_secret"`
}

func (m *tfModelSensitive) toPbAddGithubProviderRequest() *management.AddGitHubProviderRequest {
	return &management.AddGitHubProviderRequest{
		Name:            m.Name.ValueString(),
		ClientId:        m.ClientID.ValueString(),
		ClientSecret:    m.ClientSecret.ValueString(),
		Scopes:          m.Scopes,
		ProviderOptions: m.toPbIdpOptions(),
	}
}

func (m *tfModelSensitive) toPbUpdateGithubProviderRequest() *management.UpdateGitHubProviderRequest {
	return &management.UpdateGitHubProviderRequest{
		Id:              m.ID,
		Name:            m.Name.ValueString(),
		ClientId:        m.ClientID.ValueString(),
		ClientSecret:    m.ClientSecret.ValueString(),
		Scopes:          m.Scopes,
		ProviderOptions: m.toPbIdpOptions(),
	}
}

func (m *tfModelSensitive) toPbIdpOptions() *idp.Options {
	return &idp.Options{
		IsLinkingAllowed:  m.IsLinkingAllowed.ValueBool(),
		IsCreationAllowed: m.IsCreationAllowed.ValueBool(),
		IsAutoCreation:    m.IsAutoCreation.ValueBool(),
		IsAutoUpdate:      m.IsAutoUpdate.ValueBool(),
	}
}

func (m *tfModel) withClientSecret(secret string) *tfModelSensitive {
	return &tfModelSensitive{
		ID:                m.ID,
		OrgID:             m.OrgID,
		Name:              m.Name,
		ClientID:          m.ClientID,
		Scopes:            m.Scopes,
		IsLinkingAllowed:  m.IsLinkingAllowed,
		IsCreationAllowed: m.IsCreationAllowed,
		IsAutoCreation:    m.IsAutoCreation,
		IsAutoUpdate:      m.IsAutoUpdate,
		ClientSecret:      types.StringValue(secret),
	}
}

func fromPbGetProviderByIDResponse(resp *management.GetProviderByIDResponse, orgID string) *tfModel {
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetGithub()
	generalCfg := cfg.GetOptions()
	return &tfModel{
		ID:                idp.GetId(),
		OrgID:             orgID,
		Name:              types.StringValue(idp.GetName()),
		ClientID:          types.StringValue(specificCfg.GetClientId()),
		Scopes:            specificCfg.GetScopes(),
		IsLinkingAllowed:  types.BoolValue(generalCfg.GetIsLinkingAllowed()),
		IsCreationAllowed: types.BoolValue(generalCfg.GetIsCreationAllowed()),
		IsAutoCreation:    types.BoolValue(generalCfg.GetIsAutoCreation()),
		IsAutoUpdate:      types.BoolValue(generalCfg.GetIsAutoUpdate()),
	}
}
