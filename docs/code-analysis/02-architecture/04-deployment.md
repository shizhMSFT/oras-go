# Deployment & Integration Patterns

**Reading Time**: ~8 min  
**Analysis Phase**: 2 - Architecture  
**Level**: Architecture  
**Confidence**: [H] High (90%)

## Context

This document describes how ORAS Go is deployed and integrated as a library SDK into various applications and workflows. Unlike standalone services, ORAS Go provides programmatic APIs for embedding OCI artifact management capabilities directly into Go applications.

## Related Documentation

⬆️ [Overview Section](../01-overview/01-project-goal.md) | ↔️ [Architecture Overview](01-architecture.md) · [Component Diagram](02-component-diagram.md) · [Data Flow](03-data-flow.md) | ⬇️ [Modules Section](../03-modules/)

---

## SDK Usage Model

**Confidence**: [H] High (90%)

ORAS Go is designed as a **library SDK** rather than a standalone service or application. It provides Go packages that applications import and use programmatically.

### Typical Integration Pattern

```go
package main

import (
    "context"
    "fmt"
    
    "oras.land/oras-go/v3"
    "oras.land/oras-go/v3/content/file"
    "oras.land/oras-go/v3/content/oci"
    "oras.land/oras-go/v3/registry/remote"
    "oras.land/oras-go/v3/registry/remote/auth"
)

func main() {
    ctx := context.Background()
    
    // 1. Create storage backend
    store, err := oci.New("./oci-layout")
    if err != nil {
        panic(err)
    }
    
    // 2. Create remote repository client
    repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
    if err != nil {
        panic(err)
    }
    
    // 3. Configure authentication
    repo.Client = &auth.Client{
        Credential: auth.StaticCredential("ghcr.io", auth.Credential{
            Username: "myuser",
            Password: "mytoken",
        }),
    }
    
    // 4. Perform operations
    desc, err := oras.Copy(ctx, repo, "v1.0", store, "latest", oras.CopyOptions{})
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Copied: %s@%s\n", desc.MediaType, desc.Digest)
}
```

---

## Integration Points

**Confidence**: [H] High (92%)

### 1. Command-Line Tools

ORAS Go serves as the foundational library for CLI tools:

**ORAS CLI**
```go
// github.com/oras-project/oras
import "oras.land/oras-go/v3"

func pushCommand(ctx context.Context, args []string) error {
    // Parse arguments
    ref := args[0]
    files := args[1:]
    
    // Create repository client
    repo, _ := remote.NewRepository(ref)
    
    // Pack and push artifact
    desc, _ := oras.Pack(ctx, repo, artifactType, files, opts)
    
    return oras.Copy(ctx, repo, desc.Digest.String(), repo, tag, opts)
}
```

**Notation CLI**
```go
// github.com/notaryproject/notation-go
import "oras.land/oras-go/v3/registry/remote"

func signArtifact(ctx context.Context, ref string, signature []byte) error {
    // Use ORAS for registry interaction
    repo, _ := remote.NewRepository(ref)
    
    // Push signature as referrer
    return repo.Pushes(ctx, signatureDesc, bytes.NewReader(signature))
}
```

### 2. CI/CD Pipelines

Integration into continuous integration systems:

**GitHub Actions**
```yaml
- name: Push artifact
  run: |
    # Custom Go action using ORAS Go SDK
    go run ./build-tool push \
      --artifact ./dist/app.tar.gz \
      --registry ghcr.io/${{ github.repository }}
```

**GitLab CI**
```go
// Custom GitLab CI job using ORAS
package main

import "oras.land/oras-go/v3"

func main() {
    // Read GitLab CI variables
    registry := os.Getenv("CI_REGISTRY")
    image := os.Getenv("CI_REGISTRY_IMAGE")
    
    // Push build artifacts
    oras.Copy(ctx, localStore, "latest", remoteRepo, tag, opts)
}
```

### 3. Container Runtimes

Potential integration with container runtimes:

