package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Config struct {
	TestAutoStart bool
}

func TestConfig(t *testing.T, snapName string, conf Config) {
	t.Run("config", func(t *testing.T) {
		TestAutoStart(t, snapName, conf.TestAutoStart)
	})
}

func TestAutoStart(t *testing.T, snapName string, testAutoStart bool) {
	if testAutoStart {
		t.Run("autostart", func(t *testing.T) {
			TestAutostartGlobal(t, snapName)
		})
	}
}

func TestAutostartGlobal(t *testing.T, snapName string) {
	t.Run("set and unset global autostart", func(t *testing.T) {
		t.Cleanup(func() {
			SnapUnset(t, snapName, "autostart")
			SnapStop(t, snapName)
		})

		SnapStop(t, snapName)
		require.False(t, SnapServicesEnabled(t, snapName))
		require.False(t, SnapServicesActive(t, snapName))

		SnapSet(t, snapName, "autostart", "true")
		require.True(t, SnapServicesEnabled(t, snapName))
		require.True(t, SnapServicesActive(t, snapName))

		SnapUnset(t, snapName, "autostart")
		require.True(t, SnapServicesEnabled(t, snapName))
		require.True(t, SnapServicesActive(t, snapName))

		SnapSet(t, snapName, "autostart", "false")
		require.False(t, SnapServicesEnabled(t, snapName))
		require.False(t, SnapServicesActive(t, snapName))
	})
}
