package helper

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIgnoreIfNotFoundError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "nil error returns nil",
			err:     nil,
			wantNil: true,
		},
		{
			name:    "not found error returns nil",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: true,
		},
		{
			name:    "already exists error is not ignored",
			err:     status.Error(codes.AlreadyExists, "already exists"),
			wantNil: false,
		},
		{
			name:    "internal error is not ignored",
			err:     status.Error(codes.Internal, "internal error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnoreIfNotFoundError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnoreIfNotFoundError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnoreIfNotFoundError() = nil, want non-nil")
			}
		})
	}
}

func TestIgnorePreconditionError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "precondition error returns nil",
			err:     status.Error(codes.FailedPrecondition, "failed precondition"),
			wantNil: true,
		},
		{
			name:    "not found error is not ignored",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: false,
		},
		{
			name:    "internal error is not ignored",
			err:     status.Error(codes.Internal, "internal error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnorePreconditionError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnorePreconditionError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnorePreconditionError() = nil, want non-nil")
			}
		})
	}
}

func TestIgnoreAlreadyExistsError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "already exists error returns nil",
			err:     status.Error(codes.AlreadyExists, "already exists"),
			wantNil: true,
		},
		{
			name:    "not found error is not ignored",
			err:     status.Error(codes.NotFound, "not found"),
			wantNil: false,
		},
		{
			name:    "plain error is not ignored",
			err:     fmt.Errorf("some error"),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnoreAlreadyExistsError(tt.err)
			if tt.wantNil && result != nil {
				t.Errorf("IgnoreAlreadyExistsError() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("IgnoreAlreadyExistsError() = nil, want non-nil")
			}
		})
	}
}
