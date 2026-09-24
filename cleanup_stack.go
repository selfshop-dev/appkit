package appkit

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type cleanupItem struct {
	name string
	fn   func(context.Context) error
}

// CleanupStack stores cleanup functions and executes them in LIFO order.
//
// Cleanup functions are executed by [CleanupStack.Finalize]. All registered
// functions are executed even if one or more functions return an error.
// Multiple errors are combined and can be inspected with [errors.Is] and
// [errors.As].
type CleanupStack struct {
	items []cleanupItem
}

// Add registers a cleanup function.
//
// Cleanup functions are executed in reverse order of registration. A nil
// function is ignored. The cleanup name is trimmed; if it is empty,
// "<unnamed>" is used instead.
func (c *CleanupStack) Add(name string, fn func(context.Context) error) {
	if fn == nil {
		return
	}
	if name = strings.TrimSpace(name); name == "" {
		name = "<unnamed>"
	}
	c.items = append(c.items, cleanupItem{name: name, fn: fn})
}

// Finalize executes all registered cleanup functions in reverse order.
//
// All cleanup functions are executed even when one of them returns an error.
// Returned errors are wrapped with the corresponding cleanup name and
// combined into a single error.
//
// After Finalize returns, the cleanup stack is empty.
func (c *CleanupStack) Finalize(ctx context.Context) error {
	var err error
	for _, item := range slices.Backward(c.items) {
		if itemErr := item.fn(ctx); itemErr != nil {
			err = errors.Join(
				err,
				fmt.Errorf("%s: %w", item.name, itemErr),
			)
		}
	}
	c.items = nil
	return err
}
