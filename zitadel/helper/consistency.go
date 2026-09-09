package helper

import (
	"context"
	"time"
)

const (
	// eventualConsistencyTimeout bounds how long a read waits for ZITADEL's
	// eventually consistent v2 query APIs to reflect a preceding write.
	eventualConsistencyTimeout = 15 * time.Second
	// eventualConsistencyInterval is the pause between two lookups.
	eventualConsistencyInterval = 500 * time.Millisecond
)

// RetryUntilFound calls lookup until it reports the object as found, returns an
// error, or the consistency timeout elapses. ZITADEL's v2 query APIs are served
// from projections that can lag behind a write for a short moment, so an object
// created by a preceding resource may not be visible immediately. Callers keep
// their not-found semantics: when the timeout elapses, found is false and err is nil.
func RetryUntilFound(ctx context.Context, lookup func() (found bool, err error)) (bool, error) {
	deadline := time.Now().Add(eventualConsistencyTimeout)
	for {
		found, err := lookup()
		if err != nil || found {
			return found, err
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(eventualConsistencyInterval):
		}
	}
}
