package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	apiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
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

func (logger ECS) formatLog(ctx context.Context, msg string, info ...AdditionalInfo) string {
	// Add additional keys from context
	apiContext.IfSet(ctx, apiContext.KeyIsModule, func(value any) {
		info = append(info, AdditionalInfo{"event.module", value})
	})
	apiContext.IfSet(ctx, apiContext.KeyIsProvider, func(value any) {
		info = append(info, AdditionalInfo{"event.provider", value})
	})
	apiContext.IfSet(ctx, apiContext.KeyIsChannel, func(value any) {
		info = append(info, AdditionalInfo{"event.action", value})
	})
	apiContext.IfSet(ctx, apiContext.KeyIsOperation, func(value any) {
		info = append(info, AdditionalInfo{"event.reason", value})
	})
	apiContext.IfSet(ctx, apiContext.KeyIsMessage, func(value any) {
		info = append(info, AdditionalInfo{"event.original", value})
	})
	apiContext.IfSet(ctx, apiContext.KeyIsCorrelationID, func(value any) {
		info = append(info, AdditionalInfo{"trace.id", value})
	})

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

func (logger ECS) Info(ctx context.Context, msg string, info ...AdditionalInfo) {
	// Add additional keys
	info = append(info, AdditionalInfo{"log.level", "info"})

	// Print log
	log.Print(logger.formatLog(ctx, msg, info...))
}

func (logger ECS) Error(ctx context.Context, msg string, info ...AdditionalInfo) {
	// Add additional keys
	info = append(info, AdditionalInfo{"log.level", "error"})

	// Print log
	log.Print(logger.formatLog(ctx, msg, info...))
}
