# Assert Go

A lightweight assertion library for Go, designed for systems and real-time programming with structured logging and flush controls.

## Installation

```bash
go get github.com/ZanzyTHEbar/assert-lib
```

## Usage

```go
package main

import (
    "context"
    "github.com/ZanzyTHEbar/assert-lib/assert"
)

func main() {
    handler := assert.NewAssertHandler()
    
    handler.Assert(context.TODO(), false, "This should fail")
}
```

Check out the [examples](/examples/) directory for usage examples.

## Features

- **Assertions**: Assert, Nil, NotNil, NoError, Never.
- **Flush Management**: Control output flushes with AssertFlush.
- **Context-Based Logging**: Attach structured logging to your assertion calls.
- **Custom Loggers**: Use your own logger with the AssertHandler interface.

## Examples

```bash
go run examples/basic_assertion.go
go run examples/deferred_assertions.go
go run examples/custom_exit.go
go run examples/formater.go
```
