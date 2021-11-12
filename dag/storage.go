package dag

import (
	"context"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Descriptor describes the disposition of targeted content.
type Descriptor = ocispec.Descriptor

// Fetcher fetches content.
type Fetcher interface {
	// Fetch fetches the content identified by the descriptor.
	Fetch(ctx context.Context, target Descriptor) (io.ReadCloser, error)
}

// Pusher pushes content.
type Pusher interface {
	// Push pushes the content, matching the expected descriptor.
	// Reader is perferred to Writer so that the suitable buffer size can be
	// chosen by the underlying implementation. Furthermore, the implementation
	// can also do reflection on the Reader for more advanced I/O optimization.
	Push(ctx context.Context, expected Descriptor, content io.Reader) error
}

// Storage represents a content-addressable storage (CAS) where contents are
// accessed via Descriptors.
// The storage is designed to handle blobs of large sizes.
type Storage interface {
	Fetcher
	Pusher

	// Exists returns true if the described content exists.
	Exists(ctx context.Context, target Descriptor) (bool, error)
}
