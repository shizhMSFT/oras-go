# Project Goal

**Reading Time**: ~3 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: [H] High

## Context

ORAS (OCI Registry as Storage) Go is a Go SDK that enables developers to work with OCI artifacts and registries programmatically. It serves as both a library for artifact management and a client for interacting with OCI-compliant registries.

## Related Documentation

⬇️ [Entry Points](02-entry-points.md) | ↔️ [Tech Stack](07-tech-stack.md) | ⬇️ [User Interface/SDK](05-user-interface.md)

## Purpose

ORAS Go is a **Go SDK for OCI artifact management** and serves as an **OCI registry client**. It provides unified APIs for push/pull operations across:

- OCI-compliant registries (Docker Hub, GHCR, ACR, etc.)
- File systems (local storage)
- In-memory stores (testing, caching)

**Primary Sources**:
- [README.md](../../README.md)
- [AGENTS.md](../../AGENTS.md)

## Standards Compliance

ORAS Go fully implements and supports:

### OCI Image Format Specification v1.1.1
Defines the schema for container images including:
- Manifests (artifact metadata)
- Image indexes (multi-platform support)
- Filesystem layers (content storage)
- Configuration objects

### OCI Distribution Specification v1.1.1  
Defines the API protocol for content distribution including:
- Push/pull endpoints
- Blob management
- Manifest operations
- Referrers API

**Reference**: [AGENTS.md](../../AGENTS.md#L5-L7)

## Primary Use Cases

### 1. Artifact Management
Push, pull, and manage OCI artifacts programmatically across different storage backends.

### 2. Cross-Storage Operations
Copy artifacts seamlessly between:
- Remote registries ↔ Local file systems
- Registry ↔ Registry
- File system ↔ OCI layout
- Any storage ↔ Memory

### 3. Registry Client
Interact with OCI-compliant and Docker-compliant registries using HTTP APIs with full authentication support.

### 4. Content-Addressable Storage
Implement reliable content storage using cryptographic digests for:
- Content integrity verification
- Automatic deduplication
- Secure content retrieval

**Reference**: [README.md](../../README.md), [docs/tutorial/quickstart.md](../../docs/tutorial/quickstart.md)

## Critical Design Principles

### 1. Content-Addressable Storage (CAS)
All content is addressed by cryptographic digests (descriptors), enabling:
- **Reliability**: Content integrity verified via digest matching
- **Deduplication**: Identical content stored once, referenced by same digest  
- **Immutability**: Content cannot change without changing its digest

### 2. Graph-Based Model
Unlike other OCI clients, ORAS models every element of an artifact as **nodes in a Directed Acyclic Graph (DAG)**:
- **Nodes**: Descriptors (manifests, configs, layers, blobs)
- **Edges**: References between descriptors
- **Traversal**: Efficient graph operations (Successors, Predecessors)

This enables sophisticated artifact relationships including referrers, signatures, and SBOMs.

**Learn more**: [ORAS Graphs Documentation](https://oras.land/docs/client_libraries/overview#graphs)

### 3. Copy-Based Operations
ORAS models data movement as **copy operations** rather than separate push/pull:
- Unified API: `Copy(src, dst)` works regardless of source/destination type
- Consistent behavior across storage backends
- Simplified developer experience

**Learn more**: [Unified Experience](https://oras.land/docs/client_libraries/overview#unified-experience)

**Reference**: [AGENTS.md](../../AGENTS.md#L12-L16)

## Project Structure

```
ORAS Go SDK
├── Core APIs (oras package)
│   ├── Copy operations (Copy, CopyGraph, ExtendedCopy)
│   ├── Pack operations (Pack, PackManifest)
│   └── Content utilities (Fetch, Tag, Resolve)
│
├── Content Storage (content package)
│   ├── Interfaces (Storage, GraphStorage, TagResolver)
│   ├── Memory store (testing/caching)
│   ├── File store (local artifacts)
│   └── OCI layout store (standards-compliant)
│
├── Registry Client (registry package)
│   ├── Abstractions (Registry, Repository)
│   ├── Remote client (HTTP/HTTPS)
│   ├── Authentication (OAuth2, basic auth)
│   └── Credentials (Docker config, OS keychain)
│
└── Internal Utilities
    ├── CAS implementations
    ├── Graph algorithms
    ├── Concurrency utilities
    └── Platform selection
```

## Version Information

- **Current Development**: v3 (main branch) - Breaking changes expected
- **Stable Production**: v2 branch - Recommended for production use
- **Maintenance**: v1 branch - Security fixes only

**Reference**: [README.md](../../README.md#L44-L94)

## Next Steps

- **Understand the SDK**: Read [User Interface/SDK](05-user-interface.md)
- **Learn Prerequisites**: Review [Prerequisites](04-prerequisites.md) for required domain knowledge
- **Get Started**: Follow [Getting Started Guide](08-getting-started.md)
- **Explore Entry Points**: See [Entry Points](02-entry-points.md) for main APIs

## Citations

[^1]: All confidence ratings based on direct code analysis and official documentation from [README.md](../../README.md) and [AGENTS.md](../../AGENTS.md)
