package oras

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
)

func GetManifest(ctx context.Context, target Target, reference string) ([]byte, ocispec.Descriptor, error) {
	desc, err := target.Resolve(ctx, reference)
	if err != nil {
		return nil, ocispec.Descriptor{}, err
	}
	content, err := content.FetchAll(ctx, target, desc)
	if err != nil {
		return nil, ocispec.Descriptor{}, err
	}
	return content, desc, nil
}
