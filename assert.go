package assert

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Option defines a function that can modify assert configuration
type Option func(*AssertConfig)

// AssertConfig holds temporary configuration for assert operations
type AssertConfig struct {
	formatter   Formatter
	writer      io.Writer
	exitFunc    func(code int)
	deferMode   bool
	debugMode   bool
	verboseMode bool
}

// Option functions for configuring assert behavior
func WithFormatter(f Formatter) Option {
	return func(c *AssertConfig) {
		c.formatter = f
	}
}

func WithWriter(w io.Writer) Option {
	return func(c *AssertConfig) {
		c.writer = w
	}
}

func WithExitFunc(f func(int)) Option {
	return func(c *AssertConfig) {
		c.exitFunc = f
	}
}

func WithDeferMode(deferMode bool) Option {
	return func(c *AssertConfig) {
		c.deferMode = deferMode
	}
}

func WithDebugMode() Option {
	return func(c *AssertConfig) {
		c.debugMode = true
	}
}

func WithVerboseMode() Option {
	return func(c *AssertConfig) {
		c.verboseMode = true
	}
}

func WithQuietMode() Option {
	return func(c *AssertConfig) {
		c.debugMode = false
		c.verboseMode = false
	}
}

// Additional convenience options for common use cases
func WithCrashOnFailure() Option {
	return func(c *AssertConfig) {
		c.exitFunc = os.Exit
	}
}

func WithPanicOnFailure() Option {
	return func(c *AssertConfig) {
		c.exitFunc = func(code int) {
			panic(fmt.Sprintf("assertion failed with exit code %d", code))
		}
	}
}

func WithSilentMode() Option {
	return func(c *AssertConfig) {
		c.writer = io.Discard
	}
}

func WithLogLevel(level string) Option {
	return func(c *AssertConfig) {
		// TODO: will be extended to work with different log levels
	}
}

// Combine multiple options into one
func WithTestingDefaults() Option {
	return func(c *AssertConfig) {
		// Set up defaults good for testing
		c.exitFunc = func(code int) {} // No exit
		c.formatter = &TextFormatter{}
		c.debugMode = true // Include stack traces for debugging tests
	}
}

func WithProductionDefaults() Option {
	return func(c *AssertConfig) {
		// Set up defaults good for production
		c.exitFunc = func(code int) {} // No exit, just log
		c.formatter = &JSONFormatter{} // Structured logging
		c.debugMode = false            // No stack traces in production
		c.verboseMode = false          // Concise output
	}
}

// Global default handler for package-level functions
var (
	defaultHandler *AssertHandler
	initOnce       sync.Once
)

// getDefaultHandler returns the global default handler, initializing it if needed
func getDefaultHandler() *AssertHandler {
	initOnce.Do(func() {
		defaultHandler = NewAssertHandler()
		// Set safe defaults that don't crash the program
		defaultHandler.SetExitFunc(func(code int) {
			// No-op: just return instead of exiting
			// Users can opt into crashing behavior if needed
		})
		// Use stderr for default output but don't exit
		defaultHandler.ToWriter(os.Stderr)
	})
	return defaultHandler
}

