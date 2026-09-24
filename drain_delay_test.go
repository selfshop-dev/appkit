package appkit_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/appkit"
)

func TestDrainDelay_Duration(t *testing.T) {
	// Arrange
	delay := appkit.DrainDelay(5 * time.Second)

	// Act
	got := delay.Duration()

	// Assert
	assert.Equal(t, 5*time.Second, got)
}

func TestDrainDelay_String(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name  string
		delay appkit.DrainDelay
		want  string
	}{
		{
			name:  "zero",
			delay: 0,
			want:  "0s",
		},
		{
			name:  "seconds",
			delay: appkit.DrainDelay(5 * time.Second),
			want:  "5s",
		},
		{
			name:  "complex duration",
			delay: appkit.DrainDelay(1*time.Hour + 2*time.Minute + 3*time.Second),
			want:  "1h2m3s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := tc.delay.String()

			// Assert
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDrainDelay_Wait(t *testing.T) {
	t.Run("waits for delay", func(t *testing.T) {
		// Arrange
		delay := appkit.DrainDelay(10 * time.Millisecond)
		ctx := context.Background()

		// Act
		err := delay.Wait(ctx)

		// Assert
		require.NoError(t, err)
	})

	t.Run("returns context error when cancelled", func(t *testing.T) {
		// Arrange
		delay := appkit.DrainDelay(time.Hour)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Act
		err := delay.Wait(ctx)

		// Assert
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("returns nil for zero delay", func(t *testing.T) {
		// Arrange
		delay := appkit.DrainDelay(0)
		ctx := context.Background()

		// Act
		err := delay.Wait(ctx)

		// Assert
		require.NoError(t, err)
	})
}

func TestDrainDelay_Set(t *testing.T) {
	// Arrange
	var delay appkit.DrainDelay

	// Act
	err := delay.Set("5s")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, appkit.DrainDelay(5*time.Second), delay)
}

func TestDrainDelay_Set_Invalid(t *testing.T) {
	// Arrange
	delay := appkit.DrainDelay(5 * time.Second)

	// Act
	err := delay.Set("invalid")

	// Assert
	require.EqualError(
		t,
		err,
		`invalid drain delay "invalid": time: invalid duration "invalid"`,
	)
	assert.Equal(t, appkit.DrainDelay(5*time.Second), delay)
}

func TestDrainDelay_Get(t *testing.T) {
	// Arrange
	delay := appkit.DrainDelay(5 * time.Second)

	// Act
	got := delay.Get()

	// Assert
	assert.Equal(t, delay, got)
}

func TestDrainDelay_MarshalText(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name    string
		delay   appkit.DrainDelay
		want    []byte
		wantErr string
	}{
		{
			name:  "zero",
			delay: 0,
			want:  []byte("0s"),
		},
		{
			name:  "seconds",
			delay: appkit.DrainDelay(5 * time.Second),
			want:  []byte("5s"),
		},
		{
			name:  "complex duration",
			delay: appkit.DrainDelay(1*time.Hour + 2*time.Minute + 3*time.Second),
			want:  []byte("1h2m3s"),
		},
		{
			name:    "negative",
			delay:   appkit.DrainDelay(-5 * time.Second),
			wantErr: "invalid drain delay -5s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := tc.delay.MarshalText()

			// Assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDrainDelay_UnmarshalText(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name string
		text string
		want appkit.DrainDelay
	}{
		{
			name: "empty",
			text: "",
			want: 0,
		},
		{
			name: "seconds",
			text: "5s",
			want: appkit.DrainDelay(5 * time.Second),
		},
		{
			name: "milliseconds",
			text: "250ms",
			want: appkit.DrainDelay(250 * time.Millisecond),
		},
		{
			name: "complex duration",
			text: "1h2m3s",
			want: appkit.DrainDelay(1*time.Hour + 2*time.Minute + 3*time.Second),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var delay appkit.DrainDelay

			// Act
			err := delay.UnmarshalText([]byte(tc.text))

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.want, delay)
		})
	}
}

func TestDrainDelay_UnmarshalText_Invalid(t *testing.T) {
	// Arrange
	delay := appkit.DrainDelay(5 * time.Second)

	// Act
	err := delay.UnmarshalText([]byte("invalid"))

	// Assert
	require.EqualError(
		t,
		err,
		`invalid drain delay "invalid": time: invalid duration "invalid"`,
	)
	assert.Equal(t, appkit.DrainDelay(5*time.Second), delay)
}

func TestDrainDelay_UnmarshalText_Negative(t *testing.T) {
	// Arrange
	delay := appkit.DrainDelay(5 * time.Second)

	// Act
	err := delay.UnmarshalText([]byte("-5s"))

	// Assert
	require.EqualError(
		t,
		err,
		"drain delay must not be negative: -5s",
	)
	assert.Equal(t, appkit.DrainDelay(5*time.Second), delay)
}

func TestDrainDelay_UnmarshalText_NilReceiver(t *testing.T) {
	// Arrange
	var delay *appkit.DrainDelay

	// Act
	err := delay.UnmarshalText([]byte("5s"))

	// Assert
	require.EqualError(t, err, "nil DrainDelay")
}

func TestDrainDelay_JSONMarshal(t *testing.T) {
	// Arrange
	value := struct {
		DrainDelay appkit.DrainDelay `json:"drain_delay"`
	}{
		DrainDelay: appkit.DrainDelay(5 * time.Second),
	}

	// Act
	got, err := json.Marshal(value)

	// Assert
	require.NoError(t, err)
	assert.JSONEq(t, `{"drain_delay":"5s"}`, string(got))
}

func TestDrainDelay_JSONUnmarshal(t *testing.T) {
	// Arrange
	input := []byte(`{"drain_delay":"5s"}`)
	var value struct {
		DrainDelay appkit.DrainDelay `json:"drain_delay"`
	}

	// Act
	err := json.Unmarshal(input, &value)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, appkit.DrainDelay(5*time.Second), value.DrainDelay)
}
