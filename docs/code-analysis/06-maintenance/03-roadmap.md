# Roadmap

**Module**: Maintenance Analysis  
**Component**: Roadmap  
**Last Updated**: December 16, 2025

[⬆️ Back to Quality](../05-quality/01-testing.md)

---

## Overview

This document outlines ORAS Go's development roadmap, version strategy, and long-term vision. The project follows semantic versioning with clear deprecation policies and community-driven feature planning.

### Version Status

| Version | Branch | Status | Go Support | Target Release |
|---------|--------|--------|------------|----------------|
| v1.x | v1 | Maintenance | 1.23, 1.24 | Security fixes only |
| v2.x | v2 | Stable | 1.24, 1.25 | Active maintenance |
| v3.x | main | Development | 1.24, 1.25 | TBD (2026) |
| v4.x | N/A | Planned | TBD | TBD (2027+) |

---

## Version Strategy

### Support Policy

**Go Version Support**: 2 latest Go versions per [Go Security Policy](https://github.com/golang/go/security/policy)

**Version Lifecycle**:
1. **Development**: Active feature development on main branch
2. **Stable**: Production-ready, receives features + bug fixes
3. **Maintenance**: Bug fixes + security patches only
4. **End of Life**: No further updates

### Branch Strategy

**Reference**: [README.md](../../../README.md#L45-L70)

- **main**: v3 development (breaking changes allowed)
- **v2**: Stable v2.x releases (backward compatible changes only)
- **v1**: Maintenance v1.x releases (critical fixes only)

### Versioning Policy

**Semantic Versioning** (semver.org):
- **Major**: Breaking API changes (v2 → v3)
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, security patches

**Import Path Changes**: Major versions require new import path
- v2: `oras.land/oras-go/v2`
- v3: `oras.land/oras-go/v3` (future)

---

## v1.x - Maintenance Mode

### Status

**Current Version**: v1.x (latest)  
**Support Level**: Maintenance only  
**End of Life**: TBD (estimated 2025-2026)

### Support Scope

**Included**:
- ✅ Security vulnerability patches
- ✅ Critical bug fixes (data loss, corruption)
- ✅ Go version compatibility updates

**Excluded**:
- ❌ New features
- ❌ Non-critical bug fixes
- ❌ API enhancements
- ❌ Performance optimizations

### Migration Path

**Recommendation**: Upgrade to v2.x for continued feature support

**Migration Guide**: [MIGRATION_GUIDE.md](../../../MIGRATION_GUIDE.md)

**Key Changes**:
- Import path update
- API surface refinements
- Enhanced error handling

---

## v2.x - Stable Production Release

### Status

**Current Version**: v2.x (latest stable)  
**Support Level**: Active maintenance  
**Expected Support**: Through 2027+

### Features

**Core Capabilities**:
- ✅ Copy operations with graph traversal
- ✅ Manifest packing (OCI image spec compliant)
- ✅ Multiple storage backends (file, memory, OCI)
- ✅ Remote registry client (OCI distribution spec)
- ✅ Authentication & credential management
- ✅ Extended copy with referrers support
- ✅ Platform selection for multi-arch images

**Stability**:
- Production-tested
- 90.2% test coverage
- Extensive example documentation
- Community adoption (containerd, Helm, Notary)

### Maintenance Commitment

**Active Support Includes**:
- Bug fixes (all severity levels)
- Security patches
- Performance improvements (non-breaking)
- New features (backward compatible)
- Documentation updates
- Community support

**Release Cadence**:
- Patch releases: As needed (typically monthly)
- Minor releases: Quarterly (with new features)

### Backward Compatibility

**Guarantee**: v2.x releases maintain API compatibility
- No breaking changes within v2 series
- New features additive only
- Deprecations marked clearly

---

## v3.x - Active Development

### Status

**Branch**: main  
**Status**: Active development (breaking changes expected)  
**Target Release**: TBD (estimated Q1-Q2 2026)  
**Current State**: API refinements in progress

### Goals

**Primary Objectives**:

1. **API Simplification**
   - Cleaner interfaces
   - Reduced API surface
   - Better naming consistency

2. **Performance Optimization**
   - Improved concurrency controls
   - Reduced memory footprint
   - Optimized caching strategies

3. **Enhanced Capabilities**
   - Better referrers API integration
   - Improved credential management
   - More robust retry mechanisms

4. **Developer Experience**
   - Clearer error messages
   - Better debugging support
   - Comprehensive examples

### Breaking Changes from v2

**Import Path**: `oras.land/oras-go/v3`

**Known Changes** (subject to finalization):
- API signatures may change
- Helper functions relocated
- Enhanced type safety requirements
- Refined option patterns

**Migration Support**:
- [MIGRATION_GUIDE.md](../../../MIGRATION_GUIDE.md) will be updated
- Examples for all breaking changes
- Community preview period before release

### Development Timeline

| Phase | Timeline | Status |
|-------|----------|--------|
| API Design | Q4 2025 | 🟡 In Progress |
| Implementation | Q1 2026 | 🔵 Upcoming |
| Testing & Refinement | Q1-Q2 2026 | ⚪ Planned |
| API Freeze | Q2 2026 | ⚪ Planned |
| Release | Q2 2026 | ⚪ Planned |

**Note**: Timeline subject to community feedback and technical requirements

### Feature Roadmap

#### Planned Features

**High Priority**:
- ✅ Leaf node copy optimization ([Quick Win #1](02-quick-wins.md#quick-win-1))
- ⚪ Enhanced OCI validation ([Quick Win #3](02-quick-wins.md#quick-win-3))
- ⚪ Comprehensive benchmark suite ([Quick Win #2](02-quick-wins.md#quick-win-2))
- ⚪ Performance metrics API ([Quick Win #7](02-quick-wins.md#quick-win-7))

**Medium Priority**:
- ⚪ Improved progress reporting
- ⚪ Enhanced credential caching
- ⚪ Better multi-registry support
- ⚪ Streaming optimizations for large artifacts

**Community-Requested**:
- ⚪ Multi-platform artifact handling examples
- ⚪ Custom credential helper templates
- ⚪ Error recovery patterns
- ⚪ Advanced referrers usage

### API Stability

**Current State**: Unstable (main branch)
- Breaking changes may occur without notice
- Use v2 for production workloads

**Stabilization Process**:
1. Community RFC for major changes
2. Preview releases (alpha, beta, RC)
3. Feedback incorporation period
4. API freeze announcement
5. Final release

---

## v4.x - Future Planning

### Status

**Timeline**: 2027+ (tentative)  
**Current State**: Conceptual planning only  
**Dependencies**: v3 community feedback

### Planned Removals

#### Deprecated API Cleanup

From [pack.go](../../../pack.go):

**PackManifestVersion1_1_RC4** (Line 76-78)
- **Type**: Constant
- **Deprecated**: v2.0.0
- **Reason**: RC version superseded by stable v1.1
- **Migration**: Use `PackManifestVersion1_1`

**PackOptions** (Line 152)
- **Type**: Type alias
- **Deprecated**: v2.0.0
- **Reason**: Naming clarity (manifest vs. general packing)
- **Migration**: Use `PackManifestOptions`

**Pack()** (Line 188)
- **Type**: Function
- **Deprecated**: v2.0.0
- **Reason**: Replaced by more flexible `PackManifest()`
- **Migration**: Use `PackManifest()` with explicit version

#### Deprecation Policy

**Minimum Deprecation Period**: 2 major versions (v2 → v3 → v4)
- Deprecated in v2.x
- Maintained through v3.x
- Removed in v4.x

**Timeline**: ~4-5 years from deprecation to removal
- Provides ample migration time
- Supports long-term projects
- Maintains ecosystem stability

### v4 Feature Vision

#### API Streamlining

**Goals**:
- Single, consistent packing API
- Unified naming conventions
- Reduced API surface area
- Improved type safety

**Example**: Single manifest packing function

```go
// v4 vision (conceptual)
func PackManifest(ctx context.Context, pusher Pusher, version ManifestVersion, opts PackOptions) (Descriptor, error)
```

#### Enhanced Capabilities

**Potential Features** (subject to community input):

1. **Graph Operations**
   - More efficient DAG traversal algorithms
   - Better handling of large artifact graphs
   - Improved referrers API integration
   - Graph visualization utilities

2. **Storage Backends**
   - Plugin architecture for custom stores
   - S3/cloud-native backend support
   - Better streaming for massive artifacts
   - Distributed storage coordination

3. **Registry Compatibility**
   - Broader OCI registry support
   - Docker Hub API enhancements
   - Better error recovery for unreliable networks
   - Multi-registry mirroring/replication

4. **Developer Experience**
   - Interactive debugging tools
   - Performance profiling utilities
   - Better logging and observability
   - Enhanced error diagnostics

### Community Input

**Feedback Channels**:
- GitHub Discussions: Feature proposals
- CNCF Slack #oras: Real-time discussion
- GitHub Issues: Specific requests
- Community meetings: Design reviews

**RFC Process**:
- Major changes require RFC (Request for Comments)
- Community review period (minimum 2 weeks)
- Consensus-based decision making
- Transparent decision documentation

---

## Long-Term Vision

### Specification Alignment

#### Current Compliance

**OCI Specifications**:
- ✅ OCI Image Spec v1.1.1
- ✅ OCI Distribution Spec v1.1.1
- ✅ OCI Image Layout v1.1.1
- ✅ OCI Descriptor format
- ✅ OCI Manifest format

**Docker Standards**:
- ✅ Docker Registry HTTP API V2
- ✅ Docker credential helpers
- ✅ Docker token authentication

#### Future Tracking

**Monitoring**:
- Watch OCI spec repositories for updates
- Participate in OCI working groups
- Track CNCF artifact standards evolution
- Monitor Docker API changes

**Spec Evolution Strategy**:
1. Monitor v1.2+ proposals in OCI repos
2. Evaluate impact on ORAS Go
3. Implement in main branch (v3/v4)
4. Backport critical fixes to stable branches

**Key Repositories**:
- [opencontainers/image-spec](https://github.com/opencontainers/image-spec)
- [opencontainers/distribution-spec](https://github.com/opencontainers/distribution-spec)

### Platform Support

#### Current Support

**Operating Systems**:
- ✅ Linux (all major distributions)
- ✅ macOS (Intel + Apple Silicon)
- ✅ Windows (10+, Server 2016+)
- ✅ FreeBSD, OpenBSD (community support)

**Architectures**:
- ✅ amd64 (x86-64)
- ✅ arm64 (aarch64)
- ✅ arm (32-bit)
- ✅ 386, ppc64le, s390x (via Go)

#### Future Enhancements

**Planned Improvements**:
- Better Windows path handling
- ARM64 performance optimization
- Enhanced cross-platform testing
- Platform-specific optimization paths

### Ecosystem Integration

#### Current Integrations

**Major Projects Using ORAS Go**:
- **containerd**: Artifact management
- **Notary v2**: Supply chain security
- **Helm**: OCI registry support
- **Artifact Hub**: Artifact distribution

#### Future Integration Goals

**Target Ecosystems**:

1. **Kubernetes Ecosystem**
   - Operator support
   - CRD integration
   - Native image management

2. **CI/CD Platforms**
   - GitHub Actions helpers
   - GitLab CI integration
   - Jenkins plugins
   - ArgoCD artifact support

3. **Cloud-Native Tools**
   - Tekton pipeline tasks
   - Flux CD integration
   - Kustomize artifact support
   - KubeVirt image management

4. **Security Tools**
   - Sigstore integration
   - Policy enforcement (OPA/Kyverno)
   - Vulnerability scanning
   - SBOM generation/distribution

---

## Feature Priorities

### Prioritization Framework

**Criteria**:
1. **User Impact**: How many users benefit?
2. **Spec Compliance**: Required by OCI/Docker standards?
3. **Ecosystem Demand**: Community requests and adoption blockers?
4. **Technical Debt**: Does it address maintenance burden?
5. **Implementation Cost**: Development effort vs. value

**Priority Levels**:
- 🔴 **Critical**: Spec compliance, security, data integrity
- 🟡 **High**: Major user pain points, ecosystem adoption
- 🟢 **Medium**: Quality of life, performance optimization
- 🔵 **Low**: Nice-to-have, limited user impact

### Current Priorities

#### v3 Release (Immediate)

🔴 **Critical**:
- API stabilization and freeze
- Migration guide completion
- Comprehensive testing
- Security audit

🟡 **High**:
- Performance optimization ([Quick Wins](02-quick-wins.md))
- Enhanced error messages
- Benchmark suite
- Example expansion

#### Post-v3 (Short-Term)

🟢 **Medium**:
- Performance metrics API
- Enhanced observability
- Multi-registry mirroring
- Streaming optimizations

🔵 **Low**:
- Advanced debugging tools
- Interactive playground
- Graph visualization
- Custom storage plugin examples

#### v4 Planning (Long-Term)

🟡 **High**:
- Deprecation removal
- API simplification
- Community RFC process

🟢 **Medium**:
- Plugin architecture exploration
- Cloud-native backend support
- Enhanced referrers API

---

## Maintenance Timeline

### 2026 Roadmap

| Quarter | Focus Areas | Deliverables |
|---------|-------------|--------------|
| **Q1 2026** | v3 stabilization | API freeze, beta releases |
| **Q2 2026** | v3 release | v3.0.0 GA, migration guide |
| **Q3 2026** | Community feedback | v3 patches, v2 maintenance |
| **Q4 2026** | v4 planning | Deprecation timeline, RFC process |

### 2027+ Vision

| Year | Major Goals |
|------|-------------|
| **2027** | v4 development, v2 EOL planning |
| **2028** | v4 release, ecosystem expansion |
| **2029+** | Spec evolution tracking, next-gen features |

**Note**: Timeline subject to community needs, specification changes, and resource availability

---

## Community-Driven Development

### Feedback Mechanisms

#### Communication Channels

| Channel | Purpose | Response Time |
|---------|---------|---------------|
| **GitHub Issues** | Bug reports, feature requests | 1-3 business days |
| **GitHub Discussions** | Q&A, design discussions | 1-5 business days |
| **CNCF Slack #oras** | Real-time support | Best-effort (hours) |
| **Community Meetings** | Design reviews, planning | Monthly |

**Links**:
- [GitHub Issues](https://github.com/oras-project/oras-go/issues)
- [GitHub Discussions](https://github.com/oras-project/oras-go/discussions)
- [CNCF Slack](https://cloud-native.slack.com/archives/CJ1KHJM5Z)

#### Feature Request Process

1. **Submit GitHub Issue** with "enhancement" label
2. **Community Discussion**: Gather feedback, refine proposal
3. **Maintainer Review**: Assess fit with roadmap
4. **RFC (if major)**: Formal design document
5. **Implementation**: Contributor or maintainer develops
6. **Review & Merge**: Code review, testing, documentation
7. **Release**: Include in next appropriate version

### Contributor Opportunities

#### Good First Issues

**Areas for New Contributors**:
- Documentation improvements
- Example test expansion
- Error message enhancement
- Test coverage improvements

#### Advanced Contributions

**Complex Features**:
- Performance optimizations
- New storage backends
- Registry client enhancements
- Concurrency improvements

**Guidance**: [CONTRIBUTING.md](../../../CODE_OF_CONDUCT.md) and [OWNERS.md](../../../OWNERS.md)

---

## Success Metrics

### Version Adoption

**Tracking**:
- Go module proxy download statistics
- GitHub star/fork growth
- Ecosystem project integration count
- Community channel activity

**Targets**:
- v2 adoption: Maintain through 2027
- v3 adoption: 50%+ within 6 months of release
- v4 adoption: Planned based on v3 learnings

### Quality Metrics

**Ongoing Measurement**:
- Test coverage: Maintain > 90%
- GitHub issue response time: < 3 business days
- Security vulnerability patches: < 7 days
- Documentation completeness: > 95%

### Community Health

**Indicators**:
- Active contributors (monthly)
- GitHub issue closure rate
- Community meeting attendance
- Slack channel engagement

---

## Risk Management

### Technical Risks

| Risk | Mitigation |
|------|------------|
| **Spec Divergence** | Regular OCI spec tracking, working group participation |
| **Breaking Changes** | Long deprecation periods, clear migration paths |
| **Performance Regression** | Benchmark suite, continuous monitoring |
| **Security Vulnerabilities** | Regular audits, prompt patch releases |

### Community Risks

| Risk | Mitigation |
|------|------------|
| **Fragmentation** | Clear version strategy, backward compatibility |
| **Contributor Burnout** | Distributed maintenance, clear ownership |
| **Adoption Barriers** | Comprehensive docs, examples, migration support |
| **Ecosystem Conflicts** | Close coordination with major integrators |

---

## Conclusion

### Roadmap Summary

**2026 Focus**: v3 stabilization and release  
**Long-Term Vision**: Continued specification alignment, ecosystem expansion, community growth

**Key Commitments**:
- ✅ Clear version strategy with long support periods
- ✅ Community-driven feature development
- ✅ Backward compatibility where possible
- ✅ Transparent deprecation and migration paths
- ✅ Regular maintenance and security updates

### Next Steps

**Immediate (Q1 2026)**:
1. Complete v3 API refinements
2. Implement quick wins ([see list](02-quick-wins.md))
3. Expand test coverage and benchmarks
4. Prepare v3 beta releases

**Short-Term (2026)**:
1. v3.0.0 GA release
2. Community feedback incorporation
3. v4 deprecation timeline documentation
4. Ecosystem integration expansion

**Long-Term (2027+)**:
1. v4 planning and development
2. Advanced features (plugins, cloud backends)
3. Next-generation OCI spec support
4. Continued ecosystem growth

---

## Related Documentation

- [Quick Wins](02-quick-wins.md) - Near-term improvements
- [References](04-references.md) - Specifications and resources
- [MIGRATION_GUIDE.md](../../../MIGRATION_GUIDE.md) - Version migration instructions

---

**Document Version**: 1.0  
**Last Updated**: December 16, 2025  
**Next Review**: March 2026 (Quarterly)
