# Tech Stack

**Reading Time**: ~4 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: H

## Context

ORAS Go maintains a minimal dependency footprint, relying primarily on official OCI specifications and standard Go libraries. This lightweight approach ensures long-term stability and reduces supply chain risks.

## Related Documentation

- **Previous**: [External Resources](06-external-resources.md) - Storage backends
- **Next**: [Getting Started](08-getting-started.md) - Setup and verification
- **Related**: [go.mod](../../../go.mod) - Dependency declarations

## Go Version Requirements

**Required**: Go 1.24.0 or Go 1.25

**Source**: [go.mod](../../../go.mod#L3)

**Policy**: Supports 2 latest Go versions per [Go Security Policy](https://github.com/golang/go/security/policy)

**Reference**: [README.md](../../../README.md#L11-L12)

## Core Dependencies

From [go.mod](../../../go.mod#L5-L10):

| Dependency | Version | Purpose | Key Types/Functions |
|------------|---------|---------|---------------------|
| `github.com/opencontainers/go-digest` | v1.0.0 | Digest computation and validation | `digest.Digest`, `digest.Algorithm` |
| `github.com/opencontainers/image-spec` | v1.1.1 | OCI specification types | `ocispec.Descriptor`, `ocispec.Manifest`, `ocispec.Index` |
| `golang.org/x/sync` | v0.19.0 | Advanced concurrency primitives | `semaphore.Weighted`, `errgroup.Group` |

### Dependency Details

#### 1. github.com/opencontainers/go-digest v1.0.0

**Purpose**: Cryptographic digest operations for content-addressable storage

**Usage in ORAS**:
- Computing SHA256/SHA512 digests for blobs and manifests
- Validating content integrity during fetch operations
- Parsing digest strings (e.g., `sha256:abc123...`)

**Key Types**:
- `digest.Digest` - Represents a cryptographic hash
- `digest.Algorithm` - Hash algorithm (SHA256, SHA512)

#### 2. github.com/opencontainers/image-spec v1.1.1

**Purpose**: Official OCI specification types and constants

**Usage in ORAS**:
- `ocispec.Descriptor` - Content descriptors for CAS operations
- `ocispec.Manifest` - Image manifest structure
- `ocispec.Index` - Image index (manifest list)
- `ocispec.Platform` - OS/architecture specifications
- Media type constants (e.g., `MediaTypeImageManifest`)

**Spec Version**: v1.1.1 (latest stable OCI spec)

**Note**: This is the authoritative source for OCI types, ensuring full standards compliance.

#### 3. golang.org/x/sync v0.19.0

**Purpose**: Advanced concurrency patterns beyond standard library

**Usage in ORAS**:
- `semaphore.Weighted` - Limit concurrent operations (e.g., parallel blob downloads)
- `errgroup.Group` - Coordinate goroutines with error handling

**Used In**:
- [copy.go](../../../copy.go) - Concurrent copy operations
- [content/graph.go](../../../content/graph.go) - Parallel graph traversal
- [internal/syncutil/](../../../internal/syncutil/) - Custom concurrency utilities

## OCI Standards Compliance

ORAS Go implements three core OCI specifications:

| Specification | Version | Purpose | Implementation |
|---------------|---------|---------|----------------|
| **OCI Image Format** | v1.1.1 | Container image structure | Manifest/Index handling in [content/](../../../content/) |
| **OCI Distribution** | v1.1.1 | Registry API protocol | HTTP client in [registry/remote/](../../../registry/remote/) |
| **OCI Image Layout** | v1.1.1 | On-disk format | Directory structure in [content/oci/](../../../content/oci/) |

**References**:
- [OCI Image Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md)
- [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec/blob/v1.1.1/spec.md)
- [OCI Image Layout Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/image-layout.md)

## HTTP & Web Standards

ORAS Go's remote registry client adheres to these standards:

| Standard | Purpose |
|----------|---------|
| **RFC 7230-7235** | HTTP/1.1 protocol |
| **RFC 6750** | OAuth 2.0 Bearer Token authentication |
| **RFC 3339** | DateTime format for timestamps |

## Media Types

ORAS Go recognizes and generates standard OCI media types:

### OCI Media Types
- `application/vnd.oci.image.manifest.v1+json` - OCI manifest
- `application/vnd.oci.image.index.v1+json` - OCI index
- `application/vnd.oci.image.config.v1+json` - OCI config
- `application/vnd.oci.descriptor.v1+json` - OCI descriptor

### Docker Media Types (Compatibility)
- `application/vnd.docker.distribution.manifest.v2+json` - Docker manifest v2
- `application/vnd.docker.distribution.manifest.list.v2+json` - Docker manifest list
- `application/vnd.docker.container.image.v1+json` - Docker config

**Reference**: [internal/docker/mediatype.go](../../../internal/docker/mediatype.go)

## Development Tools

### Build System

**Makefile**: [Makefile](../../../Makefile)

Key targets:

| Target | Purpose | Command |
|--------|---------|---------|
| `test` | Run tests with race detection + coverage | `make test` |
| `test-coverage` | View coverage details | `make test-coverage` |
| `clean` | Remove build artifacts | `make clean` |
| `check-encoding` | Verify no CRLF line endings | `make check-encoding` |
| `fix-encoding` | Fix line endings to LF | `make fix-encoding` |
| `vendor` | Vendor dependencies | `make vendor` |

### Testing Stack

**Framework**: Go standard `testing` package

**Test Types**:
- Unit tests (`*_test.go`)
- Table-driven tests
- Example tests (godoc) (`example_*_test.go`)
- Race detection (`go test -race`)

**Test Coverage**:
- **Required**: 80%+ for CI to pass
- **Current**: ~80-90%

**Coverage Script**: [scripts/coverage.sh](../../../scripts/coverage.sh)

### CI/CD Tools

**Platform**: GitHub Actions

**Workflows**:
- Automated testing (multiple Go versions)
- CodeQL security analysis
- Codecov coverage reporting
- Dependabot dependency updates

## Module Configuration

**Module Path**: `github.com/oras-project/oras-go/v3`

**Import Path**: `oras.land/oras-go/v3` (vanity URL)

**Versioning**: Semantic versioning (MAJOR.MINOR.PATCH)

**Branches**:
- `v3`: Current development (main branch) - **breaking changes expected**
- `v2`: Stable production - **recommended for production use**
- `v1`: Maintenance only

**Reference**: [README.md](../../../README.md#L44-L94)

## No External Runtime Dependencies

**Self-Contained**: All core functionality works without external services or daemons.

**Optional External**:
- Remote OCI registries (only for remote operations)
- OS credential stores (only for native credential storage)

This design ensures:
- Easy testing and development
- Minimal deployment requirements
- High reliability (no service dependencies)
- Cross-platform compatibility

## Next Steps

Continue with [Getting Started](08-getting-started.md) to set up and verify your ORAS Go development environment.

## Citations

- [go.mod](../../../go.mod) - Dependency declarations
- [Makefile](../../../Makefile) - Build automation
- [scripts/coverage.sh](../../../scripts/coverage.sh) - Coverage reporting
- [README.md](../../../README.md) - Project overview and versioning
- [OCI Image Spec v1.1.1](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md)
- [OCI Distribution Spec v1.1.1](https://github.com/opencontainers/distribution-spec/blob/v1.1.1/spec.md)
