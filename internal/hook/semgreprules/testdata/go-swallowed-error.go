package main

import "log/slog"

func logged() {
	err := something()
	// ruleid:go-swallowed-error
	if err != nil {
		slog.Warn("something failed")
	}
}

func commented() {
	err := something()
	// ok:go-swallowed-error
	if err != nil {
		// best-effort: failure here doesn't affect the caller
		slog.Warn("something failed")
	}
}

func propagated() error {
	err := something()
	// ok:go-swallowed-error
	if err != nil {
		return err
	}
	return nil
}

func something() error { return nil }
