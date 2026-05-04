package helper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsUnimplemented returns true if the error is a gRPC Unimplemented status code.
// This is used to detect when a V2 API is not available (e.g. on Zitadel server v3.x)
// and fallback to the legacy Management/Admin API should be attempted.
func IsUnimplemented(err error) bool {
	return status.Code(err) == codes.Unimplemented
}
