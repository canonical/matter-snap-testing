package utils

import (
	"context"
	goexec "os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {

	t.Run("one command", func(t *testing.T) {
		stdout, stderr, err := exec(t, nil, `echo "hi"`, true)
		assert.NoError(t, err)
		assert.Empty(t, stderr)
		assert.Equal(t, "hi\n", stdout)
	})

	t.Run("exit after slow command", func(t *testing.T) {
		start := time.Now()
		stdout, _, err := exec(t, nil, `echo "hi" && sleep 0.1 && echo "hi2"`, true)
		// must return after 100ms±50ms
		require.WithinDuration(t,
			start.Add(100*time.Millisecond),
			time.Now(),
			50*time.Millisecond)
		assert.NoError(t, err)
		assert.Equal(t, "hi\nhi2\n", stdout)
	})

	t.Run("bad command", func(t *testing.T) {
		stdout, stderr, err := exec(nil, nil, `bad_command`, true)
		assert.Error(t, err)
		assert.Empty(t, stdout)
		assert.Contains(t, stderr, "not found")
	})

	t.Run("print to stderr", func(t *testing.T) {
		stdout, stderr, err := exec(t, nil, `echo "failing" >&2`, true)
		assert.NoError(t, err)
		assert.Empty(t, stdout)
		assert.Equal(t, "failing\n", stderr)
	})

	t.Run("stderr then stdout", func(t *testing.T) {
		stdout, stderr, err := exec(t, nil, `echo "failing" >&2; echo "succeeding"`, true)
		assert.NoError(t, err)
		assert.Equal(t, "failing\n", stderr)
		assert.Equal(t, "succeeding\n", stdout)
	})

	t.Run("context timed out", func(t *testing.T) {
		testTimeout := func(t *testing.T, command string) {
			timeout := 1 * time.Second

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			t.Cleanup(cancel)

			start := time.Now()
			_, _, err := exec(nil, ctx, command, true)

			require.NoError(t, err)
			require.WithinDuration(t, start, time.Now(), timeout+500*time.Millisecond)
		}

		t.Run("user+bash", func(t *testing.T) {
			testTimeout(t, "sleep 10")
		})

		t.Run("root+bash", func(t *testing.T) {
			testTimeout(t, "sudo sleep 10")
		})

		t.Run("root", func(t *testing.T) {
			timeout := 1 * time.Second

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			t.Cleanup(cancel)

			start := time.Now()
			out, err := goexec.CommandContext(ctx, "sudo", "sleep", "10").CombinedOutput()

			t.Logf("output: %s", out)
			require.Error(t, err)
			require.WithinDuration(t, start, time.Now(), timeout+500*time.Millisecond)
		})

	})

	t.Run("context not timed out", func(t *testing.T) {
		timeout := 3 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		t.Cleanup(cancel)

		_, _, err := exec(nil, ctx, `sleep 1`, true)

		require.NoError(t, err)
	})
}
