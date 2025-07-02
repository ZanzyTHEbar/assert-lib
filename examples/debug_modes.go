package main

import (
	"context"
	"os"

	"github.com/ZanzyTHEbar/assert-lib"
)

func main() {
	ctx := context.TODO()

	println("=== DEFAULT MODE (No Stack Traces) ===")
	assert.Assert(ctx, false, "Default mode - clean output")

	println("\n=== DEBUG MODE (With Stack Traces) ===")
	assert.Assert(ctx, false, "Debug mode - includes stack trace",
		assert.WithDebugMode())

	println("\n=== VERBOSE MODE (With Stack + Args) ===")
	assert.Assert(ctx, false, "Verbose mode - includes everything",
		assert.WithVerboseMode())

	println("\n=== QUIET MODE (Minimal Output) ===")
	assert.Assert(ctx, false, "Quiet mode - minimal output",
		assert.WithQuietMode())

	println("\n=== TESTING DEFAULTS (Debug Mode) ===")
	assert.Assert(ctx, false, "Testing defaults include debug info",
		assert.WithTestingDefaults())

	println("\n=== PRODUCTION DEFAULTS (Clean JSON) ===")
	assert.Assert(ctx, false, "Production defaults are clean",
		assert.WithProductionDefaults())

	println("\n=== CUSTOM COMBINATIONS ===")
	assert.Assert(ctx, false, "JSON format with debug info",
		assert.WithFormatter(&assert.JSONFormatter{}),
		assert.WithDebugMode())

	// Test with custom writer
	println("\n=== CUSTOM WRITER + VERBOSE ===")
	assert.Assert(ctx, false, "Custom writer with verbose mode",
		assert.WithWriter(os.Stdout),
		assert.WithVerboseMode())

	println("\nProgram completed - notice the different output levels!")
}
