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

	fmt.Fprintf(a.writer, "ARGS: %+v\n", args)

	for k, v := range a.assertData {
		data[k] = v.Dump()
	}

	stack := string(debug.Stack())

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
	slog.InfoContext(ctx, "Nil Check", "item", item)
	if item == nil {
		a.runAssert(ctx, msg, data...)
		return
	}

	slog.ErrorContext(ctx, "Nil#not nil encountered")
	a.runAssert(ctx, msg, data...)
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
