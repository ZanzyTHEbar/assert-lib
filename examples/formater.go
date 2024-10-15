package main

import (
	"context"

	"github.com/ZanzyTHEbar/assert-lib/assert"
)

func main() {
	handler := assert.NewAssertHandler()

	// Set the formatter to JSON
	handler.SetFormatter(&assert.JSONFormatter{})

	// This will output the assertion failure in JSON format
	handler.Assert(context.TODO(), false, "JSON Format Assertion Failed")

	// Set the formatter to YAML
	handler.SetFormatter(&assert.YAMLFormatter{})

	// This will output the assertion failure in YAML format
	handler.Assert(context.TODO(), false, "YAML Format Assertion Failed")
}
