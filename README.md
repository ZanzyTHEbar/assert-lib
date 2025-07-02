# Assert Go

A lightweight assertion library for Go, designed for systems and real-time programming with structured logging and flush controls. Features **safe-by-default** behavior with optional instance creation using the Function Options pattern.

[![Go Reference](https://pkg.go.dev/badge/github.com/ZanzyTHEbar/assert-lib.svg)](https://pkg.go.dev/github.com/ZanzyTHEbar/assert-lib)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZanzyTHEbar/assert-lib)](https://goreportcard.com/report/github.com/ZanzyTHEbar/assert-lib)

## Installation

```bash
go get github.com/ZanzyTHEbar/assert-lib
```

## Documentation

-   **API Documentation**: [pkg.go.dev/github.com/ZanzyTHEbar/assert-lib](https://pkg.go.dev/github.com/ZanzyTHEbar/assert-lib)
-   **Examples**: See the [examples/](./examples/) directory
-   **Source Code**: [github.com/ZanzyTHEbar/assert-lib](https://github.com/ZanzyTHEbar/assert-lib)

## Quick Start - No Instance Required!

```go
package main

import (
    "context"
    "github.com/ZanzyTHEbar/assert-lib"
)

func main() {
    ctx := context.TODO()

    // Safe by default - won't crash your program
    assert.Assert(ctx, 1 == 2, "This will log an error but continue")
    assert.NotEmpty(ctx, "hello", "String should not be empty")
    assert.Equal(ctx, 42, 42, "Numbers should match")

    println("Program continues safely!")
}
```

## Advanced Usage with Options

```go
// Custom behavior using Function Options pattern
assert.Assert(ctx, false, "JSON formatted error",
    assert.WithFormatter(&assert.JSONFormatter{}))

// Exit on failure (opt-in behavior)
assert.Assert(ctx, false, "This will exit",
    assert.WithCrashOnFailure())

// Silent mode for performance-critical sections
assert.Assert(ctx, condition, "Silent check",
    assert.WithSilentMode())

// Debug mode with stack traces (great for development)
assert.Assert(ctx, false, "Debug assertion with stack trace",
    assert.WithDebugMode())

// Verbose mode with stack traces and argument details
assert.Assert(ctx, false, "Verbose assertion with full details",
    assert.WithVerboseMode())

// Testing defaults (safe + text format + debug info)
assert.Assert(ctx, false, "Test assertion",
    assert.WithTestingDefaults())

// Production defaults (safe + JSON format + clean output)
assert.Assert(ctx, false, "Production assertion",
    assert.WithProductionDefaults())

// Production defaults (safe + JSON format)
assert.Assert(ctx, false, "Production assertion",
    assert.WithProductionDefaults())
```

## Traditional Instance-Based Usage

```go
handler := assert.NewAssertHandler()
handler.SetExitFunc(func(code int) {
    // Custom exit behavior
})
handler.Assert(context.TODO(), false, "Traditional usage")
```

## Features

-   **🛡️ Safe by Default**: Assertions log errors but don't crash your program by default
-   **⚡ Zero Instance Required**: Use package-level functions without creating handlers
-   **🔧 Function Options Pattern**: Customize behavior on a per-assertion basis
-   **📊 Multiple Formatters**: Text, JSON, YAML output formats
-   **🎯 Rich Assertion Types**: Assert, Equal, NotEmpty, Contains, True/False, Nil checks
-   **🔄 Deferred Assertions**: Batch process multiple assertions
-   **📝 Structured Logging**: Context-based logging with stack traces
-   **🧪 Testing-Friendly**: Built-in testing and production defaults

## Available Assertions

### Basic Assertions

-   `Assert(ctx, condition, msg, opts...)` - Basic truth assertion
-   `True(ctx, value, msg, opts...)` - Assert value is true
-   `False(ctx, value, msg, opts...)` - Assert value is false

### Equality & Comparison

-   `Equal(ctx, expected, actual, msg, opts...)` - Assert values are equal
-   `NotEqual(ctx, expected, actual, msg, opts...)` - Assert values are different

### Nil Checks

-   `Nil(ctx, item, msg, opts...)` - Assert item is nil
-   `NotNil(ctx, item, msg, opts...)` - Assert item is not nil

### String Assertions

-   `NotEmpty(ctx, str, msg, opts...)` - Assert string is not empty
-   `Contains(ctx, str, substr, msg, opts...)` - Assert string contains substring
-   `NotContains(ctx, str, substr, msg, opts...)` - Assert string doesn't contain substring

### Error Handling

-   `NoError(ctx, err, msg, opts...)` - Assert no error occurred
-   `Never(ctx, msg, opts...)` - Always fails (for unreachable code)

## Available Options

### Output Control

-   `WithFormatter(&assert.JSONFormatter{})` - JSON output
-   `WithFormatter(&assert.YAMLFormatter{})` - YAML output
-   `WithWriter(writer)` - Custom output writer
-   `WithSilentMode()` - No output

### Debug & Verbosity Control

-   `WithDebugMode()` - Include stack traces in output (great for development)
-   `WithVerboseMode()` - Include stack traces + argument details
-   `WithQuietMode()` - Minimal output (default behavior)

### Behavior Control

-   `WithCrashOnFailure()` - Exit program on assertion failure
-   `WithPanicOnFailure()` - Panic on assertion failure
-   `WithExitFunc(func(int))` - Custom exit function

### Convenience Presets

-   `WithTestingDefaults()` - Safe behavior + text format + debug info
-   `WithProductionDefaults()` - Safe behavior + JSON format + clean output

## Examples

Run the examples to see different usage patterns:

```bash
go run examples/ergonomic_api.go      # Optimized ergonomic usage
go run examples/debug_modes.go       # Debug/verbose mode examples
go run examples/basic_assertion.go    # Traditional instance-based usage
go run examples/deferred_assertions.go # Batch processing
go run examples/custom_exit.go        # Custom exit behavior
go run examples/formater.go          # Different output formats
```

## Output Modes Comparison

### Default Mode (Clean)

```
ASSERT
   msg=Test assertion
   area=Assert
```

### Debug Mode (With Stack Trace)

```
ASSERT
   msg=Test assertion
   area=Assert
goroutine 1 [running]:
runtime/debug.Stack()...
```

### Verbose Mode (With Stack + Args)

```
ARGS: [key1 value1 key2 value2]
ASSERT
   msg=Test assertion
   area=Assert
goroutine 1 [running]:
runtime/debug.Stack()...
```

### Production Mode (Clean JSON)

```json
{
    "assertData": {
        "area": "Assert",
        "msg": "Test assertion"
    }
}
```

## Philosophy

This library follows a **safe-by-default** philosophy:

-   Package-level functions won't crash your program by default
-   **Stack traces are disabled by default** for clean production logs
-   You opt-in to debug information and crashing behavior when needed
-   You opt-in to crashing behavior when needed
-   Rich context and structured logging help with debugging
-   Function Options pattern provides flexibility without complexity

Perfect for production systems where you want assertions for debugging but can't afford unexpected crashes.

Works very well with my [errbuilder-go](https://github.com/ZanzyTHEbar/errbuilder-go) library.
