package organization_domain

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsAlreadyVerifiedError(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
		want bool
	}{{
		name: "nil",
		err:  nil,
		want: false,
	}, {
		name: "already verified",
		err:  status.Error(codes.FailedPrecondition, "Domain is already verified (ORG-HGw21)"),
		want: true,
	}, {
		// The API reuses FailedPrecondition for rejections that have to fail
		// the apply, so the code alone must not be enough.
		name: "other precondition failure",
		err:  status.Error(codes.FailedPrecondition, "Domain is not unique (ORG-Df3bs)"),
		want: false,
	}, {
		name: "other code carrying the id",
		err:  status.Error(codes.Internal, "Domain is already verified (ORG-HGw21)"),
		want: false,
	}, {
		name: "not a status error",
		err:  errors.New("Domain is already verified (ORG-HGw21)"),
		want: false,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAlreadyVerifiedError(tt.err); got != tt.want {
				t.Errorf("isAlreadyVerifiedError() = %v, want %v", got, tt.want)
			}
		})
	}
}
