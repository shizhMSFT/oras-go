/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package content

import (
	"context"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Resolver resolves reference tags.
type Resolver interface {
	// Resolve resolves a reference to a descriptor.
	Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error)
}

// TagResolver provides reference tag indexing services.
type TagResolver interface {
	Resolver

	// Tag tags a descriptor with a reference string.
	Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error
}

// ResolveFetcher provides advanced pull with the resolving service.
type ResolveFetcher interface {
	// ResolveFetch fetches the content with a reference string.
	// It is equivalent to call `Resolve()` and then `Fetch()` but more
	// efficient or equal.
	ResolveFetch(ctx context.Context, reference string) (ocispec.Descriptor, io.ReadCloser, error)
}

// TagPusher provides advanced push with the tag service.
type TagPusher interface {
	// PushTag pushes the content with a reference string.
	// It is equivalent to call `Push()` and then `Tag()` but more efficient or
	// equal.
	PushTag(ctx context.Context, expected ocispec.Descriptor, content io.Reader, reference string) error
}
