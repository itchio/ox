package ox_test

import (
	"runtime"
	"testing"

	"github.com/itchio/ox"
	"github.com/stretchr/testify/assert"
)

// TestCurrentRuntime_Platform verifies that Platform matches runtime.GOOS
func TestCurrentRuntime_Platform(t *testing.T) {
	r := ox.CurrentRuntime()

	var expectedPlatform ox.Platform
	switch runtime.GOOS {
	case "linux":
		expectedPlatform = ox.PlatformLinux
	case "darwin":
		expectedPlatform = ox.PlatformOSX
	case "windows":
		expectedPlatform = ox.PlatformWindows
	default:
		expectedPlatform = ox.PlatformUnknown
	}

	assert.Equal(t, expectedPlatform, r.Platform, "Platform should match runtime.GOOS (%s)", runtime.GOOS)
}

// TestCurrentRuntime_Architecture verifies that Arch() matches runtime.GOARCH
// On native runners (no emulation), detected arch should equal compile-time arch
func TestCurrentRuntime_Architecture(t *testing.T) {
	r := ox.CurrentRuntime()

	// On native hardware, the detected architecture should match GOARCH
	assert.Equal(t, runtime.GOARCH, r.Arch(), "Architecture should match runtime.GOARCH on native hardware")
}

// TestCurrentRuntime_Is64 verifies Is64 is true for 64-bit architectures
func TestCurrentRuntime_Is64(t *testing.T) {
	r := ox.CurrentRuntime()

	switch runtime.GOARCH {
	case "amd64", "arm64":
		assert.True(t, r.Is64, "Is64 should be true for %s", runtime.GOARCH)
	case "386", "arm":
		assert.False(t, r.Is64, "Is64 should be false for %s", runtime.GOARCH)
	default:
		t.Logf("Unknown GOARCH: %s, skipping Is64 assertion", runtime.GOARCH)
	}
}

// TestArchMapping tests the archMapping table by verifying known mappings
func TestArchMapping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// x86_64 variants
		{"x86_64", "amd64"},
		{"AMD64", "amd64"},
		// arm64 variants
		{"arm64", "arm64"},
		{"aarch64", "arm64"},
		{"ARM64", "arm64"},
		// 32-bit x86 variants
		{"i386", "386"},
		{"i686", "386"},
		{"x86", "386"},
		// 32-bit ARM variants
		{"armv7l", "arm"},
		{"armv6l", "arm"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ox.MapArchitecture(tc.input)
			assert.Equal(t, tc.expected, result, "MapArchitecture(%q) should return %q", tc.input, tc.expected)
		})
	}
}

// TestArchMapping_Unknown tests that unknown architectures return empty string
func TestArchMapping_Unknown(t *testing.T) {
	result := ox.MapArchitecture("unknown_arch")
	assert.Equal(t, "", result, "Unknown architecture should return empty string")
}

// TestRuntime_OS verifies OS() returns GOOS-compatible strings
func TestRuntime_OS(t *testing.T) {
	tests := []struct {
		platform ox.Platform
		expected string
	}{
		{ox.PlatformLinux, "linux"},
		{ox.PlatformOSX, "darwin"},
		{ox.PlatformWindows, "windows"},
		{ox.PlatformUnknown, "unknown"},
	}

	for _, tc := range tests {
		t.Run(string(tc.platform), func(t *testing.T) {
			r := ox.Runtime{Platform: tc.platform}
			assert.Equal(t, tc.expected, r.OS(), "OS() for %s should return %q", tc.platform, tc.expected)
		})
	}
}