**Containerd**
```go
// Example containerd plugin using ORAS
import (
    "oras.land/oras-go/v3/content/oci"
    "github.com/containerd/containerd/content"
)

type orasStore struct {
    store *oci.Store
}

func (s *orasStore) Info(ctx context.Context, dgst digest.Digest) (content.Info, error) {
    // Bridge ORAS to containerd content store
    desc := ocispec.Descriptor{Digest: dgst}
    return s.store.Fetch(ctx, desc)
}
```

### 4. Developer Tools

Custom artifact management applications:

**Supply Chain Security**
```go
// SBOM distribution tool
import "oras.land/oras-go/v3"

func distributeSBOM(ctx context.Context, image, sbom string) error {
    // Attach SBOM as referrer
    repo, _ := remote.NewRepository(image)
    
    sbomDesc, _ := oras.Pack(ctx, repo, "application/vnd.cyclonedx+json", 
        []string{sbom}, oras.PackOptions{
            Subject: &imageDesc,
        })
    
    return oras.Tag(ctx, repo, sbomDesc, "sbom")
}
```

**Reference**: [README.md](../../../README.md)

---

## Storage Backend Selection

**Confidence**: [H] High (93%)

### Decision Matrix

| Backend | Use Case | Persistence | Performance | Network | Thread Safety |
|---------|----------|-------------|-------------|---------|---------------|
| **Memory** | Testing, caching | None | High | N/A | Yes (locked) |
| **File** | Custom workflows | Yes (metadata only) | Medium | N/A | Yes (locked) |
| **OCI** | Standard layout | Yes (full) | Medium | N/A | Yes (RWMutex) |
| **Remote** | Registry operations | Remote | Low | Required | Yes (HTTP) |

### Memory Store

**Best For**: Short-lived operations, unit testing, temporary caching

```go
import "oras.land/oras-go/v3/content/memory"

store := memory.New()

// Use for testing
oras.Copy(ctx, remoteRepo, "v1.0", store, "", opts)
verifyContent(store) // Run tests
// Memory released when store goes out of scope
```

**Characteristics**:
- All content stored in memory
- No persistence across restarts
- High performance (no I/O)
- Limited by available RAM

**Reference**: [content/memory/memory.go](../../../content/memory/memory.go)

### File Store

**Best For**: Custom workflows with manual metadata management

```go
import "oras.land/oras-go/v3/content/file"

store, err := file.New("./artifacts")
defer store.Close()

// Files addressed by path, not digest
store.Add(ctx, "app.tar.gz", "application/gzip", "./build/app.tar.gz")

// Metadata stored in memory (not persisted)
```

**Characteristics**:
- Files stored by path (location-addressed)
- Metadata not persisted (memory only)
- Requires external metadata restoration
- Supports tar archives with reproducibility options

**Reference**: [content/file/file.go](../../../content/file/file.go)

### OCI Layout Store

**Best For**: Standards-compliant local artifact storage

```go
import "oras.land/oras-go/v3/content/oci"

store, err := oci.New("./oci-layout")
defer store.Close()

// Standard OCI layout
// - blobs/sha256/<digest>
// - index.json
// - oci-layout

oras.Copy(ctx, remoteRepo, "v1.0", store, "v1.0", opts)
```

**Characteristics**:
- OCI Image Layout specification compliant
- Full persistence (content + metadata)
- Auto-save index option
- Garbage collection support
- Concurrent access with locking

**Reference**: [content/oci/oci.go](../../../content/oci/oci.go)

### Remote Registry

**Best For**: Direct push/pull operations with remote registries

```go
import "oras.land/oras-go/v3/registry/remote"

repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
repo.Client = &auth.Client{...}

// Direct registry interaction
oras.Copy(ctx, repo, "v1.0", repo, "latest", opts)
```

**Characteristics**:
- HTTP-based Distribution API
- No local storage
- Network-dependent performance
- Built-in authentication and retry

**Reference**: [registry/remote/repository.go](../../../registry/remote/repository.go)

---

## Authentication Strategies

