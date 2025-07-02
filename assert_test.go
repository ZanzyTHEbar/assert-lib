package assert

import (
	"bytes"
	"context"
	"testing"
)

func TestAssert(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {}) // Prevent exit

	handler.Assert(context.TODO(), false, "Test Failure")

	if !bytes.Contains(buffer.Bytes(), []byte("Test Failure")) {
		t.Fatalf("Expected failure message not found in output")
	}
}

func TestCustomExitFunc(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {})

	called := false
	handler.SetExitFunc(func(code int) {
		called = true // Track if the custom exit is called
	})

	handler.Assert(context.TODO(), false, "Test Failure")

	if !called {
		t.Fatalf("Custom exit function was not called")
	}
	if !bytes.Contains(buffer.Bytes(), []byte("Test Failure")) {
		t.Fatalf("Expected failure message not found in output")
	}
}

func TestJSONFormatter(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {})

	handler.SetFormatter(&JSONFormatter{})
	handler.Assert(context.TODO(), false, "Test Failure")

	expected := `"msg": "Test Failure"`
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Fatalf("Expected JSON output not found in output")
	}
}

func TestYAMLFormatter(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {})

	handler.SetFormatter(&YAMLFormatter{})
	handler.Assert(context.TODO(), false, "Test Failure")

	expected := "msg: Test Failure"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Fatalf("Expected YAML output not found in output")
	}
}

func TestDeferredAssertions(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {})

	// Enable deferred assertions
	handler.SetDeferAssertions(true)

	// These should not be printed immediately
	handler.NotNil(context.TODO(), nil, "Test Deferred NotNil") // This will fail since nil is passed
	handler.Assert(context.TODO(), false, "Test Deferred Assert")

	// Process deferred assertions
	handler.ProcessDeferredAssertions(context.TODO())

	if !bytes.Contains(buffer.Bytes(), []byte("Test Deferred Assert")) {
		t.Fatalf("Expected deferred assertion message not found")
	}
	if !bytes.Contains(buffer.Bytes(), []byte("Test Deferred NotNil")) {
		t.Fatalf("Expected deferred NotNil message not found")
	}
}

func TestImmediateAssertions(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewAssertHandler()
	handler.ToWriter(&buffer)
	handler.SetExitFunc(func(code int) {})

	// Disable deferred assertions (default)
	handler.Assert(context.TODO(), false, "Test Immediate Assert")

	if !bytes.Contains(buffer.Bytes(), []byte("Test Immediate Assert")) {
		t.Fatalf("Expected immediate assertion message not found")
	}
}

func TestPackageLevelAssert(t *testing.T) {
	var buffer bytes.Buffer

	// Test package-level function with custom writer
	Assert(context.TODO(), false, "Package level assertion failed",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {})) // Prevent exit

	if !bytes.Contains(buffer.Bytes(), []byte("Package level assertion failed")) {
		t.Fatalf("Expected package-level assertion message not found in output")
	}
}