// TestRuntime_String verifies String() returns human-readable format
func TestRuntime_String(t *testing.T) {
	tests := []struct {
		runtime  ox.Runtime
		expected string
	}{
		{ox.Runtime{Platform: ox.PlatformOSX, Is64: true}, "64-bit macOS"},
		{ox.Runtime{Platform: ox.PlatformOSX, Is64: false}, "32-bit macOS"},
		{ox.Runtime{Platform: ox.PlatformWindows, Is64: true}, "64-bit Windows"},
		{ox.Runtime{Platform: ox.PlatformWindows, Is64: false}, "32-bit Windows"},
		{ox.Runtime{Platform: ox.PlatformLinux, Is64: true}, "64-bit Linux"},
		{ox.Runtime{Platform: ox.PlatformLinux, Is64: false}, "32-bit Linux"},
		{ox.Runtime{Platform: ox.PlatformUnknown, Is64: true}, "64-bit Unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.runtime.String())
		})
	}
}

// TestRuntime_Equals tests the Equals() method
func TestRuntime_Equals(t *testing.T) {
	r1 := ox.Runtime{Platform: ox.PlatformLinux, Is64: true, Architecture: "amd64"}
	r2 := ox.Runtime{Platform: ox.PlatformLinux, Is64: true, Architecture: "amd64"}
	r3 := ox.Runtime{Platform: ox.PlatformLinux, Is64: true, Architecture: "arm64"}
	r4 := ox.Runtime{Platform: ox.PlatformWindows, Is64: true, Architecture: "amd64"}
	r5 := ox.Runtime{Platform: ox.PlatformLinux, Is64: false, Architecture: "386"}

	// Same runtime should be equal
	assert.True(t, r1.Equals(r2), "Identical runtimes should be equal")

	// Same Platform and Is64, different Architecture - still equal (Equals only checks Platform and Is64)
	assert.True(t, r1.Equals(r3), "Runtimes with same Platform/Is64 but different Architecture should be equal")

	// Different Platform
	assert.False(t, r1.Equals(r4), "Runtimes with different Platform should not be equal")

	// Different Is64
	assert.False(t, r1.Equals(r5), "Runtimes with different Is64 should not be equal")
}

// TestRuntime_Arch_Fallback tests that Arch() falls back correctly when Architecture field is empty
func TestRuntime_Arch_Fallback(t *testing.T) {
	// Test fallback for 64-bit (should return "amd64")
	r64 := ox.Runtime{Platform: ox.PlatformLinux, Is64: true, Architecture: ""}
	assert.Equal(t, "amd64", r64.Arch(), "Empty Architecture with Is64=true should fall back to amd64")

	// Test fallback for 32-bit (should return "386")
	r32 := ox.Runtime{Platform: ox.PlatformLinux, Is64: false, Architecture: ""}
	assert.Equal(t, "386", r32.Arch(), "Empty Architecture with Is64=false should fall back to 386")

	// Test that explicit Architecture takes precedence
	rExplicit := ox.Runtime{Platform: ox.PlatformLinux, Is64: true, Architecture: "arm64"}
	assert.Equal(t, "arm64", rExplicit.Arch(), "Explicit Architecture should take precedence")
}

// TestRuntimes_HasPlatform tests the HasPlatform method on Runtimes slice
func TestRuntimes_HasPlatform(t *testing.T) {
	runtimes := ox.Runtimes{
		{Platform: ox.PlatformLinux, Is64: true},
		{Platform: ox.PlatformWindows, Is64: true},
	}

	assert.True(t, runtimes.HasPlatform(ox.PlatformLinux), "Should have Linux platform")
	assert.True(t, runtimes.HasPlatform(ox.PlatformWindows), "Should have Windows platform")
	assert.False(t, runtimes.HasPlatform(ox.PlatformOSX), "Should not have OSX platform")
	assert.False(t, runtimes.HasPlatform(ox.PlatformUnknown), "Should not have Unknown platform")
}

// TestCurrentRuntime_Consistency verifies multiple calls return same values
func TestCurrentRuntime_Consistency(t *testing.T) {
	r1 := ox.CurrentRuntime()
	r2 := ox.CurrentRuntime()

	assert.Equal(t, r1.Platform, r2.Platform, "Platform should be consistent across calls")
	assert.Equal(t, r1.Is64, r2.Is64, "Is64 should be consistent across calls")
	assert.Equal(t, r1.Architecture, r2.Architecture, "Architecture should be consistent across calls")
}
