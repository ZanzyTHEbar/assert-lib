package main

import (
	"context"
	"os"

	"github.com/ZanzyTHEbar/assert-lib"
)

func main() {
	ctx := context.TODO()

	// === BASIC USAGE - No instance creation required ===
	// These will log errors but won't crash the program by default
	assert.Assert(ctx, 1 == 1, "Basic assertion should pass")
	assert.NotEmpty(ctx, "hello", "String should not be empty")
	assert.Nil(ctx, nil, "This should pass")
	assert.Equal(ctx, 42, 42, "Numbers should be equal")
	assert.True(ctx, true, "Value should be true")

	// This will fail but won't crash the program - just logs the error
	assert.Assert(ctx, false, "This will fail but program continues")

	// === CUSTOM BEHAVIOR WITH OPTIONS ===

	// JSON formatted output
	assert.Assert(ctx, false, "JSON formatted error",
		assert.WithFormatter(&assert.JSONFormatter{}))

	// Silent mode (no output)
	assert.Assert(ctx, false, "Silent assertion",
		assert.WithSilentMode())

	// Custom writer (could be a file, buffer, etc.)
	assert.Assert(ctx, false, "Custom writer",
		assert.WithWriter(os.Stdout))

	// === CONVENIENCE OPTION COMBINATIONS ===

	// Testing defaults (no exit, text format)
	assert.Assert(ctx, false, "Testing environment",
		assert.WithTestingDefaults())

	// Production defaults (no exit, JSON format)
	assert.Assert(ctx, false, "Production environment",
		assert.WithProductionDefaults())

	// === CONDITIONAL CRASH BEHAVIOR ===

	// You can opt into crashing behavior when needed
	// Uncomment to test (will exit the program):
	// assert.Assert(ctx, false, "This will exit the program",
	//     assert.WithCrashOnFailure())

	// Or panic instead of exit
	// assert.Assert(ctx, false, "This will panic",
	//     assert.WithPanicOnFailure())

	// === ADDITIONAL ASSERTION TYPES ===

	assert.Equal(ctx, "expected", "actual", "Strings should match")
	assert.NotEqual(ctx, "foo", "bar", "Strings should be different")
	assert.Contains(ctx, "hello world", "world", "Should contain substring")
	assert.NotContains(ctx, "hello", "xyz", "Should not contain substring")
	assert.False(ctx, false, "Value should be false")

	// === COMBINING OPTIONS ===

	// Multiple options can be combined
	assert.Assert(ctx, false, "Combined options example",
		assert.WithFormatter(&assert.YAMLFormatter{}),
		assert.WithWriter(os.Stdout),
		// Don't exit, just log
	)

	println("Program completed successfully!")
	println("Notice: All assertions above were processed without crashing the program")
	println("This demonstrates the safe-by-default behavior while maintaining flexibility")
}
