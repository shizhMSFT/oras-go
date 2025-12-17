# Entry Points

**Reading Time**: ~4 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: [H] High

## Context

ORAS Go is a **library SDK** (not an application) consumed as a Go module. This document identifies the main APIs, packages, and interfaces that serve as entry points for developers using the SDK. Understanding these entry points is essential for effectively utilizing ORAS Go in your applications.

## Related Documentation

[⬆️ Project Goal](01-project-goal.md) | [↔️ File Usage](03-file-usage.md) | [⬇️ User Interface/SDK](05-user-interface.md)

## SDK Structure

**Import Path**: `oras.land/oras-go/v3`

**Package Type**: Library (no CLI or GUI)

The SDK is organized into a top-level package with several specialized sub-packages:
- `oras` - High-level copy and pack operations
- `registry` - Registry abstractions and interfaces
- `registry/remote` - Remote registry client implementation
- `content` - Content storage abstractions
- `content/*` - Storage implementations (memory, file, oci)

## Main APIs (Top-Level Package)

### Core Copy Operations

Located in [copy.go](../../copy.go):

#### Copy
```go
func Copy(ctx, src, srcRef, dst, dstRef, opts) (ocispec.Descriptor, error)
```
**Line**: [L132](../../copy.go#L132)

**Purpose**: Copy artifact with tag/reference resolution

**Use Case**: Pull from registry, push to file system, or copy between any targets

#### CopyGraph
```go
func CopyGraph(ctx, src, dst, root, opts) (ocispec.Descriptor, error)
```
**Line**: [L176](../../copy.go#L176)

**Purpose**: Copy artifact graph by descriptor (no tag resolution)

**Use Case**: Copy when you already have the manifest descriptor

### Extended Copy Operations

Located in [extendedcopy.go](../../extendedcopy.go):

#### ExtendedCopy
```go
func ExtendedCopy(ctx, src, srcRef, dst, dstRef, opts) (ocispec.Descriptor, error)
```
**Line**: [L77](../../extendedcopy.go#L77)

**Purpose**: Copy artifact along with referrers (signatures, SBOMs) and predecessors

**Use Case**: Copy complete artifact graph including all attached metadata

#### ExtendedCopyGraph
```go
func ExtendedCopyGraph(ctx, src, dst, node, opts) (ocispec.Descriptor, error)
```
**Line**: [L109](../../extendedcopy.go#L109)

**Purpose**: Extended graph copy by descriptor

**Use Case**: Copy graph with referrers when you have the descriptor

### Artifact Packing

Located in [pack.go](../../pack.go):

#### Pack
```go
func Pack(ctx, pusher, artifactType, blobs, opts) (ocispec.Descriptor, error)
```
**Line**: [L190](../../pack.go#L190)

**Purpose**: Pack blobs into a manifest

**Use Case**: Create new artifacts from existing blobs

#### PackManifest
```go
func PackManifest(ctx, pusher, packManifestVersion, artifactType, opts) (ocispec.Descriptor, error)
```
**Line**: [L139](../../pack.go#L139)

**Purpose**: Pack manifest with specific version (v1.0 or v1.1)

**Use Case**: Create artifacts compatible with specific OCI versions

### Target Interfaces

Located in [target.go](../../target.go):

These interfaces define the abstraction for storage backends:

| Interface | Lines | Description |
|-----------|-------|-------------|
| `Target` | [L20-24](../../target.go#L20-L24) | CAS with generic tags (Storage + TagResolver) |
| `GraphTarget` | [L27-31](../../target.go#L27-L31) | CAS with tags + predecessor finding |
| `ReadOnlyTarget` | [L34-37](../../target.go#L34-L37) | Read-only Target |
| `ReadOnlyGraphTarget` | [L40-43](../../target.go#L40-L43) | Read-only GraphTarget |

## Key Sub-Packages

### 1. Registry Package

**Purpose**: Primary interfaces for registry operations

**Location**: [registry/](../../registry/)

#### Registry Interface
**File**: [registry/registry.go](../../registry/registry.go#L21-L41)

**Purpose**: Repository collection interface for multi-repo operations

#### Repository Interface
**File**: [registry/repository.go](../../registry/repository.go#L30-L48)

**Purpose**: Union of ORAS target + blob/manifest CAS + referrers API

### 2. Remote Registry Package

**Purpose**: Remote registry client implementation

**Location**: [registry/remote/](../../registry/remote/)

#### remote.Repository
**File**: [registry/remote/repository.go](../../registry/remote/repository.go#L99-L148)

**Purpose**: HTTP client to communicate with remote OCI registries

**Factory Function**:
```go
repo, err := remote.NewRepository(ref)
```

**Key Features**:
- HTTP/HTTPS support
- Authentication integration
- Retry logic
- Referrers API support

### 3. Authentication Package

**Purpose**: Authentication support for remote registries

**Location**: [registry/remote/auth/](../../registry/remote/auth/)

**Key Types**:
- **`auth.Client`**: HTTP client with authentication capabilities
- **`auth.StaticCredential`**: Static username/password authentication
- **`auth.Cache`**: Token caching for performance

### 4. Content Package

**Purpose**: Core content storage abstractions

**Location**: [content/](../../content/)

**Key Interfaces**:

| Interface | File | Purpose |
|-----------|------|---------|
| `Storage` | [content/storage.go](../../content/storage.go#L43-L46) | CAS interface (Fetch, Push, Exists) |
| `GraphStorage` | [content/graph.go](../../content/graph.go#L38-L42) | Storage + Predecessors navigation |
| `TagResolver` | [content/resolver.go](../../content/resolver.go#L41-L44) | Tag/Resolve interface |

### 5. Content Store Implementations

These packages provide concrete implementations of the content storage interfaces:

#### Memory Store
**Package**: `content/memory`  
**File**: [content/memory/memory.go](../../content/memory/memory.go#L32-L37)  
**Purpose**: In-memory storage for testing and caching

#### File Store
**Package**: `content/file`  
**File**: [content/file/file.go](../../content/file/file.go)  
**Purpose**: File system-based storage for local artifacts

#### OCI Store
**Package**: `content/oci`  
**File**: [content/oci/oci.go](../../content/oci/oci.go#L49-L52)  
**Purpose**: OCI Image Layout compliant storage

## Usage Pattern Example

From [example_test.go](../../example_test.go#L33-L65):

```go
// 1. Create target (file store, OCI store, or remote repo)
fs, _ := file.New("/tmp/artifacts")
repo, _ := remote.NewRepository("registry.example.com/myrepo")
repo.Client = &auth.Client{...} // Optional authentication

// 2. Copy between targets
manifestDesc, _ := oras.Copy(ctx, repo, "latest", fs, "latest", opts)
```

## Entry Point Decision Tree

```
Need to:
├─ Copy artifacts?
│  ├─ With tag resolution → oras.Copy()
│  ├─ With descriptor only → oras.CopyGraph()
│  └─ With referrers/signatures → oras.ExtendedCopy()
│
├─ Create artifacts?
│  ├─ From blobs → oras.Pack()
│  └─ Specific version → oras.PackManifest()
│
├─ Connect to registry?
│  ├─ Remote registry → remote.NewRepository()
│  ├─ File system → file.New()
│  ├─ OCI layout → oci.New()
│  └─ In-memory → memory.New()
│
└─ Implement custom storage?
   └─ Implement content.Storage interface
```

## Next Steps

- Read [File Usage Analysis](03-file-usage.md) to understand which files implement these interfaces
- Read [User Interface/SDK](05-user-interface.md) for detailed usage patterns and examples
- Explore [Prerequisites](04-prerequisites.md) to understand required background knowledge

## Further Reading

- [copy.go](../../copy.go) - Implementation details of copy operations
- [target.go](../../target.go) - Complete target interface definitions
- [example_test.go](../../example_test.go) - Practical usage examples
