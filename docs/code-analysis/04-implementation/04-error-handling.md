# Implementation Analysis: Error Handling

**Project:** ORAS Go (OCI Registry As Storage)  
**Phase:** 4 - Implementation Analysis  
**Document:** 04-error-handling.md  
**Analysis Date:** December 16, 2025  
**Confidence Level:** High (95-100%)

---

## Navigation

- **Phase 4 Overview:** [README](README.md)
- **Previous Document:** [03-hotspots.md](03-hotspots.md)
- **Related:** [Module Analysis](../03-modules/README.md)

---

## 1. Error Definition Strategy

### 1.1 Sentinel Errors

**Pattern:** Simple sentinel errors using `errors.New()` for common error conditions  
**Confidence:** 100%

**Location:** [errdef/errors.go](../../../errdef/errors.go)

```go
var (
    ErrAlreadyExists      = errors.New("already exists")
    ErrInvalidDigest      = errors.New("invalid digest")
    ErrInvalidReference   = errors.New("invalid reference")
    ErrInvalidMediaType   = errors.New("invalid media type")
    ErrMissingReference   = errors.New("missing reference")
    ErrNotFound           = errors.New("not found")
    ErrSizeExceedsLimit   = errors.New("size exceeds limit")
    ErrUnsupported        = errors.New("unsupported")
    ErrUnsupportedVersion = errors.New("unsupported version")
)
```

**Characteristics:**
- Exported package-level variables
- Simple string-based errors
- Used for well-defined, common error conditions
- Suitable for `errors.Is()` comparisons

---

### 1.2 Domain-Specific Errors

**Pattern:** Module-specific errors with clear naming and visibility (exported vs internal)  
**Confidence:** 100%

**Location:** [content/file/errors.go](../../../content/file/errors.go)

```go
var (
    ErrMissingName             = errors.New("missing name")
    ErrDuplicateName           = errors.New("duplicate name")
    ErrPathTraversalDisallowed = errors.New("path traversal disallowed")
    ErrOverwriteDisallowed     = errors.New("overwrite disallowed")
    ErrStoreClosed             = errors.New("store already closed")
)
var errSkipUnnamed = errors.New("unnamed descriptor skipped")  // internal use
```

**Design Decisions:**
- Domain-specific errors isolated to relevant packages
- Clear distinction between public API errors (exported) and internal errors (unexported)
- Self-documenting error messages

---

## 2. Structured Error Types

### 2.1 CopyError - Contextual Error Type

**Pattern:** Structured error with context and origin tracking  
**Confidence:** 100%

**Location:** [copyerror.go](../../../copyerror.go)

```go
type CopyErrorOrigin int

const (
    CopyErrorOriginSource      CopyErrorOrigin = 1  // Source side error
    CopyErrorOriginDestination CopyErrorOrigin = 2  // Destination side error
)

type CopyError struct {
    Op     string          // Operation that caused error
    Origin CopyErrorOrigin // Source of error
    Err    error           // Underlying error
}

func (e *CopyError) Error() string {
    switch e.Origin {
    case CopyErrorOriginSource, CopyErrorOriginDestination:
        return fmt.Sprintf("failed to perform %q on %s: %v", e.Op, e.Origin, e.Err)
    default:
        return fmt.Sprintf("failed to perform %q: %v", e.Op, e.Err)
    }
}

func (e *CopyError) Unwrap() error {
    return e.Err
}
```

**Key Features:**
- Implements `error` interface with custom `Error()` method
- Implements `Unwrap()` for error chain inspection
- Provides contextual information (operation name, origin)
- Constructor ensures nil-safety: `newCopyError` returns nil if err is nil

**Usage Benefits:**
- Clear identification of where errors occurred (source vs destination)
- Preserves original error for inspection
- Provides operation context for debugging

---

### 2.2 Registry Error Codes

**Pattern:** OCI Distribution Spec compliant error representation  
**Confidence:** 100%

