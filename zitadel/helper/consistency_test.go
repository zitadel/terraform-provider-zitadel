package helper

import (
	"context"
	"errors"
	"testing"
)

func TestRetryUntilFound_ReturnsOnceFound(t *testing.T) {
	calls := 0
	found, err := RetryUntilFound(context.Background(), func() (bool, error) {
		calls++
		return calls == 3, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected found")
	}
	if calls != 3 {
		t.Fatalf("expected 3 lookups, got %d", calls)
	}
}

func TestRetryUntilFound_PropagatesError(t *testing.T) {
	want := errors.New("boom")
	calls := 0
	found, err := RetryUntilFound(context.Background(), func() (bool, error) {
		calls++
		return false, want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
	if found {
		t.Fatal("expected not found")
	}
	if calls != 1 {
		t.Fatalf("expected a single lookup, got %d", calls)
	}
}

func TestRetryUntilFound_StopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	found, err := RetryUntilFound(ctx, func() (bool, error) {
		return false, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if found {
		t.Fatal("expected not found")
	}
}
