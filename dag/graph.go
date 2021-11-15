// Package dag stores and traverses directed acyclic graphs (DAGs).
package dag

import (
	"context"
	"encoding/json"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	artifactspec "github.com/oras-project/artifacts-spec/specs-go/v1"
	"oras.land/oras-go/internal/docker"
)

// DownEdges returns the nodes directly pointed by the current node.
// In other words, returns the "children" of the current descriptor.
func DownEdges(ctx context.Context, fetcher Fetcher, node Descriptor) ([]Descriptor, error) {
	switch node.MediaType {
	case docker.MediaTypeManifest, ocispec.MediaTypeImageManifest:
		content, err := FetchContent(ctx, fetcher, node)
		if err != nil {
			return nil, err
		}

		// docker manifest and oci manifest are equivalent for down edges.
		var manifest ocispec.Manifest
		if err := json.Unmarshal(content, &manifest); err != nil {
			return nil, err
		}
		return append([]Descriptor{manifest.Config}, manifest.Layers...), nil
	case docker.MediaTypeManifestList, ocispec.MediaTypeImageIndex:
		content, err := FetchContent(ctx, fetcher, node)
		if err != nil {
			return nil, err
		}

		// docker manifest list and oci index are equivalent for down edges.
		var index ocispec.Index
		if err := json.Unmarshal(content, &index); err != nil {
			return nil, err
		}
		return index.Manifests, nil
	case artifactspec.MediaTypeArtifactManifest:
		content, err := FetchContent(ctx, fetcher, node)
		if err != nil {
			return nil, err
		}

		var manifest artifactspec.Manifest
		if err := json.Unmarshal(content, &manifest); err != nil {
			return nil, err
		}
		var nodes []Descriptor
		if !isEmptyArtifactDescriptor(manifest.Subject) {
			nodes = append(nodes, convertArtifactDescriptorToOCI(manifest.Subject))
		}
		for _, blob := range manifest.Blobs {
			nodes = append(nodes, convertArtifactDescriptorToOCI(blob))
		}
		return nodes, nil
	}
	return nil, nil
}

func isEmptyArtifactDescriptor(desc artifactspec.Descriptor) bool {
	return desc.MediaType == "" &&
		desc.ArtifactType == "" &&
		desc.Digest == "" &&
		desc.Size == 0 &&
		desc.URLs == nil &&
		desc.Annotations == nil
}

func convertArtifactDescriptorToOCI(desc artifactspec.Descriptor) ocispec.Descriptor {
	return ocispec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Annotations: desc.Annotations,
	}
}
