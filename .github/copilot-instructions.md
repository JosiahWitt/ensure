# Ensure - AI Coding Agent Instructions

## Project Overview

`ensure` is a balanced Go testing framework (Go 1.16+) with integrated mock generation. Core components:

- **Testing library** (`ensuring` package): Fluent assertions with readable diffs
- **Mock generator CLI** (`cmd/ensure`): Wraps GoMock with automatic mock setup
- **Table-driven test plugins** (`internal/plugins`): Auto-initializes mocks and subjects using reflection

## General Instructions

- This is a library, so avoid breaking changes to public APIs unless absolutely necessary.
- Follow Go conventions and idiomatic patterns.
- Prioritize simplicity over complexity.

## Architecture

### Package Structure

- `ensure.go`: Entry point - use `ensure.New(t)` to create test instances (allows package shadowing)
- `ensuring/`: Core implementation with `E` (Ensure) type and `Chain` for assertions
- `internal/plugins/`: Table test automation via reflection
  - `mocks/`: Auto-initializes mock structs with `NEW()` methods
  - `setupmocks/`: Executes `SetupMocks func(*Mocks)` before each test
  - `subject/`: Auto-wires mocks into Subject struct fields by interface matching
- `internal/tablerunner/`: Orchestrates plugin execution for table tests
- `cmd/ensure/`: CLI for mock generation (configured via `.ensure.yml`)

### Key Design Patterns

**Package Shadowing Pattern**: Always use `ensure := ensure.New(t)`, never call `ensuring.InternalCreateDoNotCallDirectly(t)` directly. This allows `ensure` variable to shadow the package while maintaining access to types via `ensuring` package.

**Table-Driven Test Structure** (see `ensuring/run_table_test.go` for examples):

```go
type Mocks struct {
    DB *mock_db.MockDB  // Must have NEW(*gomock.Controller) *MockDB method
}

table := []struct {
    Name       string                   // Required: unique test name
    Mocks      *Mocks                   // Optional: auto-initialized via NEW()
    SetupMocks func(*Mocks)             // Optional: set expectations
    Subject    *MyService               // Optional: auto-wired with mocks
    // ... test-specific fields
}{{
    Name: "with valid input",
    SetupMocks: func(m *Mocks) {
        m.DB.EXPECT().Put("id", data)
    },
}}

ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
    entry := table[i]
    // Test logic using entry.Subject
})
```

**Mock Struct Tags**:

- `` `ensure:"ignore"` ``: Skip field during mock initialization
- `` `ensure:"ignoreunused"` ``: Mark mock optional for subject wiring (prevents "unused mock" errors)

**Mock Generation**: Uses reflection-based template system (`cmd/ensure/internal/mockgen/template.go`). Generates mocks with `NEW(*gomock.Controller)` method automatically.

## Development Workflows

### Running Tests

```bash
make test                # All tests including submodules
make test-coverage       # Generate coverage HTML in tmp/coverage.html
go test ./...            # Current module only
```

### Mock Generation

```bash
make generate-mocks      # Regenerate all mocks per .ensure.yml
ensure mocks generate    # CLI equivalent
ensure mocks tidy        # Remove stale mock files
```

### Linting

```bash
make lint                # Run staticcheck and golangci-lint
```

### Project Structure

- Root module: Test framework library
- `cmd/ensure/`: Separate Go module for CLI (has own go.mod)
- `exp/entable/`: Experimental features (separate module)
- Each has its own Makefile and test suite

## Conventions

### Library Usage Patterns (For Users of Ensure)

The ensure library provides a fluent assertion API. Users should:

- Always initialize with `ensure := ensure.New(t)`
- Use the fluent pattern: `ensure(actualValue).Method(expectedValue)`
- Available assertion methods: `Equals()`, `IsTrue()`, `IsFalse()`, `IsError()`, `IsNotError()`, `IsEmpty()`, etc.
- Example: `ensure(result).Equals(42)` instead of `if result != 42 { t.Error(...) }`

### Test Patterns for This Repository

When writing tests **for the ensure framework itself**:

- Use `<package>_test` package suffix for external tests (e.g., `ensuring_test`, `mocks_test`)
  - This prevents import cycles when testing the public API
  - Required for files in `ensuring/`, `cmd/ensure/internal/`, and other internal packages
- **Testing ensure's own code uses traditional Go testing patterns** with `if` statements and `t.Error()`
  - The library's tests verify low-level behavior, so they use standard Go testing
  - Example from `ensuring/ensuring_test.go`: `if ensure == nil { t.Error("expected ensure not to be nil") }`
- Use ensure's fluent API in integration tests and when testing higher-level components
  - Example from `internal/plugins/mocks/mocks_test.go`: `ensure(err).IsError(entry.ExpectedError)`
- Helper methods should call `c.t.Helper()` and `c.markRun()` to track usage

### File Organization

- Mock files: `internal/mocks/<package_path>/mock_<package>.go`
- Tests alongside implementation: `foo.go` → `foo_test.go`
- Internal test helpers: `internal/testhelper/`

### Error Handling

Uses `github.com/JosiahWitt/erk` for structured errors. See `cmd/ensure/internal/ensurefile/ensurefile.go` for pattern:

```go
var ErrCannotFindGoModule = erk.New(ErkCannotLoadConfig{}, "Cannot find root go.mod file...")
```

### Dependencies

- `go.uber.org/mock/gomock`: Mock framework (ensure wraps this)
- `github.com/go-test/deep`: Generates readable diffs for `Equals()`
- `github.com/JosiahWitt/erk`: Structured error handling
- `golang.org/x/mod/modfile`: Parse go.mod files

## Key Implementation Details

### Plugin System

Plugins implement `plugins.TablePlugin` interface with `ParseEntryType(reflect.Type)` for validation and `TableEntryHooks` for `BeforeEntry`/`AfterEntry` hooks. Execution order matters:

1. `mocks` plugin: Initialize mock structs
2. `setupmocks` plugin: Call SetupMocks functions
3. `subject` plugin: Wire mocks into subject

### Mock Initialization Flow

1. `mocks` plugin finds struct fields with pointer-to-struct type
2. Validates `NEW(*gomock.Controller) *T` method exists
3. Calls `NEW()` with scoped gomock.Controller before each test
4. `subject` plugin matches mock interfaces to Subject fields by type

### Reflection Usage

Extensive use of `reflect` package for table test automation. See `internal/plugins/internal/iterate/struct_fields.go` for field traversal patterns including embedded struct handling.

## Common Pitfalls

- **Don't** call functions in `SetupMocks` that modify `entry` - reflection values are read-only
- **Don't** forget `Name string` field in table structs (required by name plugin)
- **Do** use `ensure:"ignoreunused"` tag when subject doesn't use all mocks
- **Remember** `.ensure.yml` must be at module root (next to go.mod)
