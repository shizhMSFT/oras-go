# Architecture Overview

**Reading Time**: ~12 min  
**Analysis Phase**: 2 - Architecture  
**Level**: Architecture  
**Confidence**: [H] High (95%)

## Context

This document provides a comprehensive overview of the ORAS Go library architecture. ORAS Go implements a layered, interface-driven architecture following clean architecture principles with clear separation between core abstractions, storage implementations, and network protocols. The system models OCI artifacts as Directed Acyclic Graphs (DAGs) in Content-Addressable Storage (CAS).

## Related Documentation

⬆️ [Overview Section](../01-overview/01-project-goal.md) | ↔️ [Component Diagram](02-component-diagram.md) | ⬇️ [Modules Section](../03-modules/)

---

## Architectural Style

**Pattern**: Layered Architecture + Hexagonal Architecture (Ports & Adapters)

The architecture consists of four primary layers that enable separation of concerns and extensibility:

```
┌─────────────────────────────────────────────────────────┐
│         Application Layer (High-Level APIs)              │
│  Copy, Pack, ExtendedCopy, Tag, TagBytes               │
└─────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│           Core Abstraction Layer (Interfaces)            │
│  Target, Storage, GraphStorage, Repository              │
└─────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│        Storage Implementation Layer (Adapters)           │
│  Memory, File, OCI Layout, Remote Registry              │
└─────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│         Infrastructure Layer (Support Services)          │
│  Auth, Retry, Cache, Sync, Status Tracking             │
└─────────────────────────────────────────────────────────┘
```

