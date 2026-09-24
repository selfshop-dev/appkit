// Package appkit provides common application infrastructure for Go services.
//
// The package is built around a small set of application-level primitives:
// configuration types, run mode parsing, graceful shutdown timing, cleanup
// management, and deferred HTTP handlers.
//
// # Configuration
//
// [Config] contains the main application configuration, including the
// application [Runmode], name, and [ShutdownConfig].
//
// A zero Config is valid. [Config.IsDevmod] reports whether the application is
// running in development mode.
//
//	cfg := appkit.Config{
//		Runmode: appkit.RunmodeDev,
//		Name:    "my-service",
//	}
//
//	if cfg.IsDevmod() {
//		// enable development-specific behavior
//	}
//
// [ShutdownConfig] combines the main shutdown timeout with an optional drain
// delay. [ShutdownConfig.TotalTimeout] returns the total amount of time
// allocated for shutdown.
//
// # Run modes
//
// [Runmode] defines the supported application modes: [RunmodeDev] and
// [RunmodeProd].
//
// [ParseRunmode] parses a textual run mode case-insensitively. An empty value
// is treated as [RunmodeProd].
//
//	runmode, err := appkit.ParseRunmode("dev")
//	if err != nil {
//		// handle invalid configuration
//	}
//
//	fmt.Println(runmode)
//	// Output:
//	// dev
//
// [Runmode.MarshalText] and [Runmode.UnmarshalText] implement text
// serialization, making Runmode suitable for configuration systems and
// encoding packages. Invalid run modes are rejected during serialization.
//
// # Shutdown and drain delay
//
// [DrainDelay] represents the amount of time an application should wait while
// draining in-flight work during shutdown.
//
// The value can be converted to a [time.Duration] with [DrainDelay.Duration]
// and formatted with [DrainDelay.String].
//
// [DrainDelay.Wait] waits for the configured duration or returns early when
// the supplied context is cancelled.
//
//	delay := appkit.DrainDelay(5 * time.Second)
//
//	if err := delay.Wait(ctx); err != nil {
//		// shutdown was cancelled
//	}
//
// [DrainDelay] supports text serialization through [DrainDelay.MarshalText]
// and [DrainDelay.UnmarshalText]. Duration strings use the format accepted by
// [time.ParseDuration]. Empty configuration values are treated as zero.
// Negative drain delays are rejected.
//
// # Cleanup
//
// [CleanupStack] provides LIFO cleanup management for application resources.
// Cleanup functions are registered with [CleanupStack.Add] and executed in
// reverse registration order by [CleanupStack.Finalize].
//
//	var cleanup appkit.CleanupStack
//
//	cleanup.Add("database", func(ctx context.Context) error {
//		return db.Close()
//	})
//
//	cleanup.Add("server", func(ctx context.Context) error {
//		return server.Shutdown(ctx)
//	})
//
//	if err := cleanup.Finalize(ctx); err != nil {
//		// handle cleanup errors
//	}
//
// Cleanup functions are executed even when an earlier cleanup returns an
// error. Errors from multiple cleanup functions are combined and retain their
// individual causes for use with [errors.Is] and [errors.As].
//
// A nil cleanup function passed to [CleanupStack.Add] is ignored. Cleanup
// names are trimmed; an empty name is replaced with "<unnamed>". Calling
// [CleanupStack.Finalize] clears the registered cleanup functions.
//
// # Deferred HTTP handlers
//
// [NewDeferredHandler] creates an HTTP handler whose implementation can be
// replaced after the handler has been installed.
//
// This is useful when an HTTP endpoint needs to exist before the application
// has completed initialization. When created with a nil base handler, the
// handler responds with HTTP 503 Service Unavailable until a real handler is
// installed.
//
//	handler := appkit.NewDeferredHandler(nil)
//
//	// The endpoint is available but reports that the service is not ready.
//
//	handler.Set(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintln(w, "ready")
//	}))
//
// [NewDeferredHandler] uses an atomic handler reference, so replacing the
// handler with Set is safe while requests are being served concurrently.
//
// The returned handler implements [http.Handler] and can therefore be
// registered directly with an HTTP server or router.
//
// # Serialization
//
// Configuration-related types implement the standard text marshaling
// interfaces where appropriate. This allows [Runmode] and [DrainDelay] to be
// used with configuration libraries and encoding packages that rely on
// [encoding.TextMarshaler] and [encoding.TextUnmarshaler].
//
// Invalid configuration values are rejected instead of being silently
// normalized, while empty values that have an explicit default are handled
// according to the type's documented behavior.
package appkit
