# User Interface/SDK

**Reading Time**: ~6 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: [H] High

## Context

ORAS Go provides a comprehensive SDK for working with OCI artifacts through a clean, interface-driven design. This document explores the SDK architecture, core concepts, key packages, and common usage patterns that developers will encounter when using the library.

## Related Documentation

[⬆️ Prerequisites](04-prerequisites.md) | [↔️ Entry Points](02-entry-points.md) | [⬇️ Tech Stack](07-tech-stack.md)

## SDK Type

**Library SDK** - Imported as a Go module. No CLI or GUI provided.

**Import Path**:
```go
import "oras.land/oras-go/v3"
```

**Module Declaration**:
```go
module github.com/oras-project/oras-go/v3
```

## Core Concepts

### 1. Target Abstraction

**Definition**: Unified interface for content storage across different backends

**Philosophy**: Write code once, run on any storage backend (registry, file system, memory)

**Interface Hierarchy**:
```
Storage          (CAS operations: Fetch, Push, Exists)
  └─ GraphStorage    (+ Predecessors)

TagResolver      (Tag <-> Digest resolution)

Target = Storage + TagResolver
GraphTarget = GraphStorage + TagResolver
```

**Available Implementations**:

| Implementation | Package | Use Case |
|----------------|---------|----------|
| `memory.Store` | `content/memory` | Testing, ephemeral caching |
| `file.Store` | `content/file` | Local file system storage |
| `oci.Store` | `content/oci` | OCI Image Layout (standards-compliant) |
| `remote.Repository` | `registry/remote` | Remote OCI registries (production) |

**Reference**: [target.go](../../target.go), [docs/Targets.md](../../docs/Targets.md)

### 2. Copy Operations

**Philosophy**: Unified copy model replaces traditional push/pull operations

**Why Copy?**
- Conceptually simpler: `Copy(source, destination)`
- Works between any targets (registry ↔ file, memory ↔ registry, etc.)
- Consistent API regardless of direction
- Natural support for registry-to-registry copies

**Operations Available**:

