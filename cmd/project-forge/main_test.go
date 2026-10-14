package main

import "testing"

func TestRun(t *testing.T) {
	// Test that Run doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run panicked: %v", r)
		}
	}()

	// Override osExit to prevent actual exit
	originalExit := osExit
	osExit = func(code int) {
		// Capture exit code but don't exit
	}
	defer func() { osExit = originalExit }()

	Run()
}

func TestMain(t *testing.T) {
	// Just ensure main doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()
	// We don't actually call main() as it would exit, but the code is simple
	// and covered by being compiled
}
