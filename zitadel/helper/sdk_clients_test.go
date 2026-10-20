package helper

import (
	"context"
	"sync"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

// The embedded SDK can configure multiple aliases in one process. Each
// configuration must retain its own connection, including after token rotation.
func TestSDKClientsAreScopedToConfiguration(t *testing.T) {
	getters := map[string]func(context.Context, *ClientInfo) (*grpc.ClientConn, error){
		"project": func(ctx context.Context, info *ClientInfo) (*grpc.ClientConn, error) {
			c, err := GetProjectV2Client(ctx, info)
			if err != nil {
				return nil, err
			}
			return c.Connection.ClientConn, nil
		},
		"application": func(ctx context.Context, info *ClientInfo) (*grpc.ClientConn, error) {
			c, err := GetAppV2Client(ctx, info)
			if err != nil {
				return nil, err
			}
			return c.Connection.ClientConn, nil
		},
		"authorization": func(ctx context.Context, info *ClientInfo) (*grpc.ClientConn, error) {
			c, err := GetAuthorizationV2Client(ctx, info)
			if err != nil {
				return nil, err
			}
			return c.Connection.ClientConn, nil
		},
		"administrator": func(ctx context.Context, info *ClientInfo) (*grpc.ClientConn, error) {
			c, err := GetInternalPermissionV2Client(ctx, info)
			if err != nil {
				return nil, err
			}
			return c.Connection.ClientConn, nil
		},
	}
	for name, get := range getters {
		t.Run(name, func(t *testing.T) {
			var previous *grpc.ClientConn
			for _, token := range []string{"first", "rotated"} {
				info := &ClientInfo{Domain: "127.0.0.1:1", Issuer: "http://127.0.0.1:1", Options: []zitadel.Option{
					zitadel.WithInsecure(), zitadel.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})),
				}}
				first, err := get(context.Background(), info)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = first.Close() })
				if first == previous {
					t.Fatal("different configurations shared a connection")
				}
				previous = first
				var calls sync.WaitGroup
				for range 8 {
					calls.Go(func() {
						got, err := get(context.Background(), info)
						if err != nil || got != first {
							t.Error("one configuration did not reuse its connection")
						}
					})
				}
				calls.Wait()
			}
		})
	}
}
