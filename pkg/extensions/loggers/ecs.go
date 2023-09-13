package loggers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// ECS is a logger that will print logs in Elastic Common Schema format
type ECS struct{}

// NewECS creates a new ECS logger
func NewECS() ECS {
	return ECS{}
}

func insertLogIntoStruct(key, value string, m map[string]any) map[string]any {
	// Split key
	l := strings.Split(key, ".")

	// Check if there is no depth, just add it to the map
	if len(l) == 1 {
		m[key] = value
		return m
	}

	// Check if the submap exists, otherwise create it
	var subm map[string]any
	if v, ok := m[l[0]]; !ok {
		subm = make(map[string]any)
	} else {
		subm, ok = v.(map[string]any)
		if !ok {
			// Explicitely drop the old value
			subm = make(map[string]any)
		}
	}

	// Insert the log into the submap
	subm = insertLogIntoStruct(strings.Join(l[1:], "."), value, subm)

	// Insert the submap into the map
	m[l[0]] = subm

	return m
}

func structureLogs(info []extensions.LogInfo) map[string]any {
	structuredLog := make(map[string]any)
	for _, logInfo := range info {
		structuredLog = insertLogIntoStruct(logInfo.Key, fmt.Sprintf("%+v", logInfo.Value), structuredLog)
	}
	return structuredLog
}

func (logger ECS) formatLog(ctx context.Context, msg string, info ...extensions.LogInfo) string {
	// Add additional keys from context
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsProvider, func(value any) {
		info = append(info, extensions.LogInfo{Key: "asyncapi.provider", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsChannel, func(value any) {
		info = append(info, extensions.LogInfo{Key: "asyncapi.channel", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsMessageDirection, func(value any) {
		if value == "publication" {
			info = append(info, extensions.LogInfo{Key: "event.action", Value: "published-message"})
		} else if value == "reception" {
			info = append(info, extensions.LogInfo{Key: "event.action", Value: "received-message"})
		}
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsMessage, func(value any) {
		info = append(info, extensions.LogInfo{Key: "event.original", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsCorrelationID, func(value any) {
		info = append(info, extensions.LogInfo{Key: "trace.id", Value: value})
	})

	// Add additional keys
	info = append(info, extensions.LogInfo{Key: "message", Value: msg})
	info = append(info, extensions.LogInfo{Key: "@timestamp", Value: time.Now().UTC().Format(time.RFC3339Nano)})
	info = append(info, extensions.LogInfo{Key: "log.logger", Value: "github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers/ecs.go"})

	// Structure log
	sl := structureLogs(info)

	// Print log
	b, err := json.Marshal(sl)
	if err != nil {
		return "{\"error\":\"error while marshalling log\"}"
	}

	return string(b)
}

func (logger ECS) logWithLevel(ctx context.Context, level string, msg string, info ...extensions.LogInfo) {
	// Add additional keys
	info = append(info, extensions.LogInfo{Key: "log.level", Value: "info"})

	// Print log
	fmt.Println(logger.formatLog(ctx, msg, info...))
}

// Info logs a message at info level with context and additional info
func (logger ECS) Info(ctx context.Context, msg string, info ...extensions.LogInfo) {
	logger.logWithLevel(ctx, "info", msg, info...)
}

// Warning logs a message at warning level with context and additional info
func (logger ECS) Warning(ctx context.Context, msg string, info ...extensions.LogInfo) {
	logger.logWithLevel(ctx, "warning", msg, info...)
}

// Error logs a message at error level with context and additional info
func (logger ECS) Error(ctx context.Context, msg string, info ...extensions.LogInfo) {
	logger.logWithLevel(ctx, "error", msg, info...)
}
