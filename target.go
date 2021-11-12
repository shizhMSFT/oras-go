package oras

import (
	"context"
	"io"

	"oras.land/oras-go/dag"
)

// Resolver resolves reference tags.
type Resolver interface {
	// Resolve resolves a reference to a descriptor.
	Resolve(ctx context.Context, reference string) (dag.Descriptor, error)
}

// TagResolver provides reference tag indexing services.
type TagResolver interface {
	Resolver

	// Tag tags a descriptor with a reference string.
	Tag(ctx context.Context, desc dag.Descriptor, reference string) error
}

// Target is a CAS with tags.
type Target interface {
	dag.Storage
	TagResolver
}

// ParentFinder finds out the parent nodes of a given node of a directed acyclic
// graph.
// ParentFinder is an extension of Storage.
type ParentFinder interface {
	FindParent(ctx context.Context, node dag.Descriptor) ([]dag.Descriptor, error)
}

// TagPusher pushes content with a reference tag.
// TagPusher is an extension of Target.
type TagPusher interface {
	// PushWithTag pushes the content with a reference tag, matching the
	// expected descriptor.
	PushWithTag(ctx context.Context, expected dag.Descriptor, content io.Reader, reference string) error
}
