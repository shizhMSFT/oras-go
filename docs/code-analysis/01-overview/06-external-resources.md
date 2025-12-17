# External Resources

**Reading Time**: ~5 min  
**Analysis Request**: Analyze ORAS Go repository  
**Level**: Overview  
**Confidence**: H

## Context

ORAS Go supports multiple storage backends through a unified `Target` abstraction. Understanding these external resources is essential for choosing the right storage backend for your use case—whether that's production distribution via registries, local development with file systems, or testing with in-memory stores.

## Related Documentation

- **Previous**: [User Interface/SDK](05-user-interface.md) - SDK usage patterns
- **Next**: [Tech Stack](07-tech-stack.md) - Dependencies and tools
- **Related**: [Targets Documentation](../../Targets.md) - Target interface details

## Storage Backend Types

ORAS Go integrates with five types of external resources, each serving different use cases:

### 1. OCI-Compliant Registries (Remote Storage)

**Purpose**: Store and distribute OCI artifacts over HTTP/HTTPS

**Examples**:
- **Docker Hub** (`docker.io`)
- **GitHub Container Registry** (`ghcr.io`)
- **Azure Container Registry** (`<name>.azurecr.io`)
- **Google Artifact Registry** (`<region>-docker.pkg.dev`)
- **AWS ECR** (`<account>.dkr.ecr.<region>.amazonaws.com`)
- **Harbor** (self-hosted)
- **Distribution** (formerly Docker Registry v2)

**Configuration**:
```go
repo, err := remote.NewRepository("ghcr.io/myorg/myrepo")
repo.PlainHTTP = false  // Use HTTPS (default)
repo.Client = &auth.Client{...}  // Optional authentication
```

**Authentication Methods**:
- Basic auth (username/password)
- Token auth (OAuth2 bearer tokens)
- Refresh tokens
- Static credentials

**Data Operations**:
- Push/fetch manifests via Distribution API
- Push/fetch blobs (layers, configs)
- List tags
- List referrers (if supported)
- Mount blobs (cross-repository)

**Reference**: [registry/remote/repository.go](../../registry/remote/repository.go)

### 2. File System (Local Storage)

**Purpose**: Store artifacts as files on local disk

**Package**: `content/file`

**Configuration**:
```go
store, err := file.New("/path/to/artifacts")
defer store.Close()

// Optional settings
store.AllowPathTraversalOnWrite = false  // Security
store.DisableOverwrite = false           // Allow updates
```

**Storage Structure**:
```
/path/to/artifacts/
├── file1.txt  (annotated filename)
├── file2.tar.gz
└── .oras/
    └── sha256_<hash>  (blobs by digest)
```

**Data Operations**:
- Store blobs by digest
- Store files by annotation (`org.opencontainers.image.title`)
- Tar extraction support
- Path traversal protection

**Use Cases**:
- Local artifact storage
- Filesystem-based distribution
- Artifact extraction

**Reference**: [content/file/file.go](../../content/file/file.go)

### 3. OCI Image Layout (Standards-Compliant Local Storage)

**Purpose**: Store artifacts in OCI-compliant directory structure

**Package**: `content/oci`

**Configuration**:
```go
store, err := oci.New("/path/to/oci-layout")

// Optional settings
store.AutoSaveIndex = true  // Auto-save index.json (default: true)
```

**Directory Structure** (per OCI spec):
```
/path/to/oci-layout/
├── oci-layout          {"imageLayoutVersion": "1.0.0"}
├── index.json          (OCI Image Index with tagged manifests)
└── blobs/
    └── sha256/
        ├── abc123...   (manifest)
        ├── def456...   (config)
        └── 789012...   (layer)
```

**Data Operations**:
- Compliant with [OCI Image Layout Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/image-layout.md)
- Tag management via index.json
- Blob storage by algorithm/digest
- Automatic index persistence

**Use Cases**:
- Standards-compliant artifact storage
- Portable artifact archives
- Tool interoperability (skopeo, crane, etc.)

**Reference**: [content/oci/oci.go](../../content/oci/oci.go)

### 4. In-Memory Storage (Ephemeral Storage)

**Purpose**: Store artifacts in RAM (testing, caching)

**Package**: `content/memory`

**Configuration**:
```go
store := memory.New()
// No configuration needed - simple in-memory map
```

**Data Operations**:
- All operations in-memory (no persistence)
- Fast access
- Automatic cleanup on process exit

**Use Cases**:
- Unit testing
- Temporary caching
- Prototyping
- CI/CD pipelines (ephemeral artifacts)

**Reference**: [content/memory/memory.go](../../content/memory/memory.go)

### 5. Credential Stores (Authentication)

**Purpose**: Securely store registry credentials

**Package**: `registry/remote/credentials`

**Store Types**:

#### FileStore (Docker config.json)
```go
store, err := credentials.NewFileStore("~/.docker/config.json")
cred, err := store.Get(ctx, "ghcr.io")
```
- Location: `~/.docker/config.json` or custom path
- Format: Docker config JSON
- Shared with Docker CLI

#### NativeStore (OS Keychain)
```go
store := credentials.NewNativeStore("docker-credential-helper")
```
- **macOS**: Keychain
- **Windows**: Credential Manager
- **Linux**: Secret Service / pass

#### MemoryStore (In-Memory)
```go
store := credentials.NewMemoryStore()
store.Put(ctx, "ghcr.io", auth.Credential{Username: "user", Password: "pass"})
```
- Ephemeral (process lifetime)
- No persistence

**Reference**: [registry/remote/credentials/](../../registry/remote/credentials/)

## Resource Comparison

| Resource | Purpose | Configuration | Data Ops | Use Case |
|----------|---------|---------------|----------|----------|
| **Remote Registry** | OCI artifact distribution | `remote.NewRepository()` + auth | Push/fetch via HTTP | Production distribution |
| **File System** | Local file storage | `file.New(path)` | File-based storage | Local artifacts, extraction |
| **OCI Layout** | Standards-compliant local | `oci.New(path)` | OCI spec directory | Portable, interoperable |
| **Memory** | Ephemeral RAM storage | `memory.New()` | In-memory map | Testing, caching |
| **Credentials** | Auth storage | FileStore, NativeStore, etc. | Secure cred retrieval | Registry authentication |

## Choosing the Right Backend

**For Production Distribution**: Use remote registries (GitHub Container Registry, Docker Hub, etc.)

**For Local Development**: Use file system storage for quick iteration

**For CI/CD & Testing**: Use in-memory storage for fast, isolated tests

**For Tool Interoperability**: Use OCI Layout storage for compatibility with other OCI tools

**For Portable Archives**: Use OCI Layout for standards-compliant artifact bundles

## Next Steps

Continue with [Tech Stack](07-tech-stack.md) to learn about ORAS Go's dependencies and development tools.

## Citations

- [registry/remote/repository.go](../../registry/remote/repository.go) - Remote registry implementation
- [content/file/file.go](../../content/file/file.go) - File system storage
- [content/oci/oci.go](../../content/oci/oci.go) - OCI Layout implementation
- [content/memory/memory.go](../../content/memory/memory.go) - In-memory storage
- [registry/remote/credentials/](../../registry/remote/credentials/) - Credential management
- [OCI Image Layout Spec](https://github.com/opencontainers/image-spec/blob/v1.1.1/image-layout.md)
