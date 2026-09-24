package appkit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanupStack_Add(t *testing.T) {
	t.Run("adds cleanup item", func(t *testing.T) {
		// Arrange
		var stack CleanupStack
		fn := func(context.Context) error {
			return nil
		}

		// Act
		stack.Add(" cleanup ", fn)

		// Assert
		require.Len(t, stack.items, 1)
		assert.Equal(t, "cleanup", stack.items[0].name)
		assert.NotNil(t, stack.items[0].fn)
	})

	t.Run("uses unnamed for blank name", func(t *testing.T) {
		// Arrange
		var stack CleanupStack
		fn := func(context.Context) error {
			return nil
		}

		// Act
		stack.Add(" \t\n ", fn)

		// Assert
		require.Len(t, stack.items, 1)
		assert.Equal(t, "<unnamed>", stack.items[0].name)
		assert.NotNil(t, stack.items[0].fn)
	})

	t.Run("ignores nil function", func(t *testing.T) {
		// Arrange
		var stack CleanupStack

		// Act
		stack.Add("cleanup", nil)

		// Assert
		assert.Empty(t, stack.items)
	})
}

func TestCleanupStack_Finalize_Empty(t *testing.T) {
	// Arrange
	var stack CleanupStack

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.NoError(t, err)
}

func TestCleanupStack_Finalize_Order(t *testing.T) {
	// Arrange
	var stack CleanupStack
	var calls []string

	stack.Add("first", func(context.Context) error {
		calls = append(calls, "first")
		return nil
	})
	stack.Add("second", func(context.Context) error {
		calls = append(calls, "second")
		return nil
	})
	stack.Add("third", func(context.Context) error {
		calls = append(calls, "third")
		return nil
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{
		"third",
		"second",
		"first",
	}, calls)
}

func TestCleanupStack_Finalize_PassesContext(t *testing.T) {
	// Arrange
	var stack CleanupStack
	type ctxKey string

	const testKey ctxKey = "key"
	ctx := context.WithValue(context.Background(), testKey, "value")
	var received context.Context

	stack.Add("cleanup", func(ctx context.Context) error {
		received = ctx
		return nil
	})

	// Act
	err := stack.Finalize(ctx)

	// Assert
	require.NoError(t, err)
	assert.Same(t, ctx, received)
}

func TestCleanupStack_Finalize_Error(t *testing.T) {
	// Arrange
	var stack CleanupStack
	cleanupErr := errors.New("cleanup failed")

	stack.Add("database", func(context.Context) error {
		return cleanupErr
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.Error(t, err)
	assert.Equal(t, "database: cleanup failed", err.Error())
	assert.ErrorIs(t, err, cleanupErr)
}

func TestCleanupStack_Finalize_ContinuesAfterError(t *testing.T) {
	// Arrange
	var stack CleanupStack
	var calls []string
	cleanupErr := errors.New("cleanup failed")

	stack.Add("first", func(context.Context) error {
		calls = append(calls, "first")
		return cleanupErr
	})
	stack.Add("second", func(context.Context) error {
		calls = append(calls, "second")
		return nil
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, cleanupErr)
	assert.Equal(t, []string{
		"second",
		"first",
	}, calls)
}

func TestCleanupStack_Finalize_AggregatesErrors(t *testing.T) {
	// Arrange
	var stack CleanupStack
	firstErr := errors.New("first error")
	secondErr := errors.New("second error")

	stack.Add("first", func(context.Context) error {
		return firstErr
	})
	stack.Add("second", func(context.Context) error {
		return secondErr
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, firstErr)
	require.ErrorIs(t, err, secondErr)
	assert.Equal(
		t,
		"second: second error\nfirst: first error",
		err.Error(),
	)
}

func TestCleanupStack_Finalize_ClearsItems(t *testing.T) {
	// Arrange
	var stack CleanupStack

	stack.Add("first", func(context.Context) error {
		return nil
	})
	stack.Add("second", func(context.Context) error {
		return nil
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Nil(t, stack.items)
}

func TestCleanupStack_Finalize_ClearsItemsAfterError(t *testing.T) {
	// Arrange
	var stack CleanupStack
	cleanupErr := errors.New("cleanup failed")

	stack.Add("cleanup", func(context.Context) error {
		return cleanupErr
	})

	// Act
	err := stack.Finalize(context.Background())

	// Assert
	require.Error(t, err)
	assert.Nil(t, stack.items)
}

func TestCleanupStack_Finalize_CalledTwice(t *testing.T) {
	// Arrange
	var stack CleanupStack
	var calls int

	stack.Add("cleanup", func(context.Context) error {
		calls++
		return nil
	})

	// Act
	firstErr := stack.Finalize(context.Background())
	secondErr := stack.Finalize(context.Background())

	// Assert
	require.NoError(t, firstErr)
	require.NoError(t, secondErr)
	assert.Equal(t, 1, calls)
}