**Confidence**: [H] High (88%)

### 1. Static Credentials

Direct username/password or token:

```go
import "oras.land/oras-go/v3/registry/remote/auth"

repo.Client = &auth.Client{
    Credential: auth.StaticCredential(
        "ghcr.io",
        auth.Credential{
            Username: "myuser",
            Password: "ghp_personaltoken",
        },
    ),
}
```

### 2. Docker Config

Leverage existing Docker credentials:

```go
import (
    "oras.land/oras-go/v3/registry/remote/auth"
    "oras.land/oras-go/v3/registry/remote/credentials"
)

// Load from ~/.docker/config.json
store, err := credentials.NewStoreFromDockerConfig()
if err != nil {
    panic(err)
}

repo.Client = &auth.Client{
    Credential: credentials.Credential(store),
}
```

### 3. Dynamic Credential Functions

Custom credential resolution:

```go
func credentialFromVault(ctx context.Context, hostport string) (auth.Credential, error) {
    // Fetch from secret manager
    secret, err := vaultClient.Get(ctx, "registry/"+hostport)
    if err != nil {
        return auth.EmptyCredential, err
    }
    
    return auth.Credential{
        Username: secret["username"],
        Password: secret["password"],
    }, nil
}

repo.Client = &auth.Client{
    Credential: credentialFromVault,
}
```

### 4. OAuth2 / Bearer Tokens

Token-based authentication with automatic refresh:

```go
import "oras.land/oras-go/v3/registry/remote/auth"

// Token automatically fetched and cached
repo.Client = &auth.Client{
    Cache: auth.DefaultCache,
    Credential: func(ctx context.Context, hostport string) (auth.Credential, error) {
        // Return credentials for token exchange
        return auth.Credential{
            Username: "oauth2accesstoken",
            Password: os.Getenv("GITHUB_TOKEN"),
        }, nil
    },
}
```

**Features**:
- Automatic token caching by scope
- Challenge-response flow handling
- Thread-safe concurrent access
- Configurable cache expiration

**Reference**: [registry/remote/auth/client.go](../../../registry/remote/auth/client.go)

---

## Deployment Scenarios

**Confidence**: [H] High (90%)

### Scenario 1: Registry-to-Registry Copy

**Use Case**: Mirror artifacts between registries (e.g., disaster recovery, multi-region distribution)

```go
func mirrorArtifact(ctx context.Context, source, dest string) error {
    srcRepo, _ := remote.NewRepository(source)
    dstRepo, _ := remote.NewRepository(dest)
    
    // Configure auth for both
    srcRepo.Client = &auth.Client{Credential: sourceAuth}
    dstRepo.Client = &auth.Client{Credential: destAuth}
    
    // Copy with referrers
    opts := oras.ExtendedCopyOptions{
        ExtendedCopyGraphOptions: oras.ExtendedCopyGraphOptions{
            Depth: -1, // Copy all referrers
        },
    }
    
    return oras.ExtendedCopy(ctx, srcRepo, "v1.0", dstRepo, "v1.0", opts)
}
```

### Scenario 2: Registry-to-Local

**Use Case**: Download artifacts for local processing, offline access

```go
func downloadArtifact(ctx context.Context, ref, outputDir string) error {
    repo, _ := remote.NewRepository(ref)
    store, _ := oci.New(outputDir)
    defer store.Close()
    
    desc, err := oras.Copy(ctx, repo, "latest", store, "latest", oras.CopyOptions{
        CopyGraphOptions: oras.CopyGraphOptions{
            PostCopy: func(ctx context.Context, desc ocispec.Descriptor) error {
                fmt.Printf("Downloaded: %s (%d bytes)\n", desc.Digest, desc.Size)
                return nil
            },
        },
    })
    
    return err
}
```

### Scenario 3: Local Build and Push

**Use Case**: CI/CD build process creating and publishing artifacts

