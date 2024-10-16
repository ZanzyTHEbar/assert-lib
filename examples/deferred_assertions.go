package main

import (
	"context"

	"github.com/ZanzyTHEbar/assert-lib"
)

func main() {
	handler := assert.NewAssertHandler()

	// Enable deferred assertions
	handler.SetDeferAssertions(true)

	// Multiple assertions that will not fail immediately
	handler.Nil(context.TODO(), nil, "Deferred Nil Assertion")
	handler.Assert(context.TODO(), false, "Deferred Assert Failure")

	// Process all deferred assertions (will print all errors and exit)
	handler.ProcessDeferredAssertions(context.TODO())
}
