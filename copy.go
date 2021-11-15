package oras

import (
	"context"
	"errors"

	"oras.land/oras-go/dag"
)

// Copy copies a rooted directed acyclic graph (DAG) with the tagged root node
// in the source Target to the destination Target.
// The destination reference will be the same as the source reference if the
// destination reference is left blank.
// Returns the descriptor of the root node on successful copy.
func Copy(ctx context.Context, src Target, srcRef string, dst Target, dstRef string) (dag.Descriptor, error) {
	if src == nil {
		return dag.Descriptor{}, errors.New("nil source target")
	}
	if dst == nil {
		return dag.Descriptor{}, errors.New("nil destination target")
	}

	srcDesc, err := src.Resolve(ctx, srcRef)
	if err != nil {
		return dag.Descriptor{}, err
	}

}
