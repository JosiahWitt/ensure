---
name: ensure-testing
description: Comprehensive guide for using the `ensure` Go testing framework. Always use this skill when writing Go tests that utilize the `ensure` library, such as for fluent assertions, automatic mock generation, and table-driven test automation. It covers critical setup patterns for assertions and the strict structural requirements for table-driven tests and mocks.
---

# Using the `ensure` Testing Library
> **CRITICAL:** The `ensure` library is not widely used, thus, if you don't read this entire skill, you'll likely make incorrect assumptions about how the library works. Thus, it is VERY important that you read the entire file.

> **IMPORTANT:** If you are implementing **Table-Driven Tests** or **Mocks**, you MUST follow the structural requirements exactly. Failure to do so will cause tests to panic or fail to initialize.

## Packages
- `github.com/JosiahWitt/ensure`: Only used for `ensure := ensure.New(t)`, since the `ensure` namespace will be shadowed.
- `github.com/JosiahWitt/ensure/ensuring`: Used for all the other types (eg. `ensuring.E`, stored in the `ensure` variable) and methods.

## Quick Start: Basic Assertions
`ensure` uses a fluent API for assertions.

1. **Initialize**: You MUST create a local `ensure` variable to shadow the package name. This is the core design pattern of the library. This should be the first line of every `Test*` function.
   ```go
   func TestAbc(t *testing.T) {
	   ensure := ensure.New(t)
	   // ...
   }
   ```
2. **Assert**: Use the pattern: `ensure(actual).Method(expected)`. See a more complete list below. Here are some examples:
   ```go
   ensure(result).Equals(42)
   ensure(err).IsError(expectedErr)
   ensure(slice).IsEmpty()
   ```

## Gotchas
### Mandatory Shadowing
**NEVER** name the `ensure.New(t)` result or the `Run*` callback parameter anything other than `ensure`. This shadowing is required for the library's fluent API to work as intended.

### Subject Wiring
For table-driven tests, the `Subject` is auto-wired by matching mock interfaces in the `Mocks` struct to fields in the `Subject` struct. If they don't match, the subject will not be initialized correctly. NEVER initialize the subject directly.

### Mock Initializing
For table-driven tests, the `Mock` struct is auto-initialized, as detailed below. NEVER initialize the mocks directly.

### Mock Generation
`ensure mocks generate` generates mock code per `.ensure.yml`, as detailed below. Mocks need to be regenerated if interfaces change, otherwise they will cause unexpected errors.

### Incorrectly Scoped `ensure`
Test reporting is weird if an `ensure` instance from a parent scope is referenced in a subtest. Thus, it is SUPER IMPORTANT that the `ensure` variable is shadowed correctly in any `Run*` method, or in `SetupMocks` when applicable.

## Fluent Assertion API
Use these methods for assertions. The `ensure(value)` call returns a chain for the check.
*Note: This API may be extended in future versions, so be sure to double check to see if an assertion was added if you need specific functionality not in this list.*

| Method | Description | Example |
| :--- | :--- | :--- |
| `Equals(expected)` | Deep equality check with readable diffs | `ensure(user).Equals(expectedUser)` |
| `IsTrue()` | Asserts value is `true` | `ensure(isValid).IsTrue()` |
| `IsFalse()` | Asserts value is `false` | `ensure(isValid).IsFalse()` |
| `IsNil()` | Asserts value is `nil` | `ensure(ptr).IsNil()` |
| `IsNotNil()` | Asserts value is not `nil` | `ensure(ptr).IsNotNil()` |
| `IsEmpty()` | Asserts slice, map, or string is empty | `ensure(results).IsEmpty()` |
| `IsNotEmpty()` | Asserts slice, map, or string is not empty | `ensure(results).IsNotEmpty()` |
| `Contains(expected)` | Asserts collection or string contains value | `ensure(list).Contains("item")` |
| `DoesNotContain(exp)` | Asserts collection or string does not contain value | `ensure(list).DoesNotContain("item")` |
| `MatchesRegexp(pat)` | Asserts string matches regex pattern | ``ensure(str).MatchesRegexp(`^v\d+$`)`` |
| `IsError(err)` | Asserts value is an error and matches `err`, or is `nil` if `nil` is expected | `ensure(err).IsError(errNotFound)` |
| `IsNotError()` | Asserts value is `nil` (for errors) | `ensure(err).IsNotError()` |
| `MatchesAllErrors(...)` | Asserts error matches all of the provided errors – useful when asserting errors are wrapped correctly | `ensure(err).MatchesAllErrors(e1, e2)` |

