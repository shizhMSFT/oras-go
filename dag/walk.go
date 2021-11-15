package dag

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// WalkFunc is a type of the function called by Walk to visit the content
// described by each descriptor.
type WalkFunc func(Descriptor) error

// Walk walks a rooted DAG concurrently.
//
func Walk(ctx context.Context, root Descriptor, fn WalkFunc, limiter *semaphore.Weighted) error {

}
