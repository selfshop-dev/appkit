package appkit_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/appkit"
)

func TestParseRunmode(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name    string
		text    string
		want    appkit.Runmode
		wantErr string
	}{
		{
			name: "dev",
			text: "dev",
			want: appkit.RunmodeDev,
		},
		{
			name: "prod",
			text: "prod",
			want: appkit.RunmodeProd,
		},
		{
			name: "empty",
			text: "",
			want: appkit.RunmodeProd,
		},
		{
			name: "case insensitive",
			text: "DEV",
			want: appkit.RunmodeDev,
		},
		{
			name:    "invalid",
			text:    "invalid",
			want:    appkit.InvalidRunmode,
			wantErr: `unrecognized runmode: "invalid"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := appkit.ParseRunmode(tc.text)

			// Assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				assert.Equal(t, tc.want, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRunmode_String(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name string
		r    appkit.Runmode
		want string
	}{
		{
			name: "dev",
			r:    appkit.RunmodeDev,
			want: "dev",
		},
		{
			name: "prod",
			r:    appkit.RunmodeProd,
			want: "prod",
		},
		{
			name: "invalid",
			r:    appkit.InvalidRunmode,
			want: "invalid",
		},
		{
			name: "unknown",
			r:    appkit.Runmode(100),
			want: "unknown(100)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := tc.r.String()

			// Assert
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRunmode_Set(t *testing.T) {
	// Arrange
	var r appkit.Runmode

	// Act
	err := r.Set("dev")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, appkit.RunmodeDev, r)
}

func TestRunmode_Set_Invalid(t *testing.T) {
	// Arrange
	r := appkit.RunmodeDev

	// Act
	err := r.Set("invalid")

	// Assert
	require.EqualError(t, err, `unrecognized runmode: "invalid"`)
	assert.Equal(t, appkit.InvalidRunmode, r)
}

func TestRunmode_Get(t *testing.T) {
	// Arrange
	r := appkit.RunmodeDev

	// Act
	got := r.Get()

	// Assert
	assert.Equal(t, r, got)
}

func TestRunmode_MarshalText(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name    string
		r       appkit.Runmode
		want    []byte
		wantErr string
	}{
		{
			name: "dev",
			r:    appkit.RunmodeDev,
			want: []byte("dev"),
		},
		{
			name: "prod",
			r:    appkit.RunmodeProd,
			want: []byte("prod"),
		},
		{
			name:    "invalid",
			r:       appkit.InvalidRunmode,
			wantErr: "invalid runmode 1",
		},
		{
			name:    "unknown",
			r:       appkit.Runmode(100),
			wantErr: "invalid runmode 100",
		},
		{
			name:    "less than minimum",
			r:       appkit.Runmode(-2),
			wantErr: "invalid runmode -2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := tc.r.MarshalText()

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

func TestRunmode_UnmarshalText(t *testing.T) {
	// Arrange
	testCases := [...]struct {
		name string
		text string
		want appkit.Runmode
	}{
		{
			name: "dev",
			text: "dev",
			want: appkit.RunmodeDev,
		},
		{
			name: "prod",
			text: "prod",
			want: appkit.RunmodeProd,
		},
		{
			name: "empty",
			text: "",
			want: appkit.RunmodeProd,
		},
		{
			name: "upper case",
			text: "PROD",
			want: appkit.RunmodeProd,
		},
		{
			name: "mixed case",
			text: "DeV",
			want: appkit.RunmodeDev,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var r appkit.Runmode

			// Act
			err := r.UnmarshalText([]byte(tc.text))

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.want, r)
		})
	}
}

func TestRunmode_UnmarshalText_Invalid(t *testing.T) {
	// Arrange
	r := appkit.RunmodeDev

	// Act
	err := r.UnmarshalText([]byte("invalid"))

	// Assert
	require.EqualError(t, err, `unrecognized runmode: "invalid"`)
	assert.Equal(t, appkit.InvalidRunmode, r)
}

func TestRunmode_UnmarshalText_NilReceiver(t *testing.T) {
	// Arrange
	var r *appkit.Runmode

	// Act
	err := r.UnmarshalText([]byte("dev"))

	// Assert
	require.EqualError(t, err, "nil Runmode")
}

func TestRunmode_JSONMarshal(t *testing.T) {
	// Arrange
	value := struct {
		Runmode appkit.Runmode `json:"runmode"`
	}{
		Runmode: appkit.RunmodeDev,
	}

	// Act
	got, err := json.Marshal(value)

	// Assert
	require.NoError(t, err)
	assert.JSONEq(t, `{"runmode":"dev"}`, string(got))
}

func TestRunmode_JSONUnmarshal(t *testing.T) {
	// Arrange
	input := []byte(`{"runmode":"prod"}`)
	var value struct {
		Runmode appkit.Runmode `json:"runmode"`
	}

	// Act
	err := json.Unmarshal(input, &value)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, appkit.RunmodeProd, value.Runmode)
}
