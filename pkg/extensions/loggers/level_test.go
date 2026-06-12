package loggers

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout runs f and returns everything it wrote to os.Stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	f()

	require.NoError(t, w.Close())
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	return buf.String()
}

func TestECSLevelFiltersLowerSeverity(t *testing.T) {
	ctx := context.Background()

	// With an error-level logger, info and warning must be discarded but error
	// must still be printed.
	logger := NewECS(WithLevel(LevelError))

	assert.Empty(t, captureStdout(t, func() {
		logger.Info(ctx, "an info message")
	}), "info should be filtered out at error level")

	assert.Empty(t, captureStdout(t, func() {
		logger.Warning(ctx, "a warning message")
	}), "warning should be filtered out at error level")

	out := captureStdout(t, func() {
		logger.Error(ctx, "an error message")
	})
	assert.Contains(t, out, "an error message", "error should always be printed")
}

func TestECSWarningLevelKeepsWarningAndError(t *testing.T) {
	ctx := context.Background()
	logger := NewECS(WithLevel(LevelWarning))

	assert.Empty(t, captureStdout(t, func() {
		logger.Info(ctx, "an info message")
	}), "info should be filtered out at warning level")

	assert.Contains(t, captureStdout(t, func() {
		logger.Warning(ctx, "a warning message")
	}), "a warning message")

	assert.Contains(t, captureStdout(t, func() {
		logger.Error(ctx, "an error message")
	}), "an error message")
}

func TestECSDefaultLevelLogsEverything(t *testing.T) {
	ctx := context.Background()
	logger := NewECS()

	assert.Contains(t, captureStdout(t, func() {
		logger.Info(ctx, "an info message")
	}), "an info message", "default logger should print info")
}

func TestTextLevelFiltersLowerSeverity(t *testing.T) {
	ctx := context.Background()
	logger := NewText(WithLevel(LevelError))

	assert.Empty(t, captureStdout(t, func() {
		logger.Info(ctx, "an info message")
	}), "info should be filtered out at error level")

	assert.NotEmpty(t, captureStdout(t, func() {
		logger.Error(ctx, "an error message")
	}), "error should always be printed")
}
