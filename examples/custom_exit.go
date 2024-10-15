package main

import (
	"context"
	"fmt"

	"github.com/ZanzyTHEbar/assert-lib/assert"
)

func main() {
	handler := assert.NewAssertHandler()

	// Set a custom exit function (this will prevent the program from exiting)
	handler.SetExitFunc(func(code int) {
		fmt.Printf("Custom exit called with code: %d, but not exiting\n", code)
	})

	// This assertion will trigger the custom exit function but will not exit
	handler.Assert(context.TODO(), false, "Custom Exit Assertion Failed")
}
