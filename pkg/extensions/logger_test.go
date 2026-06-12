package extensions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogInfosFromContext ensures the helper extracts the asyncapi-codegen
// context values into LogInfo entries, and only the ones that are set (#133).
func TestLogInfosFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyIsChannel, "my-channel")
	ctx = context.WithValue(ctx, ContextKeyIsCorrelationID, "corr-1")

	infos := LogInfosFromContext(ctx)

	got := make(map[string]any, len(infos))
	for _, info := range infos {
		got[info.Key] = info.Value
	}

	assert.Equal(t, "my-channel", got["channel"])
	assert.Equal(t, "corr-1", got["correlationID"])

	// Values that were not set must not be present.
	_, hasProvider := got["provider"]
	assert.False(t, hasProvider)
	assert.Len(t, infos, 2)
}

// TestLogInfosFromContextEmpty ensures an empty context yields no LogInfo.
func TestLogInfosFromContextEmpty(t *testing.T) {
	assert.Empty(t, LogInfosFromContext(context.Background()))
}
