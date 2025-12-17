# References

**Module**: Maintenance Analysis  
**Component**: References  
**Last Updated**: December 16, 2025

[⬆️ Back to Quality](../05-quality/01-testing.md)

---

## Overview

This document provides a comprehensive reference guide for ORAS Go development, including external specifications, project documentation, API references, community resources, and development tools.

---

## External Specifications

### OCI Standards

#### Core Specifications

| Specification | Version | Purpose | ORAS Go Usage |
|--------------|---------|---------|---------------|
| **OCI Image Spec** | v1.1.1 | Image format definition | Core descriptor format, manifest structure |
| **OCI Distribution Spec** | v1.1.1 | Registry API specification | Registry client implementation |
| **OCI Image Layout** | v1.1.1 | Local storage format | OCI store implementation |

**Primary Reference**: [OCI Image Spec v1.1.1](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md)

#### Detailed Specifications

**OCI Descriptor** (v1.1.1)
- **URL**: [GitHub](https://github.com/opencontainers/image-spec/blob/v1.1.1/descriptor.md)
- **Usage**: Content addressing, digest verification
- **Implementation**: [content/descriptor.go](../../../content/descriptor.go)

**OCI Manifest** (v1.1.1)
- **URL**: [GitHub](https://github.com/opencontainers/image-spec/blob/v1.1.1/manifest.md)
- **Usage**: Artifact packaging, layer references
- **Implementation**: [pack.go](../../../pack.go), manifest utilities

**OCI Artifact** (v1.1.0-rc2, deprecated)
- **URL**: [GitHub](https://github.com/opencontainers/image-spec/blob/v1.1.0-rc2/artifact.md)
- **Status**: Deprecated in favor of OCI Image Manifest
- **Implementation**: [internal/spec/artifact.go](../../../internal/spec/artifact.go)
- **Note**: Maintained for backward compatibility

### Docker Standards

| Specification | Purpose | ORAS Go Usage |
|--------------|---------|---------------|
| **Docker Registry HTTP API V2** | Registry protocol | [registry/remote/](../../../registry/remote/) implementation |
| **Docker Credential Helpers** | Secure credential storage | [registry/remote/credentials/](../../../registry/remote/credentials/) |
| **Docker Token Authentication** | Bearer token auth flow | [registry/remote/auth/](../../../registry/remote/auth/) |

**Primary Reference**: [Docker Registry API](https://distribution.github.io/distribution/spec/api/)

#### Docker Documentation Links

- **Registry API**: [distribution.github.io](https://distribution.github.io/distribution/spec/api/)
- **Credential Helpers**: [docs.docker.com](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers)
- **Token Auth**: [distribution](https://distribution.github.io/distribution/spec/auth/token/)

### Go Standards

| Resource | Purpose | Relevance to ORAS Go |
|----------|---------|----------------------|
| **Go Security Policy** | Version support policy | 2 latest Go versions supported |
| **Go Modules** | Dependency management | [go.mod](../../../go.mod) structure |
| **Go Internal Packages** | Package visibility | [internal/](../../../internal/) organization |

**Links**:
- [Go Security Policy](https://github.com/golang/go/security/policy)
- [Go Modules](https://go.dev/ref/mod)
- [Go Internal Packages](https://go.dev/doc/go1.4#internalpackages)

---

## Project Documentation

### Core Documentation

| Document | Purpose | Location |
|----------|---------|----------|
| **README.md** | Project overview, quick start | [README.md](../../../README.md) |
| **MIGRATION_GUIDE.md** | v1→v2→v3 migration | [MIGRATION_GUIDE.md](../../../MIGRATION_GUIDE.md) |
| **CODE_OF_CONDUCT.md** | Community guidelines | [CODE_OF_CONDUCT.md](../../../CODE_OF_CONDUCT.md) |
| **SECURITY.md** | Security reporting | [SECURITY.md](../../../SECURITY.md) |
| **OWNERS.md** | Project maintainers | [OWNERS.md](../../../OWNERS.md) |
| **CODEOWNERS** | Code review assignments | [CODEOWNERS](../../../CODEOWNERS) |
| **AGENTS.md** | AI agent guidelines | [AGENTS.md](../../../AGENTS.md) |

### Conceptual Guides

**Location**: [docs/](../../../docs/)

| Guide | Topics Covered | Key Concepts |
|-------|---------------|--------------|
| **Modeling Artifacts** | DAG representation, CAS concepts | Directed Acyclic Graphs, Content-Addressable Storage |
| **Targets** | Storage interfaces | Storage, Target, GraphTarget |
| **Quickstart** | Step-by-step tutorial | Basic copy, pack operations (v2) |

**Links**:
- [Modeling Artifacts](../../../docs/Modeling-Artifacts.md)
- [Targets](../../../docs/Targets.md)
- [Quickstart Tutorial](../../../docs/tutorial/quickstart.md)

### Code Analysis Documentation

**Location**: [docs/code-analysis/](../code-analysis/)

Comprehensive analysis across 6 phases:

| Phase | Focus Area | Key Documents |
|-------|-----------|--------------|
| **Phase 1** | Overview | Project goal, structure, dependencies |
| **Phase 2** | Architecture | Design patterns, core components |
| **Phase 3** | Modules | Package-by-package deep dive |
| **Phase 4** | Implementation | Code patterns, error handling |
| **Phase 5** | Quality | Testing, security, performance |
| **Phase 6** | Maintenance | Unused files, quick wins, roadmap |

**Navigation**:
- [Phase 1: Overview](../01-overview/)
- [Phase 2: Architecture](../02-architecture/)
- [Phase 3: Modules](../03-modules/)
- [Phase 4: Implementation](../04-implementation/)
- [Phase 5: Quality](../05-quality/)
- [Phase 6: Maintenance](./) (this section)

---

## API Documentation

### Package Documentation

**Primary Source**: [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2)

#### Core Packages (v2)

| Package | Purpose | Documentation | Example Count |
|---------|---------|--------------|---------------|
| **oras.land/oras-go/v2** | Core operations | [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2) | 15+ |
| **registry** | Repository interface | [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2/registry) | 8+ |
| **registry/remote** | Remote registry client | [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2/registry/remote) | 12+ |
| **content** | Storage interfaces | [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2/content) | 6+ |

#### Registry Packages

| Package | Purpose | Key Types |
|---------|---------|-----------|
| **registry/remote/auth** | Authentication | `Client`, `Credential` |
| **registry/remote/credentials** | Credential management | `Store`, `FileStore` |
| **registry/remote/retry** | Retry policies | `Policy`, `DefaultPolicy` |

**Links**:
- [registry/remote/auth](https://pkg.go.dev/oras.land/oras-go/v2/registry/remote/auth)
- [registry/remote/credentials](https://pkg.go.dev/oras.land/oras-go/v2/registry/remote/credentials)

#### Content Storage Packages

| Package | Purpose | Implementation |
|---------|---------|----------------|
| **content/file** | File-based storage | Filesystem directories |
| **content/oci** | OCI layout storage | OCI Image Layout Spec |
| **content/memory** | In-memory storage | Testing, temporary operations |

**Links**:
- [content/file](https://pkg.go.dev/oras.land/oras-go/v2/content/file)
- [content/oci](https://pkg.go.dev/oras.land/oras-go/v2/content/oci)
- [content/memory](https://pkg.go.dev/oras.land/oras-go/v2/content/memory)

### v3 Documentation

**Status**: Not yet published (v3 still in development)  
**Expected**: Available after v3.0.0 GA release  
**Location**: Will be at `pkg.go.dev/oras.land/oras-go/v3`

---

## Community Resources

### Communication Channels

| Channel | Purpose | Response Time | Link |
|---------|---------|---------------|------|
| **GitHub Issues** | Bug reports, feature requests | 1-3 business days | [Issues](https://github.com/oras-project/oras-go/issues) |
| **GitHub Discussions** | Q&A, design discussions | 1-5 business days | [Discussions](https://github.com/oras-project/oras-go/discussions) |
| **CNCF Slack #oras** | Real-time support | Best-effort (hours) | [Slack](https://cloud-native.slack.com/archives/CJ1KHJM5Z) |

**Getting Help**:
1. **Search existing issues/discussions** first
2. **Provide minimal reproducible example**
3. **Include version information** (Go version, ORAS Go version, OS)
4. **Tag appropriately** (bug, enhancement, question)

### Project Resources

| Resource | Purpose | URL |
|----------|---------|-----|
| **ORAS Project Website** | Ecosystem overview | [oras.land](https://oras.land/) |
| **ORAS CLI** | Command-line tool | [GitHub](https://github.com/oras-project/oras) |
| **GitHub Actions** | CI/CD pipelines | [Actions](https://github.com/oras-project/oras-go/actions) |
| **Codecov** | Test coverage reports | [codecov.io](https://codecov.io/gh/oras-project/oras-go) |

### Related Projects

#### Ecosystem Integrations

| Project | Relationship | URL |
|---------|-------------|-----|
| **containerd** | Uses ORAS for artifact management | [containerd.io](https://containerd.io/) |
| **Notary v2** | Supply chain security with OCI | [GitHub](https://github.com/notaryproject) |
| **Helm** | OCI registry support via ORAS | [helm.sh](https://helm.sh/) |
| **Artifact Hub** | Artifact discovery and distribution | [artifacthub.io](https://artifacthub.io/) |

**Use Cases**:
- **containerd**: Content store management, image distribution
- **Notary v2**: Signature storage and verification
- **Helm**: Chart storage in OCI registries
- **Artifact Hub**: Artifact cataloging and search

---

## Development Resources

### Testing & Quality

| Tool/Resource | Purpose | Configuration |
|--------------|---------|---------------|
| **Go Test** | Unit testing | Standard `go test` |
| **Race Detector** | Concurrency bug detection | `go test -race` |
| **Coverage** | Test coverage tracking | 90.2% average |
| **Benchmarks** | Performance testing | See [Quick Wins](02-quick-wins.md#quick-win-2) |

**Testing Reference**: [Phase 5: Testing](../05-quality/01-testing.md)

#### Running Tests

```bash
# All tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -cover ./...

# Specific package
go test ./registry/remote/...

# Verbose output
go test -v ./...
```

### Build & Deploy

| Tool | Purpose | Location |
|------|---------|----------|
| **Makefile** | Build automation | [Makefile](../../../Makefile) |
| **go.mod** | Dependency management | [go.mod](../../../go.mod) |
| **GitHub Actions** | CI/CD | `.github/workflows/` |

#### Build Commands

```bash
# Build all packages
go build ./...

# Run tests
make test

# Generate coverage report
make coverage

# Run linters
make lint
```

### Dependencies

**From** [go.mod](../../../go.mod):

```go
require (
    github.com/opencontainers/go-digest v1.0.0
    github.com/opencontainers/image-spec v1.1.1
    golang.org/x/sync v0.19.0
)
```

#### Key Dependencies

**opencontainers/go-digest**
- **Purpose**: Content digest computation (SHA256, etc.)
- **Usage**: Descriptor digest calculation, content verification
- **Version**: v1.0.0 (stable)

**opencontainers/image-spec**
- **Purpose**: OCI types and specifications
- **Usage**: Descriptor, manifest, config types
- **Version**: v1.1.1 (latest stable)

**golang.org/x/sync**
- **Purpose**: Advanced concurrency primitives
- **Usage**: Semaphore (concurrency limiting), errgroup
- **Version**: v0.19.0

**Dependency Philosophy**: Minimal external dependencies for security, maintainability, and reduced attack surface

---

## Specification Evolution

### Tracking OCI Spec Changes

**Current Versions**:
- OCI Image Spec: v1.1.1 (stable)
- OCI Distribution Spec: v1.1.1 (stable)

**Future Versions**: Monitor for v1.2.0+ proposals

#### Monitoring Process

1. **Watch OCI Repositories**
   - [opencontainers/image-spec](https://github.com/opencontainers/image-spec)
   - [opencontainers/distribution-spec](https://github.com/opencontainers/distribution-spec)

2. **Evaluate Impact**
   - Review proposed changes
   - Assess ORAS Go implementation impact
   - Participate in spec discussions

3. **Implementation**
   - Implement new features in main branch (v3/v4)
   - Backport critical fixes to stable branches (v2)
   - Update documentation and examples

4. **Community Communication**
   - Announce spec compliance updates
   - Provide migration guidance
   - Update compatibility matrix

#### OCI Working Groups

**Participation**:
- OCI Technical Oversight Board (TOB)
- OCI Image Spec working group
- OCI Distribution Spec working group

**Contribution Areas**:
- Artifact manifest discussions
- Referrers API evolution
- Multi-platform image support

---

## Academic & Industry References

### Distributed Systems Concepts

| Concept | Relevance | Learn More |
|---------|-----------|------------|
| **Content-Addressable Storage** | Core storage model in ORAS | [Wikipedia](https://en.wikipedia.org/wiki/Content-addressable_storage) |
| **Directed Acyclic Graphs** | Artifact representation as DAGs | [Wikipedia](https://en.wikipedia.org/wiki/Directed_acyclic_graph) |
| **Semantic Versioning** | Version management strategy | [semver.org](https://semver.org/) |

#### Content-Addressable Storage (CAS)

**Definition**: Storage system where data is identified by its content, not location

**ORAS Implementation**:
- Blob storage using SHA256 digests
- Automatic deduplication by content
- Immutable content references
- Efficient content sharing

**Key Files**: [internal/cas/](../../../internal/cas/)

#### Directed Acyclic Graphs (DAG)

**Definition**: Graph structure with directed edges and no cycles

**ORAS Usage**:
- Artifact layers form DAG structure
- Manifest points to config and layers
- Index points to multiple manifests
- Efficient traversal algorithms

**Key Files**: [content/graph.go](../../../content/graph.go), [copy.go](../../../copy.go)

### Standards Organizations

| Organization | Role | Relevance | URL |
|-------------|------|-----------|-----|
| **Open Container Initiative (OCI)** | Specification steward | Defines image/distribution specs | [opencontainers.org](https://opencontainers.org/) |
| **Cloud Native Computing Foundation (CNCF)** | Project host | Hosts ORAS project | [cncf.io](https://cncf.io/) |
| **Internet Engineering Task Force (IETF)** | Protocol standards | HTTP/TLS used by registries | [ietf.org](https://ietf.org/) |

#### OCI (Open Container Initiative)

**Mission**: Create open industry standards for container formats and runtimes

**ORAS Relationship**: 
- Implements OCI specifications
- Provides reference Go implementation
- Participates in spec evolution

**Governance**: Part of Linux Foundation

#### CNCF (Cloud Native Computing Foundation)

**ORAS Status**: Sandbox project  
**Benefits**: Community support, infrastructure, ecosystem integration  
**Other CNCF Projects**: Kubernetes, containerd, Helm, Notary

---

## Troubleshooting Resources

### Common Issues

#### Issue: Authentication Failures

**Symptoms**: 401 Unauthorized errors  
**Causes**: Missing/incorrect credentials, expired tokens  
**Solutions**:
- Check `~/.docker/config.json`
- Use `WithAuth` option explicitly
- Verify registry URL format

**Reference**: [registry/remote/auth/](../../../registry/remote/auth/)

#### Issue: Descriptor Not Found

**Symptoms**: "not found" errors during copy  
**Causes**: Source content missing, wrong reference  
**Solutions**:
- Verify source content exists
- Check digest/tag accuracy
- Use `Exists()` to pre-validate

**Reference**: [content/storage.go](../../../content/storage.go)

#### Issue: OCI Layout Validation Errors

**Symptoms**: Invalid reference format errors  
**Causes**: Reserved names, invalid characters  
**Solutions**:
- Use alphanumeric + [._-] characters
- Avoid reserved names (., ..)
- Max 255 characters

**Reference**: [content/oci/oci.go](../../../content/oci/oci.go#L624)

### Debugging Tips

**Enable Verbose Logging**:
```go
// Context with logging
ctx := context.WithValue(ctx, "debug", true)
```

**Check Storage State**:
```go
// Verify descriptor exists
exists, err := store.Exists(ctx, desc)
```

**Inspect Manifests**:
```go
// Parse manifest
manifest, err := manifestutil.Parse(manifestBytes)
```

---

## Version Support Matrix

### Go Version Support

| ORAS Go Version | Go 1.23 | Go 1.24 | Go 1.25 | Go 1.26 |
|-----------------|---------|---------|---------|---------|
| v1.x | ✅ | ✅ | ⚠️ | ❌ |
| v2.x | ⚠️ | ✅ | ✅ | ✅ |
| v3.x | ❌ | ✅ | ✅ | ✅ |

**Legend**:
- ✅ Fully supported and tested
- ⚠️ Compatible but not officially supported
- ❌ Not compatible

**Policy**: Support 2 latest Go versions per [Go Security Policy](https://github.com/golang/go/security/policy)

### Platform Support

| OS | Architecture | v1.x | v2.x | v3.x |
|----|-------------|------|------|------|
| Linux | amd64 | ✅ | ✅ | ✅ |
| Linux | arm64 | ✅ | ✅ | ✅ |
| macOS | amd64 | ✅ | ✅ | ✅ |
| macOS | arm64 | ✅ | ✅ | ✅ |
| Windows | amd64 | ✅ | ✅ | ✅ |
| FreeBSD | amd64 | 🔶 | 🔶 | 🔶 |

**Legend**:
- ✅ Officially supported and tested
- 🔶 Community supported (not regularly tested)

---

## Additional Resources

### Books & Articles

**Container Technology**:
- "Docker Deep Dive" by Nigel Poulton
- "Kubernetes Patterns" by Bilgin Ibryam & Roland Huß
- OCI specification documents (authoritative)

**Go Programming**:
- "The Go Programming Language" by Donovan & Kernighan
- "Concurrency in Go" by Katherine Cox-Buday
- [Effective Go](https://go.dev/doc/effective_go)

### Training & Workshops

**ORAS Resources**:
- ORAS Project documentation: [oras.land](https://oras.land/)
- CNCF webinars and KubeCon talks
- Community meeting recordings

**OCI/Container**:
- OCI specification workshops
- KubeCon + CloudNativeCon sessions
- CNCF online training

### Tools & Utilities

**Development**:
- [Go Playground](https://go.dev/play/) - Test Go code online
- [godoc.org](https://pkg.go.dev/) - Go documentation
- [goreportcard.com](https://goreportcard.com/) - Code quality

**Registry Tools**:
- [ORAS CLI](https://github.com/oras-project/oras) - Command-line client
- [Docker Hub](https://hub.docker.com/) - Public registry
- [Harbor](https://goharbor.io/) - Private registry

---

## Maintenance References

### Maintenance Documentation

Related maintenance documents:

- [Unused Files](01-unused-files.md) - Zero unused files analysis
- [Quick Wins](02-quick-wins.md) - 8 high-impact improvements
- [Roadmap](03-roadmap.md) - Version strategy and timeline

### Maintenance Checklist

#### Monthly Tasks
- Review new GitHub issues
- Update dependencies (`go get -u`)
- Run full test suite with race detection
- Check security advisories
- Review TODO comments

#### Quarterly Tasks
- Evaluate deprecated API usage
- Review performance benchmarks
- Update migration guides
- Review OCI spec changes
- Community feedback analysis

**Reference**: [Roadmap - Maintenance Timeline](03-roadmap.md#maintenance-timeline)

---

## Conclusion

This reference guide provides comprehensive resources for ORAS Go development, covering:

✅ **Specifications**: OCI, Docker, Go standards  
✅ **Documentation**: Project docs, API references, code analysis  
✅ **Community**: Communication channels, related projects  
✅ **Development**: Testing, building, dependencies  
✅ **Evolution**: Spec tracking, version support  

**Stay Updated**:
- Watch ORAS Go repository
- Join CNCF Slack #oras channel
- Follow OCI specification repos
- Attend community meetings

---

## Quick Reference Links

### Most Used Resources

**Documentation**:
- [README.md](../../../README.md) - Start here
- [MIGRATION_GUIDE.md](../../../MIGRATION_GUIDE.md) - Version migration
- [pkg.go.dev](https://pkg.go.dev/oras.land/oras-go/v2) - API docs

**Specifications**:
- [OCI Image Spec v1.1.1](https://github.com/opencontainers/image-spec/blob/v1.1.1/spec.md)
- [OCI Distribution Spec v1.1.1](https://github.com/opencontainers/distribution-spec/blob/v1.1.1/spec.md)

**Community**:
- [GitHub Issues](https://github.com/oras-project/oras-go/issues)
- [CNCF Slack #oras](https://cloud-native.slack.com/archives/CJ1KHJM5Z)
- [ORAS Website](https://oras.land/)

---

**Document Version**: 1.0  
**Last Updated**: December 16, 2025  
**Maintained By**: ORAS Go Maintainers
