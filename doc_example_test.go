package appkit_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/selfshop-dev/appkit"
)

func ExampleCleanupStack() {
	var stack appkit.CleanupStack

	stack.Add("database", func(ctx context.Context) error {
		fmt.Println("close database")
		return nil
	})

	stack.Add("server", func(ctx context.Context) error {
		fmt.Println("stop server")
		return nil
	})

	err := stack.Finalize(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// stop server
	// close database
}

func ExampleCleanupStack_error() {
	var stack appkit.CleanupStack

	stack.Add("database", func(ctx context.Context) error {
		return errors.New("connection close failed")
	})

	err := stack.Finalize(context.Background())
	fmt.Println(err)

	// Output:
	// database: connection close failed
}

func ExampleConfig_IsDevmod() {
	config := appkit.Config{
		Runmode: appkit.RunmodeDev,
		Name:    "my-app",
	}

	fmt.Println(config.IsDevmod())

	// Output:
	// true
}

func ExampleShutdownConfig_TotalTimeout() {
	config := appkit.ShutdownConfig{
		Timeout:    5 * time.Second,
		DrainDelay: appkit.DrainDelay(2 * time.Second),
	}

	fmt.Println(config.TotalTimeout())

	// Output:
	// 7s
}

func ExampleParseRunmode() {
	runmode, err := appkit.ParseRunmode("DEV")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(runmode)

	// Output:
	// dev
}

func ExampleRunmode_Set() {
	var runmode appkit.Runmode

	if err := runmode.Set("prod"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(runmode)

	// Output:
	// prod
}

func ExampleRunmode_MarshalText() {
	runmode := appkit.RunmodeDev

	value, err := runmode.MarshalText()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(value))

	// Output:
	// dev
}

func ExampleDrainDelay() {
	delay := appkit.DrainDelay(5 * time.Second)

	fmt.Println(delay.Duration())
	fmt.Println(delay.String())

	// Output:
	// 5s
	// 5s
}

func ExampleDrainDelay_UnmarshalText() {
	var delay appkit.DrainDelay

	if err := delay.UnmarshalText([]byte("2s")); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(delay)

	// Output:
	// 2s
}

func ExampleDrainDelay_Wait() {
	delay := appkit.DrainDelay(0)

	if err := delay.Wait(context.Background()); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("delay completed")

	// Output:
	// delay completed
}

func ExampleNewDeferredHandler() {
	// Arrange
	handler := appkit.NewDeferredHandler(nil)

	// Act
	handler.Set(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ready") //nolint:errcheck // example: ignore write error
	}))

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	// Assert / Output
	fmt.Println(recorder.Code)
	fmt.Print(recorder.Body.String())

	// Output:
	// 200
	// ready
}
