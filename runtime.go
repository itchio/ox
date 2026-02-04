package ox

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Platform string

// these coincide with the namings used in the itch.io backend
const (
	PlatformOSX     Platform = "osx"
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
	PlatformUnknown Platform = "unknown"
)

// Runtime describes an os-arch combo in a convenient way
type Runtime struct {
	Platform Platform `json:"platform"`
	Is64     bool     `json:"is64"`
	Architecture string `json:"arch,omitempty"` // "amd64", "arm64", "386", "arm"
}

type Runtimes []Runtime

func (rs Runtimes) HasPlatform(platform Platform) bool {
	for _, r := range rs {
		if r.Platform == platform {
			return true
		}
	}
	return false
}

func (r Runtime) String() string {
	var arch string
	if r.Is64 {
		arch = "64-bit"
	} else {
		arch = "32-bit"
	}
	var platform = "Unknown"
	switch r.Platform {
	case PlatformLinux:
		platform = "Linux"
	case PlatformOSX:
		platform = "macOS"
	case PlatformWindows:
		platform = "Windows"
	}
	return fmt.Sprintf("%s %s", arch, platform)
}

// OS returns the operating system in GOOS format
func (r Runtime) OS() string {
	switch r.Platform {
	case PlatformLinux:
		return "linux"
	case PlatformOSX:
		return "darwin"
	case PlatformWindows:
		return "windows"
	default:
		return "unknown"
	}
}

// Arch returns the architecture in GOARCH format
func (r Runtime) Arch() string {
	if r.Architecture != "" {
		return r.Architecture
	}
	// Fallback for backwards compatibility with old cached values
	if r.Is64 {
		return "amd64"
	}
	return "386"
}

func (r Runtime) Equals(other Runtime) bool {
	return r.Is64 == other.Is64 && r.Platform == other.Platform
}

var cachedRuntime *Runtime

func CurrentRuntime() Runtime {
	if cachedRuntime == nil {
		arch := detectArch()
		is64 := arch == "amd64" || arch == "arm64"

		var platform Platform
		switch runtime.GOOS {
		case "linux":
			platform = PlatformLinux
		case "darwin":
			platform = PlatformOSX
		case "windows":
			platform = PlatformWindows
		default:
			platform = PlatformUnknown
		}

		cachedRuntime = &Runtime{
			Is64:         is64,
			Platform:     platform,
			Architecture: arch,
		}
	}
	return *cachedRuntime
}

// archMapping maps OS-reported architecture names to GOARCH format
var archMapping = map[string]string{
	// x86_64
	"x86_64": "amd64",
	"AMD64":  "amd64",
	"IA64":   "amd64",
	// arm64
	"arm64":   "arm64",
	"aarch64": "arm64",
	"ARM64":   "arm64",
	// 32-bit x86
	"i386": "386",
	"i686": "386",
	"x86":  "386",
	// 32-bit arm
	"armv7l": "arm",
	"armv6l": "arm",
}

// MapArchitecture converts an OS-reported architecture name to GOARCH format.
// Returns empty string if the architecture is not recognized.
func MapArchitecture(osArch string) string {
	return archMapping[osArch]
}

var hasDeterminedArch = false
var cachedArch string

func detectArch() string {
	switch runtime.GOOS {
	case "darwin":
		return detectDarwinArch()
	case "linux":
		if !hasDeterminedArch {
			cachedArch = detectLinuxArch()
			hasDeterminedArch = true
		}
		return cachedArch
	case "windows":
		return detectWindowsArch()
	}

	// unsupported platform - fall back to compile-time arch
	return runtime.GOARCH
}

func detectDarwinArch() string {
	// Check for Apple Silicon - sysctl returns "1" on ARM Macs even under Rosetta
	output, err := exec.Command("sysctl", "-n", "hw.optional.arm64").Output()
	if err == nil && strings.TrimSpace(string(output)) == "1" {
		return "arm64"
	}
	// Fall back to uname for older Intel Macs
	output, err = exec.Command("uname", "-m").Output()
	if err == nil {
		machine := strings.TrimSpace(string(output))
		if arch, ok := archMapping[machine]; ok {
			return arch
		}
	}
	return runtime.GOARCH
}

func detectLinuxArch() string {
	output, err := exec.Command("uname", "-m").Output()
	if err == nil {
		machine := strings.TrimSpace(string(output))
		if arch, ok := archMapping[machine]; ok {
			return arch
		}
	}

	output, err = exec.Command("arch").Output()
	if err == nil {
		machine := strings.TrimSpace(string(output))
		if arch, ok := archMapping[machine]; ok {
			return arch
		}
	}

	// Fall back to compile-time arch
	return runtime.GOARCH
}

func detectWindowsArch() string {
	// If we're running as a 64-bit executable, check what kind
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
		return runtime.GOARCH
	}

	// 32-bit binary running on potentially 64-bit OS - check env vars
	// PROCESSOR_ARCHITEW6432 is set when a 32-bit process runs on 64-bit Windows
	if archEnv := os.Getenv("PROCESSOR_ARCHITEW6432"); archEnv != "" {
		if arch, ok := archMapping[archEnv]; ok {
			return arch
		}
	}

	if archEnv := os.Getenv("PROCESSOR_ARCHITECTURE"); archEnv != "" {
		if arch, ok := archMapping[archEnv]; ok {
			return arch
		}
	}

	// Fall back to compile-time arch
	return runtime.GOARCH
}