// applyOptions creates a temporary handler with options applied
func applyOptions(opts ...Option) *AssertHandler {
	handler := getDefaultHandler()

	if len(opts) == 0 {
		return handler
	}

	// Create temporary config
	config := &AssertConfig{
		formatter:   handler.formatter,
		writer:      handler.writer,
		exitFunc:    handler.exitFunc,
		deferMode:   handler.deferAssertions,
		debugMode:   handler.debugMode,
		verboseMode: handler.verboseMode,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create temporary handler with modified config
	tempHandler := &AssertHandler{
		flushes:         handler.flushes,
		assertData:      handler.assertData,
		writer:          config.writer,
		exitFunc:        config.exitFunc,
		formatter:       config.formatter,
		deferredErrors:  []string{},
		deferAssertions: config.deferMode,
		debugMode:       config.debugMode,
		verboseMode:     config.verboseMode,
	}

	return tempHandler
}

// Define the AssertHandler to encapsulate state
type AssertHandler struct {
	flushes         []AssertFlush
	assertData      map[string]AssertData
	writer          io.Writer
	flushLock       sync.Mutex
	exitFunc        func(code int)
	formatter       Formatter
	deferredErrors  []string
	deferAssertions bool
	debugMode       bool
	verboseMode     bool
}

// Define interfaces for logging/asserting
type AssertData interface {
	Dump() string
}

type AssertFlush interface {
	Flush()
}

// Add a constructor for the handler
func NewAssertHandler() *AssertHandler {
	return &AssertHandler{
		flushes:         []AssertFlush{},
		assertData:      make(map[string]AssertData),
		writer:          os.Stderr,
		exitFunc:        os.Exit,          // Default exit behavior
		formatter:       &TextFormatter{}, // Default to text formatter
		deferredErrors:  []string{},
		deferAssertions: false,
		debugMode:       false, // Default: no stack traces
		verboseMode:     false, // Default: concise output
	}
}

// SetDeferAssertions allows toggling deferred assertion mode
func (a *AssertHandler) SetDeferAssertions(deferMode bool) {
	a.deferAssertions = deferMode
}

func (a *AssertHandler) SetFormatter(formatter Formatter) {
	a.formatter = formatter
}

func (a *AssertHandler) SetExitFunc(exitFunc func(int)) {
	a.exitFunc = exitFunc
}

func (a *AssertHandler) AddAssertData(key string, value AssertData) {
	a.assertData[key] = value
}

func (a *AssertHandler) RemoveAssertData(key string) {
	delete(a.assertData, key)
}

func (a *AssertHandler) AddAssertFlush(flusher AssertFlush) {
	a.flushes = append(a.flushes, flusher)
}

func (a *AssertHandler) ToWriter(w io.Writer) {
	a.writer = w
}

func (a *AssertHandler) SetDebugMode(debugMode bool) {
	a.debugMode = debugMode
}

func (a *AssertHandler) SetVerboseMode(verboseMode bool) {
	a.verboseMode = verboseMode
}

func (a *AssertHandler) runAssert(ctx context.Context, msg string, args ...interface{}) {
	a.flushLock.Lock()
	defer a.flushLock.Unlock()

	// Check if the context has been canceled
	if err := ctx.Err(); err != nil {
		fmt.Fprintln(a.writer, "Context canceled:", err)
		return
	}

	// Prevent re-entrancy by skipping further flushes
	for _, f := range a.flushes {
		f.Flush()
	}

	data := map[string]interface{}{
		"msg":  msg,
		"area": "Assert",
	}

	// append the args to the data
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		data[args[i].(string)] = args[i+1]
	}

	// Only print ARGS in verbose mode
	if a.verboseMode {
		fmt.Fprintf(a.writer, "ARGS: %+v\n", args)
	}

	for k, v := range a.assertData {
		data[k] = v.Dump()
	}

	var stack string
	// Only capture stack trace in debug or verbose mode
	if a.debugMode || a.verboseMode {
		stack = string(debug.Stack())
	}

	formattedOutput := a.formatter.Format(data, stack)

	fmt.Fprintln(a.writer, "ASSERT")
	fmt.Fprintln(a.writer, formattedOutput)

	// If we are in deferred mode, store the error and return
	if a.deferAssertions {
		a.deferredErrors = append(a.deferredErrors, formattedOutput)
		return
	}

	// Use the custom exit function instead of os.Exit directly
	a.exitFunc(1)
}

// Process all deferred assertions at once, logging or exiting if needed
func (a *AssertHandler) ProcessDeferredAssertions(ctx context.Context) {
	if len(a.deferredErrors) > 0 {
		// Combine all errors into a single string
		combinedErrors := strings.Join(a.deferredErrors, "\n---\n")
		fmt.Fprintln(a.writer, combinedErrors)

		// Clear the deferred errors after processing
		a.deferredErrors = []string{}

		// Exit after processing if it's an ERROR level
		a.exitFunc(1)
	}
}

func (a *AssertHandler) Assert(ctx context.Context, truth bool, msg string, data ...any) {
	if !truth {
		a.runAssert(ctx, msg, data...)
	}
}