## Subtests
You can run subtests using the `ensure` instance, which handles the shadowing pattern for you. ALWAYS name the callback parameter `ensure` to maintain the shadowing pattern.
- `ensure.Run(name, func(ensure ensuring.E) { ... })`
- `ensure.RunParallel(name, func(ensure ensuring.E) { ... })`

## Table-Driven Tests
`ensure` provides a table runner that auto-initializes mocks and wires them into your subject.

### 1. The Table Structure
To use `ensure.RunTableByIndex`, you MUST follow this specific struct layout for readability and correctness:
* `Name`: Required unique name for the test scenario.
* `Mocks`: Optional struct that defines mocks, as detailed below.
* `SetupMocks`: Optional function that is invoked to allow configuring mock assertions. Supports an optional `ensure ensuring.E` secondary parameter, which you MUST only use if you need to reference `ensure` from within `SetupMocks`.
* `Subject`: Optional struct that can be injected with the `Mocks` dependencies based on fields in the `Mocks` that satisfy interfaces of fields in the `Subject`.

### 2. Running the Table
The callback function for `ensure.RunTableByIndex` MUST name its first parameter `ensure` to maintain the shadowing pattern.

### The Mocks Struct
Mocks must be pointers to structs that MUST have a `NEW(*gomock.Controller) *T` or `NEW() *T` method. Mocks should generally be generated automatically by the `ensure` CLI, as detailed below.

#### Mock Struct Tags
- `` `ensure:"-"` ``: Skips initialization of this field.
- `` `ensure:"ignoreunused"` ``: Prevents "unused mock" errors if the `Subject` doesn't use this mock.

### Example
The following naming, casing, layout, and groupings are conventional. You MUST use this as a template for most table-driven tests, unless explicitly instructed not to.

```go
type Mocks struct {
  DB  *mock_db.MockDB
  API *mock_api.MockAPI
}

table := []struct {
  Name string

  InputVal string

  ExpectedResult string
  ExpectedError  error

  Mocks      *Mocks
  SetupMocks func(*Mocks)
  Subject    *MyService
}{
  {
    Name: "success case",
  
    InputVal: "id",
    
    ExpectedResult: "some data",
    
    SetupMocks: func(m *Mocks) {
      m.DB.EXPECT().Get("id").Return(data, nil)
    },
  },
}

ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
  entry := table[i]
  
  res, err := entry.Subject.DoWork(entry.InputVal)
  ensure(err).IsError(entry.ExpectedError)
  ensure(res).Equals(entry.ExpectedResult)
})
```

## Sync Tests
`ensure.RunTableByIndexSync` and `ensure.RunSync` provide equivalent variants that automatically wrap their blocks with `synctest.Test` from the Go stdlib. 

The `*Sync` variants are useful for testing concurrent code in isolated "bubbles". It is also useful for testing code that leverages `time.Now()`, since time is frozen and only incremented via sleeping, which happens instantly in realtime.

For example, if using the `*Sync` variants with something like: 
```go
start := time.Now()
time.Sleep(10*time.Minute)
result := time.Since(start)
```
The `result` would be 10 minutes, but the test would take microseconds to execute.

## Additional Methods
The `ensure` instance provides helper methods to access the underlying test context or fail the test manually:
- `T()`: Returns the `*testing.T` instance. You MUST use this instead of `t`, to guarantee you're using the correctly scoped `testing.T`.
- `InterfaceT()`: NEVER use this, always use `T()` instead.
- `GoMockController()`: Returns the `*gomock.Controller` scoped to the test.
- `Failf(format, args...)`: Fails the test immediately with a formatted message.

## Generating Mocks
The `ensure mocks generate` CLI command uses the `.ensure.yml` file at the repo root, which is structured as follows:
```yaml
mocks:
  # Used as the directory path relative to the root of the module
  # for any interfaces that are not within internal directories.
  # Optional, defaults to "internal/mocks".
  primaryDestination: internal/mocks

  # Used as the directory path relative to internal directories within the project.
  # Optional, defaults to "mocks".
  internalDestination: mocks

  # Tidy mocks after generation completes.
  # Automatically runs 'ensure mocks tidy' after 'ensure mocks generate' completes.
  # Tidy removes any files that would not be generated by the provided packages list.
  # Optional, defaults to true.
  tidyAfterGenerate: true

  # Disable enhanced matcher failure messages.
  # By default, ensure wraps gomock matchers to provide pretty-printed diffs on failure.
  # Set this to true to revert to standard gomock failure messages.
  # Optional, defaults to false.
  disableEnhancedMatcherFailures: false

  # Packages with interfaces for which to generate mocks
  packages:
    - path: github.com/my/app/some/pkg
      interfaces: [Iface1, Iface2]
```
