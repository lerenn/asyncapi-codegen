package log

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

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

func structureLogs(info []AdditionalInfo) map[string]any {
	structuredLog := make(map[string]any)
	for _, logInfo := range info {
		structuredLog = insertLogIntoStruct(logInfo.Key, fmt.Sprintf("%+v", logInfo.Value), structuredLog)
	}
	return structuredLog
}

func (logger ECS) formatLog(ctx Context, msg string, info ...AdditionalInfo) string {
	// Add additional keys from context
	if ctx.Module != "" {
		info = append(info, AdditionalInfo{"event.module", ctx.Module})
	}
	if ctx.Provider != "" {
		info = append(info, AdditionalInfo{"event.provider", ctx.Provider})
	}
	if ctx.Action != "" {
		info = append(info, AdditionalInfo{"event.action", ctx.Action})
	}
	if ctx.Operation != "" {
		info = append(info, AdditionalInfo{"event.reason", ctx.Operation})
	}
	if ctx.Message != nil {
		info = append(info, AdditionalInfo{"event.original", ctx.Message})
	}
	if ctx.CorrelationID != "" {
		info = append(info, AdditionalInfo{"trace.id", ctx.CorrelationID})
	}

	// Add additional keys
	info = append(info, AdditionalInfo{"message", msg})
	info = append(info, AdditionalInfo{"@timestamp", time.Now().UTC().Format(time.RFC3339Nano)})
	info = append(info, AdditionalInfo{"log.logger", "github.com/lerenn/asyncapi-codegen/pkg/loggers/ecs.go"})

	// Structure log
	sl := structureLogs(info)

	// Print log
	b, err := json.Marshal(sl)
	if err != nil {
		return "{\"error\":\"error while marshalling log\"}"
	}

	return string(b)
}

func (logger ECS) Info(ctx Context, msg string, info ...AdditionalInfo) {
	// Add additional keys
	info = append(info, AdditionalInfo{"log.level", "info"})

	// Print log
	log.Print(logger.formatLog(ctx, msg, info...))
}

func (logger ECS) Error(ctx Context, msg string, info ...AdditionalInfo) {
	// Add additional keys
	info = append(info, AdditionalInfo{"log.level", "error"})

	// Print log
	log.Print(logger.formatLog(ctx, msg, info...))
}
