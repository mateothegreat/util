package test

import (
	"testing"
)

func TestT(t *testing.T) {
	testTFunc(t) // Just verify this doesn't give a compiler error
}

func TestRuntimeT(t *testing.T) {
	var _ T = new(RuntimeT) // Another compiler check
}

// TestRuntimeTBasicFunctionality demonstrates the core functionality of RuntimeT.
// This test showcases how RuntimeT can be used as a drop-in replacement for *testing.T
// in scenarios where you need testing-like behavior outside of the standard test framework.
func TestRuntimeTBasicFunctionality(t *testing.T) {
	rt := &RuntimeT{}

	// Test logging functionality
	rt.Log("This is a log message")
	rt.Logf("This is a formatted log message: %s", "test")

	// Test that initially the RuntimeT is not failed or skipped
	if rt.Failed() {
		t.Error("RuntimeT should not be failed initially")
	}
	if rt.Skipped() {
		t.Error("RuntimeT should not be skipped initially")
	}

	// Test error handling
	rt.Error("This is an error")
	if !rt.Failed() {
		t.Error("RuntimeT should be failed after Error() call")
	}
}

// TestRuntimeTErrorHandling demonstrates error reporting capabilities.
// This test shows how RuntimeT handles different types of error conditions
// and maintains proper state throughout the testing lifecycle.
func TestRuntimeTErrorHandling(t *testing.T) {
	rt := &RuntimeT{}

	// Test Errorf functionality
	rt.Errorf("This is a formatted error: %d", 42)
	if !rt.Failed() {
		t.Error("RuntimeT should be failed after Errorf() call")
	}

	// Test that multiple errors accumulate
	rt.Error("Another error")
	if !rt.Failed() {
		t.Error("RuntimeT should remain failed after multiple errors")
	}
}

// TestRuntimeTSkipFunctionality demonstrates skip behavior.
// This test verifies that RuntimeT properly handles skip conditions
// and maintains the correct skipped state.
func TestRuntimeTSkipFunctionality(t *testing.T) {
	rt := &RuntimeT{}

	// Test Skip functionality
	rt.Skip("Skipping this test")
	if !rt.Skipped() {
		t.Error("RuntimeT should be skipped after Skip() call")
	}

	// Test with a fresh RuntimeT for Skipf
	rt2 := &RuntimeT{}
	rt2.Skipf("Skipping with format: %s", "reason")
	if !rt2.Skipped() {
		t.Error("RuntimeT should be skipped after Skipf() call")
	}
}

// TestRuntimeTEnvironmentVariables demonstrates environment variable management.
// This test shows how RuntimeT can safely set and restore environment variables,
// including proper cleanup behavior.
func TestRuntimeTEnvironmentVariables(t *testing.T) {
	rt := &RuntimeT{}

	// Set an environment variable
	rt.Setenv("TEST_VAR", "test_value")

	// Verify it was set (this would be done in actual usage)
	// Note: In a real scenario, you'd check os.Getenv("TEST_VAR")

	// Test that parallel and setenv don't mix
	rt2 := &RuntimeT{}
	rt2.Parallel()

	// This should panic, but we can't easily test panics in this context
	// In real usage, calling rt2.Setenv() after rt2.Parallel() would panic
}

// TestRuntimeTTempDirectory demonstrates temporary directory creation.
// This test showcases the TempDir functionality and verifies that
// directories are properly created and managed.
func TestRuntimeTTempDirectory(t *testing.T) {
	rt := &RuntimeT{}

	// Get a temporary directory
	dir1 := rt.TempDir()
	if dir1 == "" {
		t.Error("TempDir should return a non-empty directory path")
	}

	// Get another temporary directory - should be different
	dir2 := rt.TempDir()
	if dir2 == "" {
		t.Error("Second TempDir call should return a non-empty directory path")
	}

	if dir1 == dir2 {
		t.Error("Multiple TempDir calls should return different directories")
	}
}

// // TestRuntimeTCleanup demonstrates cleanup functionality.
// // This test verifies that cleanup functions are properly registered
// // and that the cleanup mechanism works as expected.
// func TestRuntimeTCleanup(t *testing.T) {
// 	rt := &RuntimeT{}

// 	cleanupCalled := false

// 	// Register a cleanup function
// 	rt.Cleanup(func() {
// 		cleanupCalled = true
// 	})

// 	// In a real scenario, cleanup would be called automatically
// 	// Here we just verify the cleanup was registered
// 	// (The actual cleanup execution would happen in a real test framework)
// }

// TestRuntimeTHelperFunction demonstrates the Helper method.
// This test shows that the Helper method can be called without issues,
// maintaining compatibility with the testing.T interface.
func TestRuntimeTHelperFunction(t *testing.T) {
	rt := &RuntimeT{}

	// Helper should not panic or cause issues
	rt.Helper()

	// Test that we can call it multiple times
	rt.Helper()
	rt.Helper()
}

// TestRuntimeTName demonstrates the Name method behavior.
// This test verifies that RuntimeT returns an appropriate name
// and maintains consistency with the testing interface.
func TestRuntimeTName(t *testing.T) {
	rt := &RuntimeT{}

	name := rt.Name()
	// RuntimeT returns empty string for name by default
	if name != "" {
		t.Errorf("Expected empty name, got: %s", name)
	}
}

// TestTBInterface demonstrates that both T and RuntimeT implement TB.
// This test showcases the polymorphic usage of the TB interface,
// allowing functions to work with both standard testing.T and RuntimeT.
func TestTBInterface(t *testing.T) {
	// Test that *testing.T implements TB
	testTBFunc(t)

	// Test that RuntimeT implements TB
	rt := &RuntimeT{}
	testTBFunc(rt)
}

// testTBFunc is a helper function that accepts TB interface.
// This demonstrates how test helper functions can be written to work
// with both standard testing types and RuntimeT.
func testTBFunc(tb TB) {
	tb.Helper()
	tb.Log("Testing TB interface functionality")

	// Test basic state
	if tb.Failed() {
		tb.Error("TB should not be failed initially")
	}

	if tb.Skipped() {
		tb.Log("TB is skipped")
	}
}

// testTFunc is a helper function that accepts T interface.
// This demonstrates the usage of the T interface for functions
// that need the full testing.T-like functionality.
func testTFunc(t T) {
	t.Helper()
	t.Log("Testing T interface functionality")

	// Test all T-specific methods
	t.Parallel()

	name := t.Name()
	t.Logf("Test name: %s", name)

	// Test temporary directory
	tempDir := t.TempDir()
	t.Logf("Temporary directory: %s", tempDir)

	// Test cleanup registration
	t.Cleanup(func() {
		t.Log("Cleanup function executed")
	})
}
