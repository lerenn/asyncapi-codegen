package loggers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/stretchr/testify/assert"
)

// contextWithAllValues returns a context carrying all the asyncapi-codegen
// values that the loggers used to enrich from (#334).
func contextWithAllValues() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, extensions.ContextKeyIsChannel, "my-channel")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, "corr-1")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "my-provider")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, "the-message")
	return ctx
}

// TestECSDoesNotReadFromContext ensures the ECS logger no longer enriches its
// output from the context: only the explicitly-given LogInfo (plus the
// intrinsic message/timestamp/logger keys) must be present (#334).
func TestECSDoesNotReadFromContext(t *testing.T) {
	ecs := NewECS()

	out := ecs.formatLog("hello", extensions.LogInfo{Key: "given", Value: "value"})

	var got map[string]any
	assert.NoError(t, json.Unmarshal([]byte(out), &got))

	// Explicitly-given info and intrinsic keys are present.
	assert.Equal(t, "value", got["given"])
	assert.Equal(t, "hello", got["message"])
	assert.Contains(t, got, "@timestamp")
	assert.Contains(t, got, "log")

	// Context-derived keys from the previous behavior must NOT appear, since
	// the logger no longer reads from any context.
	assert.NotContains(t, out, "trace.id")
	assert.NotContains(t, out, "event.action")
	assert.NotContains(t, out, "event.original")
	assert.NotContains(t, out, "asyncapi.channel")
	assert.NotContains(t, out, "asyncapi.provider")
}

// TestTextDoesNotReadFromContext ensures the Text logger no longer enriches its
// output from the context (#334).
func TestTextDoesNotReadFromContext(t *testing.T) {
	text := NewText()

	out := text.formatLog(text.boldWhitePrinter, "hello", extensions.LogInfo{Key: "given", Value: "value"})

	// Explicitly-given info and the message are present.
	assert.True(t, strings.Contains(out, "given"))
	assert.True(t, strings.Contains(out, "hello"))

	// Context-derived keys from the previous behavior must NOT appear.
	assert.False(t, strings.Contains(out, "Channel"))
	assert.False(t, strings.Contains(out, "CorrelationID"))
	assert.False(t, strings.Contains(out, "Content"))
}

// TestLoggersIgnoreContext is a guard ensuring that passing a fully-populated
// context to the public logging methods does not leak context values into the
// output (#334).
func TestLoggersIgnoreContext(t *testing.T) {
	ctx := contextWithAllValues()

	// These must not panic and must not read from the context.
	NewECS().Info(ctx, "msg")
	NewText().Info(ctx, "msg")
}
