package loggers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
)

// ECS is a logger that will print logs in Elastic Common Schema ECS format.
type ECS struct{}

// NewECS creates a new ECS logger.
func NewECS() ECS {
	return ECS{}
}

func (ecs ECS) setInfoFromContext(ctx context.Context, msg string, info ...extensions.LogInfo) []extensions.LogInfo {
	// Add additional keys from context
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsProvider, func(value any) {
		info = append(info, extensions.LogInfo{Key: "asyncapi.provider", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsChannel, func(value any) {
		info = append(info, extensions.LogInfo{Key: "asyncapi.channel", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsDirection, func(value any) {
		if value == "publication" {
			info = append(info, extensions.LogInfo{Key: "event.action", Value: "published-message"})
		} else if value == "reception" {
			info = append(info, extensions.LogInfo{Key: "event.action", Value: "received-message"})
		}
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsBrokerMessage, func(value any) {
		info = append(info, extensions.LogInfo{Key: "event.original", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsCorrelationID, func(value any) {
		info = append(info, extensions.LogInfo{Key: "trace.id", Value: value})
	})

	// Add additional keys
	info = append(info, extensions.LogInfo{
		Key:   "message",
		Value: msg,
	})
	info = append(info, extensions.LogInfo{
		Key:   "@timestamp",
		Value: time.Now().UTC().Format(time.RFC3339Nano),
	})
	info = append(info, extensions.LogInfo{
		Key:   "log.logger",
		Value: "github.com/TheSadlig/asyncapi-codegen/pkg/extensions/loggers/ecs.go",
	})

	// Return info
	return info
}

func (ecs ECS) formatLog(ctx context.Context, msg string, info ...extensions.LogInfo) string {
	// Set additional fields
	info = ecs.setInfoFromContext(ctx, msg, info...)

	// Structure log
	sl := structureLogs(info)

	// Print log
	b, err := json.Marshal(sl)
	if err != nil {
		return "{\"error\":\"error while marshalling log\"}"
	}

	return string(b)
}

func (ecs ECS) logWithLevel(ctx context.Context, level string, msg string, info ...extensions.LogInfo) {
	// Add additional keys
	info = append(info, extensions.LogInfo{Key: "log.level", Value: level})

	// Print log
	fmt.Println(ecs.formatLog(ctx, msg, info...))
}

// Info logs a message at info level with context and additional info.
func (ecs ECS) Info(ctx context.Context, msg string, info ...extensions.LogInfo) {
	ecs.logWithLevel(ctx, "info", msg, info...)
}

// Warning logs a message at warning level with context and additional info.
func (ecs ECS) Warning(ctx context.Context, msg string, info ...extensions.LogInfo) {
	ecs.logWithLevel(ctx, "warning", msg, info...)
}

// Error logs a message at error level with context and additional info.
func (ecs ECS) Error(ctx context.Context, msg string, info ...extensions.LogInfo) {
	ecs.logWithLevel(ctx, "error", msg, info...)
}
