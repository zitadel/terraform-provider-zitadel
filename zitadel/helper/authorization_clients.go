package helper

import (
	"context"
	"fmt"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	authorization "github.com/zitadel/zitadel-go/v3/pkg/client/authorization/v2"
	permission "github.com/zitadel/zitadel-go/v3/pkg/client/internal_permission/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

// These connections belong to one configured provider, never to a package-wide
// singleton: different provider aliases can use different instances or callers.
func GetAuthorizationV2Client(ctx context.Context, info *ClientInfo) (*authorization.Client, error) {
	info.bindingClientMu.Lock()
	defer info.bindingClientMu.Unlock()
	if info.authorizationClient == nil {
		client, err := authorization.NewClient(ctx, info.Issuer, info.Domain,
			[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, info.Options...)
		if err != nil {
			return nil, fmt.Errorf("failed to start authorization v2 client: %w", err)
		}
		info.authorizationClient = client
	}
	return info.authorizationClient, nil
}

func GetInternalPermissionV2Client(ctx context.Context, info *ClientInfo) (*permission.Client, error) {
	info.bindingClientMu.Lock()
	defer info.bindingClientMu.Unlock()
	if info.internalPermissionClient == nil {
		client, err := permission.NewClient(ctx, info.Issuer, info.Domain,
			[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, info.Options...)
		if err != nil {
			return nil, fmt.Errorf("failed to start internal permission v2 client: %w", err)
		}
		info.internalPermissionClient = client
	}
	return info.internalPermissionClient, nil
}

// BindingClientFactory keeps resource lifecycle tests on the real gRPC contract
// while allowing an in-process service instead of an authenticated instance.
type BindingClientFactory[T any] func(context.Context, interface{}) (T, error)
