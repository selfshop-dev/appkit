package appkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeferredHandler(t *testing.T) {
	// Arrange
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// Act
	handler := NewDeferredHandler(base)

	// Assert
	require.NotNil(t, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestNewDeferredHandler_NilBase(t *testing.T) {
	// Arrange
	handler := NewDeferredHandler(nil)

	// Act
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	assert.Equal(
		t,
		"Service Unavailable\n",
		recorder.Body.String(),
	)
}

func TestDeferredHandler_Set(t *testing.T) {
	// Arrange
	first := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	second := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	handler := NewDeferredHandler(first)

	// Act
	handler.Set(second)

	// Assert
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusAccepted, recorder.Code)
}

func TestDeferredHandler_ServeHTTP(t *testing.T) {
	// Arrange
	var gotMethod string
	var gotPath string

	handler := NewDeferredHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotPath = r.URL.Path

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/test", gotPath)
}

func TestDeferredHandler_Set_ReplacesHandler(t *testing.T) {
	// Arrange
	handler := NewDeferredHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}),
	)

	handler.Set(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	// Act
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusAccepted, recorder.Code)
}

func TestDeferredHandler_Set_Nil(t *testing.T) {
	// Arrange
	handler := NewDeferredHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}),
	)

	// Act
	handler.Set(nil)

	// Assert
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	assert.Panics(t, func() {
		handler.ServeHTTP(recorder, request)
	})
}
