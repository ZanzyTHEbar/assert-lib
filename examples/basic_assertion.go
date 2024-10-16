package main

import (
	"context"

	"github.com/ZanzyTHEbar/assert-lib"
)

func main() {
	handler := assert.NewAssertHandler()

	// Basic assertion (this will fail and exit)
	handler.Assert(context.TODO(), 2 == 1, "Basic Assertion Failed: 2 is not equal to 1")
}
