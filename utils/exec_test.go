package utils

import (
	"context"
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

	t.Run("context with cancel", func(t *testing.T) {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		stdout, stderr, err := exec(t, ctx, "sleep 2", true)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if stdout != "" {
			t.Fatalf("Unexpected stdout: %s", stdout)
		}

		if stderr != "" {
			t.Fatalf("Unexpected stderr: %s", stderr)
		}
	})

	t.Run("context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), <-time.After(2*time.Second))
		defer cancel()

		stdout, stderr, err := exec(nil, ctx, `sleep infinity`, true)
		assert.Error(t, err)
		t.Logf("stdout: %s", stdout)
		t.Logf("stdeer: %s", stderr)
	})
}
