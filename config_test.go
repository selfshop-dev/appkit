package appkit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/selfshop-dev/appkit"
)

func TestConfig_IsDevmod(t *testing.T) {
	t.Run("returns true for dev mode", func(t *testing.T) {
		// Arrange
		config := appkit.Config{
			Runmode: appkit.RunmodeDev,
		}

		// Act
		result := config.IsDevmod()

		// Assert
		assert.True(t, result)
	})

	t.Run("returns false for non-dev mode", func(t *testing.T) {
		// Arrange
		config := appkit.Config{
			Runmode: appkit.RunmodeProd,
		}

		// Act
		result := config.IsDevmod()

		// Assert
		assert.False(t, result)
	})
}

func TestShutdownConfig_TotalTimeout(t *testing.T) {
	t.Run("returns sum of timeout and drain delay", func(t *testing.T) {
		// Arrange
		config := appkit.ShutdownConfig{
			Timeout:    5 * time.Second,
			DrainDelay: appkit.DrainDelay(2 * time.Second),
		}

		// Act
		result := config.TotalTimeout()

		// Assert
		assert.Equal(t, 7*time.Second, result)
	})

	t.Run("returns timeout when drain delay is zero", func(t *testing.T) {
		// Arrange
		config := appkit.ShutdownConfig{
			Timeout:    5 * time.Second,
			DrainDelay: 0,
		}

		// Act
		result := config.TotalTimeout()

		// Assert
		assert.Equal(t, 5*time.Second, result)
	})

	t.Run("returns drain delay when timeout is zero", func(t *testing.T) {
		// Arrange
		config := appkit.ShutdownConfig{
			Timeout:    0,
			DrainDelay: appkit.DrainDelay(2 * time.Second),
		}

		// Act
		result := config.TotalTimeout()

		// Assert
		assert.Equal(t, 2*time.Second, result)
	})

	t.Run("returns zero for empty config", func(t *testing.T) {
		// Arrange
		var config appkit.ShutdownConfig

		// Act
		result := config.TotalTimeout()

		// Assert
		assert.Zero(t, result)
	})
}