func (a *AssertHandler) AssertWithTimeout(ctx context.Context, timeout time.Duration, truth bool, msg string, data ...any) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if !truth {
		a.runAssert(ctx, msg, data...)
	}
}

func (a *AssertHandler) Nil(ctx context.Context, item any, msg string, data ...any) {
	if item != nil {
		slog.ErrorContext(ctx, "Nil#not nil encountered")
		a.runAssert(ctx, msg, data...)
	}
}

func (a *AssertHandler) NotNil(ctx context.Context, item any, msg string, data ...any) {
	if item == nil {
		slog.ErrorContext(ctx, "NotNil#nil encountered")
		a.runAssert(ctx, msg, data...)
	}
}

func (a *AssertHandler) Never(ctx context.Context, msg string, data ...any) {
	a.runAssert(ctx, msg, data...)
}

func (a *AssertHandler) NoError(ctx context.Context, err error, msg string, data ...any) {
	if err != nil {
		data = append(data, "error", err)
		a.runAssert(ctx, msg, data...)
	}
}

// Package-level functions for ergonomic usage

// Assert checks if the condition is true, failing if false
func Assert(ctx context.Context, truth bool, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.Assert(ctx, truth, msg)
}

// AssertWithTimeout checks if the condition is true with a timeout, failing if false
func AssertWithTimeout(ctx context.Context, timeout time.Duration, truth bool, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.AssertWithTimeout(ctx, timeout, truth, msg)
}

// Nil checks if the item is nil, failing if not nil
func Nil(ctx context.Context, item any, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.Nil(ctx, item, msg)
}

// NotNil checks if the item is not nil, failing if nil
func NotNil(ctx context.Context, item any, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.NotNil(ctx, item, msg)
}

// Never always fails - use for code paths that should never be reached
func Never(ctx context.Context, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.Never(ctx, msg)
}

// NoError checks if the error is nil, failing if not nil
func NoError(ctx context.Context, err error, msg string, opts ...Option) {
	handler := applyOptions(opts...)
	handler.NoError(ctx, err, msg)
}

// NotEmpty checks if a string is not empty, failing if empty
func NotEmpty(ctx context.Context, str string, msg string, opts ...Option) {
	if str == "" {
		// String is empty, assertion should fail
		handler := applyOptions(opts...)
		handler.Assert(ctx, false, msg, "value", str)
	}
	// String is not empty, assertion passes (no action needed)
}

// Equal checks if two values are equal, failing if they're not
func Equal(ctx context.Context, expected, actual any, msg string, opts ...Option) {
	if expected == actual {
		return // Values are equal, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "expected", expected, "actual", actual)
}

// NotEqual checks if two values are not equal, failing if they are equal
func NotEqual(ctx context.Context, expected, actual any, msg string, opts ...Option) {
	if expected != actual {
		return // Values are not equal, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "expected_not", expected, "actual", actual)
}

// Contains checks if a string contains a substring, failing if it doesn't
func Contains(ctx context.Context, str, substr string, msg string, opts ...Option) {
	if strings.Contains(str, substr) {
		return // String contains substring, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "string", str, "substring", substr)
}

// NotContains checks if a string does not contain a substring, failing if it does
func NotContains(ctx context.Context, str, substr string, msg string, opts ...Option) {
	if !strings.Contains(str, substr) {
		return // String does not contain substring, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "string", str, "substring", substr)
}

// True checks if a value is true, failing if false
func True(ctx context.Context, value bool, msg string, opts ...Option) {
	if value {
		return // Value is true, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "value", value)
}

// False checks if a value is false, failing if true
func False(ctx context.Context, value bool, msg string, opts ...Option) {
	if !value {
		return // Value is false, assertion passes
	}
	handler := applyOptions(opts...)
	handler.Assert(ctx, false, msg, "value", value)
}

// ProcessDeferredAssertions processes any deferred assertions on the default handler
func ProcessDeferredAssertions(ctx context.Context, opts ...Option) {
	handler := applyOptions(opts...)
	handler.ProcessDeferredAssertions(ctx)
}
