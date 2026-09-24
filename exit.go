package appkit

// ExitOnError runs run and terminates the application with exit code 1 when
// run returns an error.
//
// The exit function is injected by the caller, which allows the behavior to
// be tested without calling [os.Exit] directly.
//
// If run returns nil, exit is not called.
func ExitOnError(
	run func() error, exit func(code int),
) {
	if err := run(); err != nil {
		exit(1)
	}
}