**Reference Files**:
- [target.go](../../../target.go#L1-L43) - Core abstractions
- [content/storage.go](../../../content/storage.go#L1-L82) - Storage interfaces
- [copy.go](../../../copy.go#L1-L534) - High-level copy operations
- [pack.go](../../../pack.go#L1-L449) - High-level pack operations

---

## Layer 1: Application Layer

**Purpose**: Provides high-level, user-facing APIs for common operations

### Key Components

#### Copy Operations
Located in [copy.go](../../../copy.go):
- `Copy()` - Copy artifact with tagging
- `CopyGraph()` - Copy DAG without tagging
- `ExtendedCopy()` - Copy with referrers/predecessors

#### Pack Operations
Located in [pack.go](../../../pack.go):
- `Pack()` - Create and push manifests (deprecated)
- `PackManifest()` - Create OCI manifests (v1.0, v1.1)
- `packManifestV1_0()`, `packManifestV1_1()` - Version-specific packing

#### Tag Operations
Located in [content.go](../../../content.go):
- `Tag()` - Tag single reference
- `TagN()` - Tag multiple references
- `TagBytes()` - Tag raw content

### Design Patterns

- **Facade Pattern**: Simplifies complex graph operations
- **Options Pattern**: Flexible configuration via `*Options` structs
- **Template Method**: Common flow with customizable hooks (PreCopy, OnCopySkipped)

---

## Layer 2: Core Abstraction Layer

**Purpose**: Defines contracts for storage and artifact management

### Storage Hierarchy

Located in [content/storage.go](../../../content/storage.go):

```go
Fetcher → ReadOnlyStorage → Storage
                           ↓
                      GraphStorage
```

### Target Hierarchy

Located in [target.go](../../../target.go):

```go
ReadOnlyTarget → Target
               ↓
          GraphTarget
```

### Registry Hierarchy

Located in [registry/repository.go](../../../registry/repository.go):

```go
Repository → BlobStore, ManifestStore
```

### Interface Responsibilities

| Interface | Responsibility | Key Methods |
|-----------|---------------|-------------|
| `Fetcher` | Content retrieval | `Fetch` |
| `Pusher` | Content storage | `Push` |
| `Storage` | Full CAS operations | `Fetch`, `Push`, `Exists` |
| `GraphStorage` | Storage + predecessor finding | All Storage methods + graph ops |
| `Target` | Storage + tag resolution | All Storage methods + `Resolve` |
| `Repository` | Blob + manifest stores with referrers | Union of all storage capabilities |

**Reference Files**:
- [target.go](../../../target.go#L19-L44)
- [content/storage.go](../../../content/storage.go#L24-L54)
- [content/graph.go](../../../content/graph.go#L33-L43)
- [registry/repository.go](../../../registry/repository.go#L44-L103)

---

## Layer 3: Storage Implementation Layer

**Purpose**: Concrete implementations of storage abstractions

### 1. Memory Store

**Implementation**: [content/memory/memory.go](../../../content/memory/memory.go)

**Features**:
- In-memory CAS using `internal/cas/memory.go`
- Graph indexing via `internal/graph/memory.go`
- Tag resolution via `internal/resolver/memory.go`

**Use Case**: Testing, temporary storage, caching

### 2. File Store

**Implementation**: [content/file/file.go](../../../content/file/file.go)

**Features**:
- File system-based CAS
- Location-addressed by file paths
- Virtual CAS with memory-stored metadata
- Tar handling with reproducibility options

**Use Case**: Local artifact storage, OCI layout alternative

### 3. OCI Layout Store

**Implementation**: [content/oci/oci.go](../../../content/oci/oci.go)

**Features**:
- OCI Image Layout spec compliant
- Directory structure: `blobs/sha256/`, `index.json`
- Auto-save index and garbage collection options

**Use Case**: Standard OCI artifact storage

### 4. Remote Registry

**Implementation**: [registry/remote/repository.go](../../../registry/remote/repository.go)

**Features**:
- HTTP client to remote OCI registries
- Distribution API implementation
- Referrers API support
- Manifest and blob operations

**Use Case**: Push/pull from remote registries

---

## Layer 4: Infrastructure Layer

**Purpose**: Supporting services and utilities

### 1. Authentication

**Location**: [registry/remote/auth/](../../../registry/remote/auth/)

**Features**:
- OAuth2, basic auth, bearer token
- Credential management
- Challenge handling
- Static and dynamic credential functions

### 2. Caching

**Implementation**: [internal/cas/proxy.go](../../../internal/cas/proxy.go)

**Features**:
- Proxy pattern for transparent caching
- Memory-limited cache for metadata
- Fetch-time caching with size limits

### 3. Concurrency Control

**Location**: [internal/syncutil/](../../../internal/syncutil/)

**Features**:
- Semaphore-based concurrency limiting
- Limited regions for goroutine coordination
- Merge operations for parallel work

### 4. Status Tracking

**Implementation**: [internal/status/tracker.go](../../../internal/status/tracker.go)

**Features**:
- Content operation tracking
- Duplicate work prevention
- Descriptor-based status management

### 5. Retry Logic

**Location**: [registry/remote/retry/](../../../registry/remote/retry/)

**Features**:
- Configurable retry policies
- Backoff strategies
- Network error handling

---

## Key Architectural Patterns

### Interface Segregation Pattern

The architecture heavily uses small, focused interfaces:

```go
// Content operations
type Fetcher interface {
    Fetch(ctx, desc) (ReadCloser, error)
}

type Pusher interface {
    Push(ctx, desc, reader) error
}

// Combined into larger interfaces
type Storage interface {
    ReadOnlyStorage
    Pusher
}
```

**Benefits**:
- Fine-grained capability composition
- Easy mocking for testing
- Flexible implementation requirements

**Reference**: [content/storage.go](../../../content/storage.go#L24-L54)

### Proxy Pattern (Caching)

```go
type Proxy struct {
    ReadOnlyStorage
    Cache       content.Storage
    StopCaching bool
}
```

**Purpose**: Transparent metadata caching for repeated fetches  
**Implementation**: [internal/cas/proxy.go](../../../internal/cas/proxy.go#L28-L100)

### Options Pattern

```go
type CopyOptions struct {
    CopyGraphOptions
    MapRoot func(context.Context, ...) (Descriptor, error)
}

type CopyGraphOptions struct {
    Concurrency      int
    MaxMetadataBytes int64
    PreCopy          func(...) error
    PostCopy         func(...) error
    OnCopySkipped    func(...) error
}
```

**Benefits**:
- Backward compatible API evolution
- Optional parameters with sensible defaults
- Extensible configuration

**Reference**: [copy.go](../../../copy.go#L49-L123)

### DAG Traversal Pattern

```go
// Recursive parallel traversal with status tracking
fn := func(ctx, region, desc) error {
    done, committed := tracker.TryCommit(desc)
    if !committed { return nil } // Skip if being processed
    defer close(done)
    
    successors := FindSuccessors(ctx, proxy, desc)
    syncutil.Go(ctx, limiter, fn, successors...)
    
    return copyNode(ctx, src, dst, desc, opts)
}
```

**Features**:
- Depth-first with concurrent execution
- Deduplication via status tracker
- Semaphore-based concurrency control

**Reference**: [copy.go](../../../copy.go#L214-L270)

---

## Concurrency Architecture

### Concurrency Control Mechanisms

**Components**:

1. **Semaphore Limiter** (`golang.org/x/sync/semaphore`)
   - Default concurrency: 3 (configurable)
   - Prevents resource exhaustion
   - Fair scheduling

2. **Status Tracker** ([internal/status/tracker.go](../../../internal/status/tracker.go))
   - Descriptor-based work tracking
   - `sync.Map` for concurrent access
   - Channel-based completion signaling

3. **Limited Regions** ([internal/syncutil/](../../../internal/syncutil/))
   - Scoped concurrency control
   - Deadlock prevention
   - Hierarchical task spawning

**Reference**: [copy.go](../../../copy.go#L188-L211)

### Thread Safety Guarantees

| Storage Implementation | Thread Safety Approach |
|------------------------|------------------------|
| **Memory Store** | Thread-safe via internal locking |
| **File Store** | Concurrent reads, synchronized writes |
| **OCI Store** | RWMutex for concurrent ops, exclusive delete |
| **Remote Registry** | HTTP client handles concurrency |

**Reference**: [content/oci/oci.go](../../../content/oci/oci.go#L67-L77)

---

## Design Principles

### SOLID Principles

- **S** - Single Responsibility: Each interface has one clear purpose
- **O** - Open/Closed: Extensible via interfaces, closed for modification
- **L** - Liskov Substitution: All storage implementations interchangeable
- **I** - Interface Segregation: Small, focused interfaces (Fetcher, Pusher)
- **D** - Dependency Inversion: Depends on abstractions, not concrete types

### Clean Architecture

- Clear layer separation
- Dependency rule: Inner layers don't depend on outer layers
- Business logic (copy, pack) independent of I/O details

### Go Idioms

- Interface-based design
- Error values (not exceptions)
- Context-based cancellation
- Reader/Writer streaming
- Channel-based synchronization

---

## Architectural Strengths

✅ **Clean Separation of Concerns**: Layers are well-defined with minimal coupling  
✅ **Interface-Driven Design**: Enables testing, mocking, and extensibility  
✅ **Multiple Storage Backends**: Unified API across memory, file, OCI, and remote  
✅ **Concurrency Safety**: Built-in synchronization and deduplication  
✅ **Streaming I/O**: Memory-efficient handling of large artifacts  
✅ **Extensibility**: Hook points for custom behavior  
✅ **Minimal Dependencies**: Only essential external packages

---

## Architectural Considerations

⚠️ **Metadata-only Memory Storage**: File store doesn't persist full CAS (needs external restoration)  
⚠️ **Concurrency Default**: Default of 3 may be low for high-throughput scenarios  
⚠️ **Error Context**: Could benefit from structured error types with more context  
⚠️ **Metrics/Observability**: No built-in metrics or tracing (left to users)

---

## Key Architecture Files

| File | Purpose | Criticality |
|------|---------|-------------|
| [target.go](../../../target.go) | Core abstractions | Critical |
| [copy.go](../../../copy.go) | Copy operations | Critical |
| [pack.go](../../../pack.go) | Manifest creation | Critical |
| [content/storage.go](../../../content/storage.go) | Storage interfaces | Critical |
| [content/graph.go](../../../content/graph.go) | Graph operations | Critical |
| [registry/repository.go](../../../registry/repository.go) | Repository interface | Critical |
| [registry/remote/repository.go](../../../registry/remote/repository.go) | Remote client | Critical |
| [content/oci/oci.go](../../../content/oci/oci.go) | OCI layout | High |
| [content/file/file.go](../../../content/file/file.go) | File storage | High |
| [content/memory/memory.go](../../../content/memory/memory.go) | Memory storage | Medium |
| [internal/cas/proxy.go](../../../internal/cas/proxy.go) | Caching proxy | High |
| [internal/status/tracker.go](../../../internal/status/tracker.go) | Status tracking | High |

---

## Next Steps

- [Component Diagram](02-component-diagram.md) - Detailed component interactions and data flows
- [Module Analysis](../03-modules/) - Deep dive into individual modules

---

## Citations

**Source Document**: [phase-2-findings.md](../.workspace/phase-2-findings.md) - Section 1  
**Analysis Date**: December 16, 2025  
**Confidence Level**: 95%