**Location:** [registry/remote/errcode/errors.go](../../../registry/remote/errcode/errors.go)

```go
type Error struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Detail  any    `json:"detail,omitempty"`
}

type Errors []Error  // Collection of errors

func (errs Errors) Unwrap() error {
    if len(errs) == 1 {
        return errs[0]
    }
    return nil
}
```

**Standard Codes:**
- `BLOB_UNKNOWN`
- `MANIFEST_INVALID`
- `UNAUTHORIZED`
- `UNSUPPORTED`
- And others per OCI Distribution specification

**Purpose:**
- Standards-compliant error handling for registry operations
- Structured error information for API responses
- Support for multiple errors in a single response

---

## 3. Error Wrapping Patterns

### 3.1 Usage in Copy Operations

**Pattern:** Context-aware error wrapping with origin tracking  
**Confidence:** 100%

**Location:** [copy.go](../../../copy.go)

```go
// Pattern 1: Wrapping with origin context
func Copy(ctx context.Context, src ReadOnlyTarget, srcRef string, 
          dst Target, dstRef string, opts CopyOptions) (ocispec.Descriptor, error) {
    if src == nil {
        return ocispec.Descriptor{}, newCopyError("Copy", CopyErrorOriginSource, 
            errors.New("nil source target"))
    }
    if dst == nil {
        return ocispec.Descriptor{}, newCopyError("Copy", CopyErrorOriginDestination, 
            errors.New("nil destination target"))
    }
    // ...
}

// Pattern 2: Wrapping operation errors
if err := dst.Exists(ctx, desc); err != nil {
    return newCopyError("Exists", CopyErrorOriginDestination, err)
}

if err := opts.FindSuccessors(ctx, proxy, desc); err != nil {
    return newCopyError("FindSuccessors", CopyErrorOriginSource, err)
}
```

**Benefits:**
- Immediate identification of error location
- Consistent error handling across operations
- Preserved error chain for inspection

---

### 3.2 fmt.Errorf with %w Verb

**Pattern:** Add context while preserving error chain  
**Confidence:** 100%

**Location:** [registry/remote/repository.go](../../../registry/remote/repository.go)

```go
// Pattern: Add context while preserving error chain
return fmt.Errorf("%s %q: failed to decode response: %w", 
    resp.Request.Method, resp.Request.URL, err)

return fmt.Errorf("failed to query referrers API: %w", errdef.ErrUnsupported)

return fmt.Errorf("%w: current capability = %v, new capability = %v",
    ErrReferrersCapabilityAlreadySet, fact == referrersStateSupported, capable)
```

**Best Practice:**
- Always use `%w` verb (not `%v`) to maintain error chain
- Enables `errors.Is()` and `errors.As()` to work correctly
- Adds contextual information (HTTP method, URL, state values)
- Follows Go 1.13+ error wrapping conventions

---

## 4. Error Inspection

### 4.1 errors.Is Usage

**Pattern:** Check for specific sentinel errors in error chain  
**Confidence:** 100%

**Examples from:** [registry/remote/repository.go](../../../registry/remote/repository.go)

```go
if errors.Is(err, errdef.ErrUnsupported) {
    // Handle unsupported operation
}

if errors.Is(err, errdef.ErrNotFound) {
    // Create new referrers index
}
```

**Purpose:**
- Check for specific error types anywhere in the error chain
- Works with wrapped errors (via `%w` or `Unwrap()`)
- More robust than direct equality checks

---

### 4.2 errors.As Usage

**Pattern:** Extract typed errors from error chain  
**Confidence:** 100%

**Example from:** [registry/remote/repository_test.go](../../../registry/remote/repository_test.go)

```go
var errResp *errcode.ErrorResponse
if !errors.As(err, &errResp) {
    t.Fatal("expected ErrorResponse")
}
```

**Purpose:**
- Extract structured error information from error chain
- Type-safe access to error-specific fields
- Enables detailed error handling based on error type

---

