package loggers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// ECS is a logger that will print logs in Elastic Common Schema ECS format.
type ECS struct {
	level Level
}

// NewECS creates a new ECS logger.
func NewECS(options ...Option) ECS {
	ecs := ECS{}
	for _, option := range options {
		option(&ecs)
	}
	return ecs
}

func (ecs *ECS) setLevel(level Level) {
	ecs.level = level
}

func (ecs ECS) addStandardInfo(msg string, info ...extensions.LogInfo) []extensions.LogInfo {
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
		Value: "github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers/ecs.go",
	})

	// Return info
	return info
}

func (ecs ECS) formatLog(msg string, info ...extensions.LogInfo) string {
	// Set additional fields
	info = ecs.addStandardInfo(msg, info...)

	// Structure log
	sl := structureLogs(info)

	// Print log
	b, err := json.Marshal(sl)
	if err != nil {
		return "{\"error\":\"error while marshalling log\"}"
	}

	return string(b)
}

func (ecs ECS) logWithLevel(level string, msg string, info ...extensions.LogInfo) {
	// Add additional keys
	info = append(info, extensions.LogInfo{Key: "log.level", Value: level})

	// Print log
	fmt.Println(ecs.formatLog(msg, info...))
}

// Info logs a message at info level with context and additional info.
func (ecs ECS) Info(_ context.Context, msg string, info ...extensions.LogInfo) {
	if ecs.level > LevelInfo {
		return
	}
	ecs.logWithLevel("info", msg, info...)
}

// Warning logs a message at warning level with context and additional info.
func (ecs ECS) Warning(_ context.Context, msg string, info ...extensions.LogInfo) {
	if ecs.level > LevelWarning {
		return
	}
	ecs.logWithLevel("warning", msg, info...)
}

// Error logs a message at error level with context and additional info.
func (ecs ECS) Error(_ context.Context, msg string, info ...extensions.LogInfo) {
	ecs.logWithLevel("error", msg, info...)
}
