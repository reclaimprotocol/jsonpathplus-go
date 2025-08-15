package jsonpathplus

import (
	"fmt"
	"runtime"
)

const (
	// Version represents the current version of the library.
	Version = "1.1.1"

	// APIVersion represents the API version for compatibility tracking.
	APIVersion = "v1"

	// MinSupportedGoVersion is the minimum Go version required.
	MinSupportedGoVersion = "1.19"
)

// BuildInfo contains build and version information.
type BuildInfo struct {
	Version    string `json:"version"`
	APIVersion string `json:"api_version"`
	GoVersion  string `json:"go_version"`
	GitCommit  string `json:"git_commit,omitempty"`
	BuildTime  string `json:"build_time,omitempty"`
	Platform   string `json:"platform"`
	Compiler   string `json:"compiler"`
}

// GetBuildInfo returns build information.
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:    Version,
		APIVersion: APIVersion,
		GoVersion:  runtime.Version(),
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Compiler:   runtime.Compiler,
	}
}

// IsCompatible checks if the given API version is compatible.
func IsCompatible(apiVersion string) bool {
	return apiVersion == APIVersion || apiVersion == "v1.0" || apiVersion == "1.0"
}

// CheckGoVersion verifies the Go runtime version meets minimum requirements.
func CheckGoVersion() error {
	version := runtime.Version()

	// Simple version check - in production you'd want more robust parsing
	if len(version) < 4 || version[:4] < "go1.19" {
		return fmt.Errorf("go version %s is not supported, minimum required: %s",
			version, MinSupportedGoVersion)
	}

	return nil
}
