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
	handler.Nil(context.TODO(), nil, "Test Deferred Nil")
	handler.Assert(context.TODO(), false, "Test Deferred Assert")

	// Process deferred assertions
	handler.ProcessDeferredAssertions(context.TODO())

	if !bytes.Contains(buffer.Bytes(), []byte("Test Deferred Assert")) {
		t.Fatalf("Expected deferred assertion message not found")
	}
	if !bytes.Contains(buffer.Bytes(), []byte("Test Deferred Nil")) {
		t.Fatalf("Expected deferred nil message not found")
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
