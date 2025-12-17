# Prerequisites & Learning Path

**Reading Time**: ~6 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: [H] High

## Context

Understanding ORAS Go requires knowledge of several foundational concepts from container technology and distributed systems. This document outlines the required background knowledge and provides a structured learning path for developers new to OCI artifacts and content-addressable storage.

## Related Documentation

[⬆️ File Usage](03-file-usage.md) | [↔️ User Interface/SDK](05-user-interface.md) | [⬇️ Getting Started](08-getting-started.md)

## Required Background Knowledge

### 1. OCI (Open Container Initiative) Concepts

**What is OCI?**

A set of open standards for container formats and runtimes, maintained by the Linux Foundation. OCI ensures containers work consistently across different tools and platforms.

**Key OCI Specifications**:

#### OCI Image Format Specification v1.1.1
Defines how container images are structured:
- **Manifests**: Artifact metadata (config + layers)
- **Image Indexes**: Multi-platform support
- **Filesystem Layers**: Content storage
- **Configuration Objects**: Runtime configuration

#### OCI Distribution Specification v1.1.1
Defines how to distribute container images via registries:
- Push/pull endpoints
- Blob management
- Manifest operations
- Referrers API

**Learn More**:
- [OCI Image Spec v1.1.1](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md)
- [OCI Distribution Spec v1.1.1](https://github.com/opencontainers/distribution-spec/blob/v1.1.1/spec.md)

### 2. Content-Addressable Storage (CAS)

**Definition**: 

A storage system where content is accessed via a cryptographic hash of the content itself, rather than by location or name.

**Key Concepts**:

| Concept | Description |
|---------|-------------|
| **Digest** | Cryptographic hash (e.g., `sha256:abc123...`) that uniquely identifies content |
| **Immutability** | Content cannot change without changing its digest |
| **Deduplication** | Identical content stored once, referenced by same digest |
| **Verification** | Content integrity guaranteed by matching digest |

**Benefits**:
- ✅ Automatic content integrity verification
- ✅ Efficient storage (no duplicates)
- ✅ Reliable content retrieval
- ✅ Tamper detection

**ORAS Implementation**: 

All blobs and manifests in ORAS are stored and retrieved via descriptors containing cryptographic digests. This ensures data integrity throughout copy operations.

**Reference**: [docs/Modeling-Artifacts.md](../../docs/Modeling-Artifacts.md#L3-L5)

### 3. Directed Acyclic Graph (DAG) Model

**Definition**: 

A graph structure where nodes are connected by directed edges with no cycles. This is a fundamental data structure for modeling artifact relationships.

**In ORAS Context**:

| Element | Description |
|---------|-------------|
| **Nodes** | OCI descriptors (manifests, configs, layers, blobs) |
| **Edges** | References between descriptors |
| **Root** | Manifest (entry point to artifact) |
| **Leaves** | Blobs (actual content) |

**Example Simple Graph**:
```
Manifest (root)
├─→ Config blob
├─→ Layer blob 1
└─→ Layer blob 2
```

**Example with Referrers (Subject)**:
```
Signature Manifest
└─→ Subject: Artifact Manifest
    ├─→ Config blob
    └─→ Layer blob
```

**Why DAG?**
- Models artifact relationships naturally
- Supports referrers and signatures
- Enables efficient graph traversal
- Handles shared nodes (same config referenced multiple times)
- Prevents circular dependencies

**ORAS Uniqueness**: Unlike other OCI clients, ORAS models **every element** of an artifact as nodes in a DAG, enabling sophisticated graph operations.

**Reference**: [docs/Modeling-Artifacts.md](../../docs/Modeling-Artifacts.md), [ORAS Graphs](https://oras.land/docs/client_libraries/overview#graphs)

### 4. OCI Descriptors

**Definition**: 

A metadata object that references content in CAS. Descriptors are the fundamental building blocks of OCI artifacts.

**Required Fields** ([docs/Targets.md](../../docs/Targets.md#L8-L11)):

| Field | Type | Description |
|-------|------|-------------|
| `mediaType` | string | MIME type of content (e.g., `application/vnd.oci.image.manifest.v1+json`) |
| `digest` | string | Content digest (e.g., `sha256:5b0bcabd...`) |
| `size` | int64 | Content size in bytes |

**Optional Fields**:

| Field | Purpose |
|-------|---------|
| `annotations` | Key-value metadata |
| `platform` | OS/architecture information |
| `urls` | External URLs for content |
| `data` | Embedded data for small content |
| `artifactType` | Artifact type identifier |
| `subject` | Parent artifact reference (referrers) |

**Example Descriptor**:
```json
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "size": 7682,
  "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
  "annotations": {
    "org.opencontainers.image.created": "2025-12-16T10:00:00Z"
  }
}
```

**Reference**: [OCI Descriptor Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/descriptor.md)

### 5. OCI Manifests & Image Index

#### Image Manifest

Root descriptor for an artifact, containing references to all content:

**Key Fields**:
- `config`: Configuration descriptor
- `layers`: Array of layer descriptors
- `subject`: (Optional) descriptor of parent artifact
- `artifactType`: Type of artifact
- `annotations`: Metadata

**Purpose**: Describes the structure and contents of an artifact

#### Image Index

Collection of manifests (e.g., multi-platform images):

**Key Fields**:
- `manifests`: Array of manifest descriptors
- `annotations`: Metadata

**Purpose**: Enables multi-platform support (Linux/Windows, amd64/arm64, etc.)

**Reference**: [docs/Modeling-Artifacts.md](../../docs/Modeling-Artifacts.md#L7-L150)

### 6. OCI Distribution API

**Purpose**: HTTP-based API for interacting with OCI registries

**Key Endpoints**:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/v2/<name>/manifests/<reference>` | GET | Fetch manifest by tag/digest |
| `/v2/<name>/manifests/<reference>` | PUT | Push manifest |
| `/v2/<name>/blobs/<digest>` | GET | Fetch blob |
| `/v2/<name>/blobs/uploads/` | POST | Initiate blob upload |
| `/v2/<name>/blobs/uploads/<uuid>` | PUT | Complete blob upload |
| `/v2/<name>/referrers/<digest>` | GET | List referrers |

**Authentication**: Supports OAuth2 bearer tokens and basic auth

**Reference**: [OCI Distribution Spec](https://distribution.github.io/distribution/spec/api/)

### 7. Referrers

**Definition**: 

Artifacts that reference another artifact via the `subject` field. This enables attaching metadata to existing artifacts without modifying them.

**Use Cases**:

| Use Case | Artifact Type | Description |
|----------|---------------|-------------|
| Signatures | `application/vnd.oci.image.signature` | Cryptographic signatures |
| SBOMs | `application/spdx+json` | Software Bill of Materials |
| Scan Results | `application/vnd.security.scan` | Vulnerability scan reports |
| Attestations | `application/vnd.in-toto+json` | Build attestations |

**API**: 

`GET /v2/<name>/referrers/<digest>` returns a list of artifacts referencing the specified digest.

**Example**:
```
Image Manifest (sha256:abc...)
├─ Signature 1 (subject: sha256:abc...)
├─ Signature 2 (subject: sha256:abc...)
└─ SBOM (subject: sha256:abc...)
```

**Reference**: [docs/Modeling-Artifacts.md](../../docs/Modeling-Artifacts.md#L60-L145)

## Terminology Glossary

| Term | Definition |
|------|------------|
| **Artifact** | Collection of files/data packaged with an OCI manifest |
| **Blob** | Raw content stored in CAS (layer, config, etc.) |
| **CAS** | Content-Addressable Storage - storage accessed by content hash |
| **DAG** | Directed Acyclic Graph - graph structure with no cycles |
| **Descriptor** | Metadata referencing content (mediaType, digest, size) |
| **Digest** | Cryptographic hash of content (e.g., `sha256:...`) |
| **Index** | OCI Image Index - collection of manifests for multi-platform support |
| **Manifest** | OCI Image Manifest - artifact metadata containing config + layers |
| **MediaType** | MIME type identifying content format |
| **OCI** | Open Container Initiative - container standards organization |
| **Predecessor** | Node pointing TO current node (parent in graph) |
| **Referrer** | Artifact referencing another via subject field |
| **Registry** | Remote OCI-compliant storage service (e.g., Docker Hub) |
| **Subject** | Artifact being referenced by a referrer |
| **Successor** | Node pointed BY current node (child in graph) |
| **Tag** | Human-readable name for manifest (e.g., "latest", "v1.0.0") |
| **Target** | Storage abstraction in ORAS (registry, file system, memory) |

## Structured Learning Path

### Stage 1: Fundamentals (Required)

**Time Investment**: ~1 hour

**Objectives**: Understand ORAS architecture and design philosophy

**Reading Order**:
1. [README.md](../../README.md) - Project overview and quick examples
2. [AGENTS.md](../../AGENTS.md) - Design principles and philosophy
3. [docs/Modeling-Artifacts.md](../../docs/Modeling-Artifacts.md) - DAG/CAS model deep dive
4. [docs/Targets.md](../../docs/Targets.md) - Storage interface abstractions

**Outcome**: Clear mental model of how ORAS works and why it's architected this way

### Stage 2: Hands-On Tutorial (Recommended)

**Time Investment**: ~30 minutes

**Objectives**: Practical experience with SDK usage

**Activities**:
1. Follow [docs/tutorial/quickstart.md](../../docs/tutorial/quickstart.md)
2. Run example code from [example_test.go](../../example_test.go)
3. Experiment with:
   - Creating a memory store
   - Packing a simple artifact
   - Copying between stores

**Outcome**: Ability to write basic ORAS Go programs

### Stage 3: Deep Dive (Optional)

**Time Investment**: ~2-3 hours

**Objectives**: Advanced understanding of implementation details

**Reading Order**:
1. [copy.go](../../copy.go) - Copy operations internals and graph traversal
2. [pack.go](../../pack.go) - Manifest packing and artifact creation
3. [content/graph.go](../../content/graph.go) - Graph traversal algorithms
4. [registry/remote/repository.go](../../registry/remote/repository.go) - HTTP client implementation

**Outcome**: Ability to contribute to ORAS or implement advanced features

### Stage 4: External Resources (Optional)

**Time Investment**: ~4-6 hours

**Objectives**: Deep understanding of OCI standards

**Resources**:
1. [OCI Image Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md) - Complete image format specification
2. [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec/blob/v1.1.1/spec.md) - Registry API specification
3. [ORAS Project Website](https://oras.land/) - Broader ORAS ecosystem
4. [OCI Artifacts](https://github.com/opencontainers/artifacts) - Artifact use cases

**Outcome**: Expert-level understanding of container artifact standards

## Prerequisites Checklist

Before diving into ORAS Go development:

- [ ] Understand basic Go programming (interfaces, goroutines)
- [ ] Familiar with HTTP/REST APIs
- [ ] Understand cryptographic hashes (SHA-256)
- [ ] Know what a container image is (Docker/Podman experience helpful)
- [ ] Understand JSON data structures
- [ ] Familiar with command-line tools
- [ ] (Optional) Experience with Docker registries

## Common Misconceptions

### Misconception 1: ORAS is just for containers
**Reality**: ORAS works with **any** OCI artifact, not just container images. Use cases include Helm charts, WebAssembly modules, ML models, etc.

### Misconception 2: Push/Pull are the primary operations
**Reality**: ORAS uses a **unified copy model**. There's no separate "push" or "pull" - it's all `Copy(source, destination)`.

### Misconception 3: You need a registry to use ORAS
**Reality**: ORAS works with file systems, OCI layouts, and memory stores - no registry required for many use cases.

### Misconception 4: ORAS modifies existing artifacts
**Reality**: All content is **immutable** in CAS. New artifacts are created rather than modifying existing ones.

## Next Steps

- Begin with [Stage 1: Fundamentals](#stage-1-fundamentals-required) of the learning path
- Review [User Interface/SDK](05-user-interface.md) for detailed API documentation
- Try the [Getting Started](08-getting-started.md) guide to set up your environment

## Further Reading

- [Content-Addressable Storage Explained](https://en.wikipedia.org/wiki/Content-addressable_storage)
- [Graph Theory Basics](https://en.wikipedia.org/wiki/Directed_acyclic_graph)
- [CNCF ORAS Project](https://oras.land/)
- [OCI Specification Index](https://github.com/opencontainers/image-spec)