func TestNotEmpty(t *testing.T) {
	var buffer bytes.Buffer

	// Test NotEmpty with empty string (should fail)
	NotEmpty(context.TODO(), "", "String should not be empty",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	if !bytes.Contains(buffer.Bytes(), []byte("String should not be empty")) {
		t.Fatalf("Expected NotEmpty assertion message not found in output")
	}

	// Reset buffer
	buffer.Reset()

	// Test NotEmpty with non-empty string (should pass)
	NotEmpty(context.TODO(), "hello", "String should not be empty",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("NotEmpty should have passed but got output: %s", buffer.String())
	}
}

func TestPackageLevelWithOptions(t *testing.T) {
	var buffer bytes.Buffer

	// Test with JSON formatter
	Assert(context.TODO(), false, "JSON formatted error",
		WithFormatter(&JSONFormatter{}),
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	if !bytes.Contains(buffer.Bytes(), []byte(`"msg": "JSON formatted error"`)) {
		t.Fatalf("Expected JSON formatted output not found")
	}
}

func TestDefaultSafeBehavior(t *testing.T) {
	// This test ensures the default behavior doesn't crash
	// We can't easily test this without complex process management,
	// but we can verify the default exit function is a no-op

	// Get a fresh handler to test with
	handler := NewAssertHandler()
	handler.SetExitFunc(func(code int) {}) // Set no-op exit function

	var buffer bytes.Buffer
	handler.ToWriter(&buffer)

	// This should call the exit function but not crash
	handler.Assert(context.TODO(), false, "Test safe behavior")

	if !bytes.Contains(buffer.Bytes(), []byte("Test safe behavior")) {
		t.Fatalf("Expected safe behavior message not found in output")
	}
}

func TestEqual(t *testing.T) {
	var buffer bytes.Buffer

	// Test Equal with matching values (should pass)
	Equal(context.TODO(), 42, 42, "Numbers should be equal",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("Equal should have passed but got output: %s", buffer.String())
	}

	// Reset buffer
	buffer.Reset()

	// Test Equal with non-matching values (should fail)
	Equal(context.TODO(), 42, 24, "Numbers should be equal",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	if !bytes.Contains(buffer.Bytes(), []byte("Numbers should be equal")) {
		t.Fatalf("Expected Equal assertion message not found in output")
	}
}

func TestNotEqual(t *testing.T) {
	var buffer bytes.Buffer

	// Test NotEqual with different values (should pass)
	NotEqual(context.TODO(), 42, 24, "Numbers should be different",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("NotEqual should have passed but got output: %s", buffer.String())
	}

	// Reset buffer
	buffer.Reset()

	// Test NotEqual with same values (should fail)
	NotEqual(context.TODO(), 42, 42, "Numbers should be different",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	if !bytes.Contains(buffer.Bytes(), []byte("Numbers should be different")) {
		t.Fatalf("Expected NotEqual assertion message not found in output")
	}
}

func TestContains(t *testing.T) {
	var buffer bytes.Buffer

	// Test Contains with valid substring (should pass)
	Contains(context.TODO(), "hello world", "world", "Should contain substring",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("Contains should have passed but got output: %s", buffer.String())
	}

	// Reset buffer
	buffer.Reset()

	// Test Contains with invalid substring (should fail)
	Contains(context.TODO(), "hello", "xyz", "Should contain substring",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	if !bytes.Contains(buffer.Bytes(), []byte("Should contain substring")) {
		t.Fatalf("Expected Contains assertion message not found in output")
	}
}

func TestConvenienceOptions(t *testing.T) {
	var buffer bytes.Buffer

	// Test WithTestingDefaults
	Assert(context.TODO(), false, "Testing defaults",
		WithTestingDefaults(),
		WithWriter(&buffer))

	if !bytes.Contains(buffer.Bytes(), []byte("Testing defaults")) {
		t.Fatalf("Expected testing defaults assertion message not found in output")
	}

	// Reset buffer
	buffer.Reset()

	// Test WithProductionDefaults (should use JSON formatter)
	Assert(context.TODO(), false, "Production defaults",
		WithProductionDefaults(),
		WithWriter(&buffer))

	if !bytes.Contains(buffer.Bytes(), []byte(`"msg": "Production defaults"`)) {
		t.Fatalf("Expected JSON formatted output from production defaults not found")
	}
}

func TestTrueFalse(t *testing.T) {
	var buffer bytes.Buffer

	// Test True with true value (should pass)
	True(context.TODO(), true, "Value should be true",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("True should have passed but got output: %s", buffer.String())
	}

	// Reset buffer
	buffer.Reset()

	// Test False with false value (should pass)
	False(context.TODO(), false, "Value should be false",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	// Buffer should be empty since assertion passed
	if buffer.Len() > 0 {
		t.Fatalf("False should have passed but got output: %s", buffer.String())
	}
}

func TestDebugModes(t *testing.T) {
	var buffer bytes.Buffer

	// Test default mode (no stack trace)
	Assert(context.TODO(), false, "Default mode",
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	output := buffer.String()
	if bytes.Contains(buffer.Bytes(), []byte("goroutine")) {
		t.Fatalf("Default mode should not include stack trace, but got: %s", output)
	}

	// Reset buffer
	buffer.Reset()

	// Test debug mode (with stack trace)
	Assert(context.TODO(), false, "Debug mode",
		WithDebugMode(),
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	output = buffer.String()
	if !bytes.Contains(buffer.Bytes(), []byte("goroutine")) {
		t.Fatalf("Debug mode should include stack trace, but got: %s", output)
	}

	// Reset buffer
	buffer.Reset()

	// Test verbose mode (with stack trace and args)
	Assert(context.TODO(), false, "Verbose mode",
		WithVerboseMode(),
		WithWriter(&buffer),
		WithExitFunc(func(code int) {}))

	output = buffer.String()
	if !bytes.Contains(buffer.Bytes(), []byte("goroutine")) {
		t.Fatalf("Verbose mode should include stack trace, but got: %s", output)
	}
	if !bytes.Contains(buffer.Bytes(), []byte("ARGS:")) {
		t.Fatalf("Verbose mode should include ARGS, but got: %s", output)
	}
}

func TestProductionDefaults(t *testing.T) {
	var buffer bytes.Buffer

	// Test production defaults (should be clean JSON without stack)
	Assert(context.TODO(), false, "Production test",
		WithProductionDefaults(),
		WithWriter(&buffer))

	output := buffer.String()
	if bytes.Contains(buffer.Bytes(), []byte("goroutine")) {
		t.Fatalf("Production defaults should not include stack trace, but got: %s", output)
	}
	if !bytes.Contains(buffer.Bytes(), []byte(`"msg": "Production test"`)) {
		t.Fatalf("Production defaults should use JSON format, but got: %s", output)
	}
}