```go
func buildAndPublish(ctx context.Context, registry string) error {
    // Create file store for build artifacts
    buildStore, _ := file.New("./build")
    defer buildStore.Close()
    
    // Add build outputs
    buildStore.Add(ctx, "app", "application/octet-stream", "./build/app")
    buildStore.Add(ctx, "config.yaml", "application/yaml", "./build/config.yaml")
    
    // Create remote repository
    repo, _ := remote.NewRepository(registry + "/myapp")
    repo.Client = &auth.Client{Credential: ciAuth}
    
    // Pack into manifest
    layers := []ocispec.Descriptor{
        {Digest: digest.FromBytes(appBytes), Size: appSize},
        {Digest: digest.FromBytes(configBytes), Size: configSize},
    }
    
    desc, err := oras.PackManifest(ctx, repo, oras.PackManifestVersion1_1, 
        "application/vnd.myapp.v1", oras.PackManifestOptions{
            Layers: layers,
        })
    if err != nil {
        return err
    }
    
    // Tag with version
    return oras.TagBytes(ctx, repo, "v1.0", desc.MediaType, manifestJSON)
}
```

### Scenario 4: Artifact Signing Workflow

**Use Case**: Sign artifacts and attach signatures as referrers

```go
func signAndAttach(ctx context.Context, ref, signatureFile string) error {
    repo, _ := remote.NewRepository(ref)
    repo.Client = &auth.Client{Credential: registryAuth}
    
    // 1. Resolve artifact to sign
    artifactDesc, err := repo.Resolve(ctx, "v1.0")
    if err != nil {
        return err
    }
    
    // 2. Create signature (notation, cosign, etc.)
    signature, err := signArtifact(artifactDesc)
    if err != nil {
        return err
    }
    
    // 3. Push signature with subject reference
    signatureDesc, err := oras.PackManifest(ctx, repo, oras.PackManifestVersion1_1,
        "application/vnd.cncf.notary.signature", oras.PackManifestOptions{
            Subject: &artifactDesc,
            Layers: []ocispec.Descriptor{
                {
                    MediaType: "application/jose+json",
                    Digest:    digest.FromBytes(signature),
                    Size:      int64(len(signature)),
                },
            },
        })
    if err != nil {
        return err
    }
    
    fmt.Printf("Signature attached: %s\n", signatureDesc.Digest)
    return nil
}
```

### Scenario 5: SBOM Distribution

**Use Case**: Attach Software Bill of Materials to container images

```go
func attachSBOM(ctx context.Context, imageRef, sbomFile string) error {
    repo, _ := remote.NewRepository(imageRef)
    
    // Read SBOM
    sbomData, _ := os.ReadFile(sbomFile)
    
    // Resolve image
    imageDesc, _ := repo.Resolve(ctx, "latest")
    
    // Create SBOM layer
    sbomDesc := ocispec.Descriptor{
        MediaType: "application/vnd.cyclonedx+json",
        Digest:    digest.FromBytes(sbomData),
        Size:      int64(len(sbomData)),
    }
    
    // Push SBOM content
    repo.Push(ctx, sbomDesc, bytes.NewReader(sbomData))
    
    // Create manifest referencing image
    sbomManifest, _ := oras.PackManifest(ctx, repo, oras.PackManifestVersion1_1,
        "application/vnd.oci.image.sbom.v1+json", oras.PackManifestOptions{
            Subject: &imageDesc,
            Layers:  []ocispec.Descriptor{sbomDesc},
        })
    
    // Tag SBOM
    return oras.Tag(ctx, repo, sbomManifest, "latest.sbom")
}
```

**Reference**: Example tests in [example_test.go](../../../example_test.go)

---

## Concurrency Configuration

**Confidence**: [H] High (94%)

### Default Settings

```go
// Default concurrency: 3 parallel workers
opts := oras.CopyOptions{
    CopyGraphOptions: oras.CopyGraphOptions{
        Concurrency: 3, // default
    },
}
```

### Tuning Guidelines

**High-Bandwidth, Low-Latency Networks**:
```go
// Increase concurrency for faster copying
opts.Concurrency = 10

// Increase metadata cache
opts.MaxMetadataBytes = 32 * 1024 * 1024 // 32 MiB
```

