package loggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
	"github.com/fatih/color"
)

// Text is a logger that will print logs in Elastic Common Schema format.
type Text struct {
	boldRedPrinter    *color.Color
	boldOrangePrinter *color.Color
	boldWhitePrinter  *color.Color
	greyPrinter       *color.Color
}

// NewText creates a new Human logger.
func NewText() Text {
	// Create red color
	red := color.New(color.FgHiRed)
	boldRed := red.Add(color.Bold)

	// Create orange color
	orange := color.New(color.FgHiYellow)
	boldOrange := orange.Add(color.Bold)

	// Create white color
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)

	return Text{
		boldRedPrinter:    boldRed,
		boldOrangePrinter: boldOrange,
		boldWhitePrinter:  boldWhite,
		greyPrinter:       color.New(color.FgHiBlack),
	}
}

func (tl Text) humanizeStructuredLogs(sl map[string]any, msgFmt *color.Color, prefixes ...string) string {
	var s string
	joinedPrefixes := strings.Join(prefixes, "")

	// Put timestamp and message first if it tsExists
	ts, tsExists := sl["@Timestamp"]
	msg, msgExists := sl["Message"]
	if tsExists && msgExists {
		s += msgFmt.Sprintf("> %s%s: %s\n", joinedPrefixes, ts, msg)
		delete(sl, "@Timestamp")
		delete(sl, "Message")
		return s + tl.humanizeStructuredLogs(sl, msgFmt, append(prefixes, "  ")...)
	}

	// Generate other keys
	for k, v := range sl {
		switch tv := v.(type) {
		case map[string]any:
			children := tl.humanizeStructuredLogs(tv, msgFmt, append(prefixes, "  ")...)
			s += tl.greyPrinter.Sprintf("%s%s:\n%s", joinedPrefixes, k, children)
		default:
			s += tl.greyPrinter.Sprintf("%s%s: %v\n", joinedPrefixes, k, tv)
		}
	}
	return s
}

func (tl Text) setInfoFromContext(ctx context.Context, msg string, info ...extensions.LogInfo) []extensions.LogInfo {
	// Add additional keys from context
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsChannel, func(value any) {
		info = append(info, extensions.LogInfo{Key: "Channel", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsCorrelationID, func(value any) {
		info = append(info, extensions.LogInfo{Key: "CorrelationID", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsBrokerMessage, func(value any) {
		info = append(info, extensions.LogInfo{Key: "Content", Value: value})
	})

	// Add additional keys
	info = append(info, extensions.LogInfo{
		Key:   "Message",
		Value: msg,
	})
	info = append(info, extensions.LogInfo{
		Key:   "@Timestamp",
		Value: time.Now().UTC().Format(time.RFC3339Nano),
	})

	// Return info
	return info
}

func (tl Text) formatLog(ctx context.Context, msgFmt *color.Color, msg string, info ...extensions.LogInfo) string {
	// Set additional fields
	info = tl.setInfoFromContext(ctx, msg, info...)

	// Structure log
	sl := structureLogs(info)

	// Humanize structured logs
	return tl.humanizeStructuredLogs(sl, msgFmt)
}

// Info logs a message at info level with context and additional info.
func (tl Text) Info(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(tl.formatLog(ctx, tl.boldWhitePrinter, msg, info...))
}

// Warning logs a message at warning level with context and additional info.
func (tl Text) Warning(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(tl.formatLog(ctx, tl.boldOrangePrinter, msg, info...))
}

// Error logs a message at error level with context and additional info.
func (tl Text) Error(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(tl.formatLog(ctx, tl.boldRedPrinter, msg, info...))
}
