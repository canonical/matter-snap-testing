package env

import (
	"os"
	"strconv"
)

// Environment variables, used to override defaults
const (
	// Channel/Revision of the service snap (has default)
	EnvSnapChannel = "SNAP_CHANNEL"

	// Path to snap instead, used for testing a local snap instead of
	// downloading from the store
	EnvSnapPath = "SNAP_PATH"

	// Toggle the teardown operations during tests (has default)
	EnvTeardown = "TEARDOWN"
)

var (
	// Defaults
	snapChannel = "latest/edge"
	snapPath    = ""
	teardown    = true
)

// SnapChannel returns the set snap channel
func SnapChannel() string {
	return snapChannel
}

// SnapPath returns the set path to a local snap
func SnapPath() string {
	return snapPath
}

// SkipTeardownRemoval return
func Teardown() (skip bool) {
	return teardown
}

func init() {
	loadEnvVars()
}

// Read environment variables and perform type conversion/casting
func loadEnvVars() {

	if v := os.Getenv(EnvSnapChannel); v != "" {
		snapChannel = v
	}

	if v := os.Getenv(EnvSnapPath); v != "" {
		snapPath = v
	}

	if v := os.Getenv(EnvTeardown); v != "" {
		var err error
		teardown, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	}
}