**Rate-Limited Registries**:
```go
// Reduce concurrency to avoid rate limits
opts.Concurrency = 1

// Add retry policy
repo.Client.Client = &retry.Client{
    Policy: retry.Policy{
        MaxRetries: 5,
        Backoff:    retry.ExponentialBackoff,
    },
}
```

**Memory-Constrained Environments**:
```go
// Disable metadata caching
opts.MaxMetadataBytes = 0

// Reduce concurrency
opts.Concurrency = 2
```

**Reference**: [copy.go](../../../copy.go#L49-L123)

---

## Error Handling

**Confidence**: [M] Medium (82%)

### Error Types

```go
import "oras.land/oras-go/v3/errdef"

// Standard errors
var (
    ErrNotFound       = errdef.ErrNotFound
    ErrAlreadyExists  = errdef.ErrAlreadyExists
    ErrBadRequest     = errdef.ErrBadRequest
    ErrUnauthorized   = errdef.ErrUnauthorized
    ErrForbidden      = errdef.ErrForbidden
)

// Check error types
if errors.Is(err, errdef.ErrNotFound) {
    // Handle not found
}
```

### Copy-Specific Errors

```go
import "oras.land/oras-go/v3"

_, err := oras.Copy(ctx, src, srcRef, dst, dstRef, opts)
if err != nil {
    var copyErr *oras.CopyError
    if errors.As(err, &copyErr) {
        // Determine origin
        if copyErr.Origin == oras.CopyErrorOriginSource {
            fmt.Println("Error reading from source")
        } else {
            fmt.Println("Error writing to destination")
        }
    }
}
```

**Reference**: [copyerror.go](../../../copyerror.go)

---

## Best Practices

**Confidence**: [H] High (88%)

### 1. Always Use Context

```go
// Pass context for cancellation and timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

oras.Copy(ctx, src, srcRef, dst, dstRef, opts)
```

### 2. Close Resources

```go
store, err := oci.New("./oci-layout")
if err != nil {
    return err
}
defer store.Close() // Important!
```

### 3. Configure Retry Logic

```go
repo.Client.Client = &retry.Client{
    Policy: retry.Policy{
        MaxRetries:  3,
        MinDelay:    1 * time.Second,
        MaxDelay:    30 * time.Second,
        Backoff:     retry.ExponentialBackoff,
    },
}
```

### 4. Use Hooks for Observability

```go
opts := oras.CopyOptions{
    CopyGraphOptions: oras.CopyGraphOptions{
        PreCopy: func(ctx context.Context, desc ocispec.Descriptor) error {
            log.Printf("Copying: %s (%d bytes)", desc.Digest, desc.Size)
            return nil
        },
        PostCopy: func(ctx context.Context, desc ocispec.Descriptor) error {
            metrics.IncrementCopied(desc.Size)
            return nil
        },
    },
}
```

### 5. Handle Authentication Errors

```go
_, err := oras.Copy(ctx, repo, srcRef, dst, dstRef, opts)
if err != nil {
    if errors.Is(err, errdef.ErrUnauthorized) {
        // Prompt for credentials or refresh token
        return refreshAuthAndRetry()
    }
    return err
}
```

---

## Summary

ORAS Go's deployment model as a library SDK provides maximum flexibility for integration into diverse applications and workflows. The choice of storage backend, authentication strategy, and concurrency settings allows developers to optimize for their specific use cases, from CI/CD automation to supply chain security tools.

**Key Takeaways**:
- Library SDK for embedding into Go applications
- Multiple storage backends for different use cases
- Flexible authentication strategies
- Configurable concurrency and caching
- Rich error handling and observability hooks

---

**Document Status**
- **Confidence**: [H] High (90% overall)
- **Completeness**: Comprehensive deployment scenarios and integration patterns
- **Last Updated**: December 16, 2025
- **Source**: phase-2-findings.md (Sections 4, 5, 7, 9, 15)