| Function | Purpose | Location |
|----------|---------|----------|
| `Copy()` | Copy with tag resolution | [copy.go](../../copy.go#L132) |
| `CopyGraph()` | Copy by descriptor (no tags) | [copy.go](../../copy.go#L176) |
| `ExtendedCopy()` | Copy with referrers | [extendedcopy.go](../../extendedcopy.go#L77) |
| `ExtendedCopyGraph()` | Extended copy by descriptor | [extendedcopy.go](../../extendedcopy.go#L109) |

**Reference**: [copy.go](../../copy.go), [extendedcopy.go](../../extendedcopy.go)

### 3. Descriptor-Driven Operations

**Philosophy**: All operations use OCI descriptors for content addressing

**Key Type**: `ocispec.Descriptor`

**Source**: `github.com/opencontainers/image-spec/specs-go/v1`

**Structure**:
```go
type Descriptor struct {
    MediaType   string            // Content type
    Digest      digest.Digest     // sha256:...
    Size        int64             // bytes
    URLs        []string          // optional
    Annotations map[string]string // optional
    Data        []byte            // optional
    Platform    *Platform         // optional
    ArtifactType string           // optional
    Subject     *Descriptor       // optional (referrers)
}
```

**Benefits**:
- Type-safe content references
- Built-in integrity verification
- Immutable content addressing
- Rich metadata support

**Reference**: [content/descriptor.go](../../content/descriptor.go)

## Key Packages

### Package: `oras` (Top-Level)

**Purpose**: High-level operations for artifact management

**Import**: `import "oras.land/oras-go/v3"`

**Key Functions**:

| Category | Functions |
|----------|-----------|
| **Copy** | `Copy()`, `CopyGraph()`, `ExtendedCopy()`, `ExtendedCopyGraph()` |
| **Pack** | `Pack()`, `PackManifest()` |
| **Content** | `Fetch()`, `FetchBytes()` |
| **Tagging** | `Tag()`, `TagN()` |
| **Resolution** | `Resolve()` |

**Usage Pattern**:
```go
// Copy from remote registry to local file
manifestDesc, err := oras.Copy(ctx, remoteRepo, "v1.0.0", fileStore, "v1.0.0", opts)
```

### Package: `registry`

**Purpose**: Registry abstractions and interfaces

**Import**: `import "oras.land/oras-go/v3/registry"`

**Key Types**:

#### Registry Interface
**Location**: [registry/registry.go](../../registry/registry.go#L21-L41)

**Purpose**: Repository collection management

**Methods**:
- `Repository(ctx, name) (Repository, error)` - Access specific repository
- `Repositories(ctx) ([]string, error)` - List repositories

#### Repository Interface
**Location**: [registry/repository.go](../../registry/repository.go#L30-L48)

**Purpose**: Union of Target + BlobStore + ManifestStore + ReferrerLister

**Extends**: `Target` interface

**Methods**: Inherits all Target methods plus referrers API

#### Reference Type
**Location**: [registry/reference.go](../../registry/reference.go)

**Purpose**: Parse and validate registry references

**Example**:
```go
ref, err := registry.ParseReference("ghcr.io/myorg/myrepo:v1.0.0")
// ref.Registry = "ghcr.io"
// ref.Repository = "myorg/myrepo"
// ref.Reference = "v1.0.0"
```

### Package: `registry/remote`

**Purpose**: Remote registry HTTP client implementation

**Import**: `import "oras.land/oras-go/v3/registry/remote"`

**Key Types**:

#### remote.Repository
**Location**: [registry/remote/repository.go](../../registry/remote/repository.go#L99-L148)

**Purpose**: Full-featured HTTP client for OCI registries

**Creation**:
```go
repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
```

**Configuration**:
```go
repo.PlainHTTP = false           // Use HTTPS (default)
repo.SkipTLSVerify = false       // Verify TLS certificates
repo.Client = &auth.Client{...}  // Optional authentication
repo.MaxMetadataBytes = 4 * 1024 * 1024  // Metadata size limit
```

**Features**:
- HTTP/HTTPS support
- Authentication integration
- Automatic retry logic
- Referrers API support
- Blob mounting
- Chunked uploads

#### remote.Registry
**Location**: [registry/remote/registry.go](../../registry/remote/registry.go)

**Purpose**: Multi-repository registry client

**Example**:
```go
reg, err := remote.NewRegistry("ghcr.io")
repo, err := reg.Repository(ctx, "myorg/myrepo")
```

### Package: `registry/remote/auth`

**Purpose**: Authentication for remote registries

**Import**: `import "oras.land/oras-go/v3/registry/remote/auth"`

**Key Types**:

#### auth.Client
**Purpose**: HTTP client with automatic authentication

**Features**:
- Basic auth (username/password)
- Token auth (OAuth2 bearer tokens)
- Refresh token support
- Token caching
- Challenge-response handling

**Example**:
```go
repo.Client = &auth.Client{
    Client: retry.DefaultClient,
    Cache:  auth.NewCache(),
    Credential: auth.StaticCredential("ghcr.io", auth.Credential{
        Username: "user",
        Password: "token",
    }),
}
```

#### auth.Credential
**Purpose**: Username/password credentials

**Structure**:
```go
type Credential struct {
    Username string
    Password string
    RefreshToken string  // optional
    AccessToken  string  // optional
}
```

#### auth.Cache
**Purpose**: Token caching for performance

**Usage**:
```go
cache := auth.NewCache()
client := &auth.Client{Cache: cache, ...}
```

### Package: `registry/remote/credentials`

**Purpose**: Credential storage management

**Import**: `import "oras.land/oras-go/v3/registry/remote/credentials"`

**Store Types**:

| Store | Purpose | Location |
|-------|---------|----------|
| `FileStore` | Docker config.json format | `~/.docker/config.json` |
| `NativeStore` | OS keychain integration | macOS Keychain, Windows Credential Manager |
| `MemoryStore` | Ephemeral in-memory storage | Process lifetime only |

**Example**:
```go
store, err := credentials.NewFileStore("~/.docker/config.json")
cred, err := store.Get(ctx, "ghcr.io")

repo.Client = &auth.Client{
    Credential: credentials.Credential(store),
}
```

### Package: `content`

**Purpose**: Content storage abstractions

**Import**: `import "oras.land/oras-go/v3/content"`

**Key Interfaces**:

#### Storage Interface
**Location**: [content/storage.go](../../content/storage.go#L43-L46)

**Methods**:
```go
type Storage interface {
    Fetch(ctx, target Descriptor) (io.ReadCloser, error)
    Push(ctx, expected Descriptor, content io.Reader) error
    Exists(ctx, target Descriptor) (bool, error)
}
```

#### GraphStorage Interface
**Location**: [content/graph.go](../../content/graph.go#L38-L42)

**Extends**: `Storage`

**Methods**:
```go
type GraphStorage interface {
    Storage
    Predecessors(ctx, node Descriptor) ([]Descriptor, error)
}
```

#### TagResolver Interface
**Location**: [content/resolver.go](../../content/resolver.go#L41-L44)

**Methods**:
```go
type TagResolver interface {
    Resolve(ctx, reference string) (Descriptor, error)
    Tag(ctx, desc Descriptor, reference string) error
}
```

**Key Functions**:

| Function | Purpose |
|----------|---------|
| `Successors()` | Get child nodes (config, layers) |
| `FetchAll()` | Fetch and verify content |
| `NewDescriptorFromBytes()` | Create descriptor from content |

### Package: `content/memory`

**Purpose**: In-memory storage implementation

**Import**: `import "oras.land/oras-go/v3/content/memory"`

**Key Type**: `memory.Store`

**Example**:
```go
store := memory.New()
desc := content.NewDescriptorFromBytes("text/plain", []byte("hello"))
err := store.Push(ctx, desc, bytes.NewReader([]byte("hello")))
```

**Use Cases**:
- Unit testing
- Temporary caching
- Prototype development
- CI/CD ephemeral artifacts

### Package: `content/file`

**Purpose**: File system storage implementation

**Import**: `import "oras.land/oras-go/v3/content/file"`

**Key Type**: `file.Store`

**Example**:
```go
store, err := file.New("/tmp/artifacts")
defer store.Close()

// Optional configuration
store.AllowPathTraversalOnWrite = false
store.DisableOverwrite = false
```

**Storage Structure**:
```
/tmp/artifacts/
├── myfile.txt              (named by annotation)
├── archive.tar.gz
└── .oras/
    └── sha256_abc123...    (blob by digest)
```

**Features**:
- Annotation-based file naming (`org.opencontainers.image.title`)
- Tar extraction support
- Path traversal protection
- Automatic directory creation

### Package: `content/oci`

**Purpose**: OCI Image Layout storage implementation

**Import**: `import "oras.land/oras-go/v3/content/oci"`

**Key Type**: `oci.Store`

**Example**:
```go
store, err := oci.New("/path/to/oci-layout")

// Optional configuration
store.AutoSaveIndex = true  // Auto-save index.json
```

**OCI Layout Structure**:
```
/path/to/oci-layout/
├── oci-layout             {"imageLayoutVersion": "1.0.0"}
├── index.json             (OCI Image Index)
└── blobs/
    └── sha256/
        ├── abc123...      (manifest)
        ├── def456...      (config)
        └── 789012...      (layer)
```

**Features**:
- Full OCI Image Layout Spec compliance
- Tag management via index.json
- Automatic index persistence
- Tool interoperability (skopeo, crane, etc.)

## Usage Patterns

### Pattern 1: Pull from Registry to File System

**Scenario**: Download artifact from remote registry to local storage

```go
import (
    "context"
    "oras.land/oras-go/v3"
    "oras.land/oras-go/v3/content/file"
    "oras.land/oras-go/v3/registry/remote"
    "oras.land/oras-go/v3/registry/remote/auth"
    "oras.land/oras-go/v3/registry/remote/retry"
)

// Create file store
fs, err := file.New("/tmp/artifacts")
if err != nil {
    panic(err)
}
defer fs.Close()

// Connect to remote repository
repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
if err != nil {
    panic(err)
}

// Configure authentication
repo.Client = &auth.Client{
    Client: retry.DefaultClient,
    Cache:  auth.NewCache(),
    Credential: auth.StaticCredential("ghcr.io", auth.Credential{
        Username: "myuser",
        Password: "mytoken",
    }),
}

// Copy artifact
ctx := context.Background()
desc, err := oras.Copy(ctx, repo, "v1.0.0", fs, "v1.0.0", oras.DefaultCopyOptions)
if err != nil {
    panic(err)
}

fmt.Printf("Copied %s\n", desc.Digest)
```

**Reference**: [example_test.go](../../example_test.go#L33-L65)

### Pattern 2: Pack and Push Artifact

**Scenario**: Create new artifact and push to registry

```go
import (
    "bytes"
    "context"
    "oras.land/oras-go/v3"
    "oras.land/oras-go/v3/content"
    "oras.land/oras-go/v3/registry/remote"
)

ctx := context.Background()

// Connect to repository
repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
if err != nil {
    panic(err)
}

// Prepare layer content
layerData := []byte("Hello, ORAS!")
layerDesc := content.NewDescriptorFromBytes("application/octet-stream", layerData)

// Push layer blob
err = repo.Push(ctx, layerDesc, bytes.NewReader(layerData))
if err != nil {
    panic(err)
}

// Pack manifest
artifactType := "application/vnd.example.artifact"
manifestDesc, err := oras.Pack(ctx, repo, artifactType, 
    []ocispec.Descriptor{layerDesc}, 
    oras.PackOptions{
        PackOptions: oras.PackManifestOptions{
            Layers: []ocispec.Descriptor{layerDesc},
        },
    })
if err != nil {
    panic(err)
}

// Tag manifest
err = repo.Tag(ctx, manifestDesc, "v1.0.0")
if err != nil {
    panic(err)
}

fmt.Printf("Pushed artifact: %s\n", manifestDesc.Digest)
```

**Reference**: [example_pack_test.go](../../example_pack_test.go)

### Pattern 3: Copy with Referrers

**Scenario**: Copy artifact along with all attached signatures and SBOMs

```go
import (
    "context"
    "oras.land/oras-go/v3"
    "oras.land/oras-go/v3/registry/remote"
)

ctx := context.Background()

srcRepo, _ := remote.NewRepository("ghcr.io/source/repo")
dstRepo, _ := remote.NewRepository("ghcr.io/dest/repo")

// Copy artifact with all referrers
desc, err := oras.ExtendedCopy(ctx, srcRepo, "latest", dstRepo, "latest",
    oras.ExtendedCopyOptions{
        ExtendedCopyGraphOptions: oras.ExtendedCopyGraphOptions{
            Depth: -1,  // Copy all referrers (signatures, SBOMs, etc.)
        },
    })
if err != nil {
    panic(err)
}

fmt.Printf("Copied artifact with referrers: %s\n", desc.Digest)
```

**Reference**: [extendedcopy.go](../../extendedcopy.go)

### Pattern 4: List Referrers

**Scenario**: Find all artifacts referencing a manifest (signatures, SBOMs)

```go
import (
    "context"
    "fmt"
    "oras.land/oras-go/v3/registry"
    "oras.land/oras-go/v3/registry/remote"
)

ctx := context.Background()
repo, _ := remote.NewRepository("ghcr.io/myorg/myrepo")

// Resolve manifest descriptor
manifestDesc, err := repo.Resolve(ctx, "v1.0.0")
if err != nil {
    panic(err)
}

// List all referrers
referrers, err := registry.Referrers(ctx, repo, manifestDesc, "")
if err != nil {
    panic(err)
}

// Display referrers
for _, ref := range referrers {
    fmt.Printf("Referrer: %s\n", ref.Digest)
    fmt.Printf("  Type: %s\n", ref.ArtifactType)
    if ref.Annotations != nil {
        fmt.Printf("  Title: %s\n", ref.Annotations["org.opencontainers.image.title"])
    }
}
```

**Reference**: [registry/repository.go](../../registry/repository.go)

### Pattern 5: Traverse Artifact Graph

**Scenario**: Explore the structure of an artifact (manifest → config + layers)

```go
import (
    "context"
    "fmt"
    "oras.land/oras-go/v3/content"
    "oras.land/oras-go/v3/registry/remote"
)

ctx := context.Background()
repo, _ := remote.NewRepository("ghcr.io/myorg/myrepo")

// Resolve manifest
manifestDesc, _ := repo.Resolve(ctx, "v1.0.0")

// Get all children (successors)
successors, err := content.Successors(ctx, repo, manifestDesc)
if err != nil {
    panic(err)
}

fmt.Printf("Manifest: %s\n", manifestDesc.Digest)
for _, child := range successors {
    fmt.Printf("  Child: %s (%s)\n", child.Digest, child.MediaType)
}

// For GraphStorage, can also find predecessors
if graphStore, ok := repo.(content.GraphStorage); ok {
    predecessors, err := graphStore.Predecessors(ctx, manifestDesc)
    if err != nil {
        panic(err)
    }
    
    for _, parent := range predecessors {
        fmt.Printf("  Parent: %s\n", parent.Digest)
    }
}
```

**Reference**: [content/graph.go](../../content/graph.go)

## Design Patterns Summary

| Pattern | Use Case | Key APIs |
|---------|----------|----------|
| **Target Abstraction** | Storage-agnostic operations | `Target`, `Storage`, `TagResolver` |
| **Copy Operations** | Unified data movement | `Copy()`, `CopyGraph()` |
| **Descriptor-Driven** | Type-safe content refs | `ocispec.Descriptor` |
| **Graph Traversal** | Explore relationships | `Successors()`, `Predecessors()` |
| **Referrers** | Attach metadata | `ExtendedCopy()`, `Referrers()` |

## Next Steps

- Explore [Tech Stack](07-tech-stack.md) for dependency details
- Review [Getting Started](08-getting-started.md) to set up your development environment
- Study [example_test.go](../../example_test.go) for more usage examples

## Further Reading

- [docs/Targets.md](../../docs/Targets.md) - Deep dive on Target abstractions
- [docs/tutorial/quickstart.md](../../docs/tutorial/quickstart.md) - Step-by-step tutorial
- [ORAS Project Docs](https://oras.land/) - Broader ecosystem documentation