## 5. Error Handling Guidelines

**Confidence:** 95% (inferred from consistent patterns across codebase)

Based on analysis of 50+ error usage patterns throughout the codebase:

### 5.1 Core Principles

1. **Return errors, don't panic**
   - Panics only in test setup and unrecoverable scenarios
   - Production code always returns errors for caller handling

2. **Wrap errors with context**
   - Use `fmt.Errorf` with `%w` for simple wrapping
   - Use custom error types (e.g., `CopyError`) for structured context

3. **Preserve error chains**
   - Always implement `Unwrap()` on custom error types
   - Enable `errors.Is()` and `errors.As()` inspection

4. **Use sentinel errors**
   - For common, well-defined error conditions
   - Defined at package level with clear names

5. **Check nil before wrapping**
   - Constructor pattern: `if err == nil { return nil }`
   - Prevents unnecessary wrapper allocation

6. **Provide operation context**
   - Include operation name (e.g., "Copy", "Exists", "Fetch")
   - Include affected component (source vs destination)
   - Add relevant state information for debugging

### 5.2 Error Constructor Pattern

```go
func newCopyError(op string, origin CopyErrorOrigin, err error) error {
    if err == nil {
        return nil
    }
    return &CopyError{
        Op:     op,
        Origin: origin,
        Err:    err,
    }
}
```

**Key Elements:**
- Nil check prevents unnecessary wrapping
- Returns pointer to error type
- Consistent signature across codebase

---

## 6. Error Handling Anti-Patterns (Avoided)

Based on codebase analysis, the following anti-patterns are consistently avoided:

### 6.1 What NOT to Do

❌ **Don't discard error context:**
```go
// Bad
return errors.New("operation failed")
```

✅ **Do preserve error chain:**
```go
// Good
return fmt.Errorf("operation failed: %w", err)
```

❌ **Don't use %v for wrapping:**
```go
// Bad - breaks error chain
return fmt.Errorf("failed: %v", err)
```

✅ **Do use %w for wrapping:**
```go
// Good - preserves error chain
return fmt.Errorf("failed: %w", err)
```

❌ **Don't panic in library code:**
```go
// Bad
if err != nil {
    panic(err)
}
```

✅ **Do return errors:**
```go
// Good
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

---

## 7. Summary

The ORAS Go project demonstrates a mature and consistent approach to error handling:

1. **Layered Error Strategy:**
   - Sentinel errors for common conditions
   - Domain-specific errors for modules
   - Structured error types for complex scenarios

2. **Standard Compliance:**
   - Follows Go 1.13+ error wrapping conventions
   - Implements OCI Distribution Spec error codes
   - Consistent `error` interface implementation

3. **Error Chain Preservation:**
   - All custom error types implement `Unwrap()`
   - Consistent use of `%w` in `fmt.Errorf`
   - Enables robust error inspection

4. **Contextual Information:**
   - Operation names included in error messages
   - Origin tracking (source vs destination)
   - Relevant state information for debugging

5. **Production-Ready:**
   - No panics in production code
   - Nil-safe error constructors
   - Clear error ownership and responsibility

**Overall Assessment:** The error handling implementation is well-designed, follows Go best practices, and provides excellent support for debugging and error recovery in production environments.

---

## References

- **Source Files Analyzed:**
  - [errdef/errors.go](../../../errdef/errors.go)
  - [copyerror.go](../../../copyerror.go)
  - [content/file/errors.go](../../../content/file/errors.go)
  - [registry/remote/errcode/errors.go](../../../registry/remote/errcode/errors.go)
  - [copy.go](../../../copy.go)
  - [registry/remote/repository.go](../../../registry/remote/repository.go)

- **Related Documentation:**
  - [Design Patterns](01-patterns.md)
  - [Module Analysis](../03-modules/README.md)

---

**Document Status:** Complete  
**Last Updated:** December 16, 2025  
**Extracted From:** phase-4-findings-part2.md, Section 4
