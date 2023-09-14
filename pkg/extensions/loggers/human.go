package loggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Human is a logger that will print logs in Elastic Common Schema format.
type Human struct {
	boldRedPrinter    *color.Color
	boldOrangePrinter *color.Color
	boldWhitePrinter  *color.Color
	greyPrinter       *color.Color
}

// NewHuman creates a new Human logger.
func NewHuman() Human {
	// Create red color
	red := color.New(color.FgHiRed)
	boldRed := red.Add(color.Bold)

	// Create orange color
	orange := color.New(color.FgHiYellow)
	boldOrange := orange.Add(color.Bold)

	// Create white color
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)

	return Human{
		boldRedPrinter:    boldRed,
		boldOrangePrinter: boldOrange,
		boldWhitePrinter:  boldWhite,
		greyPrinter:       color.New(color.FgHiBlack),
	}
}

func (lh Human) humanizeStructuredLogs(sl map[string]any, msgFmt *color.Color, prefixes ...string) string {
	var s string
	joinedPrefixes := strings.Join(prefixes, "")

	// Put timestamp and message first if it tsExists
	ts, tsExists := sl["@Timestamp"]
	msg, msgExists := sl["Message"]
	if tsExists && msgExists {
		s += msgFmt.Sprintf("> %s%s: %s\n", joinedPrefixes, ts, msg)
		delete(sl, "@Timestamp")
		delete(sl, "Message")
		return s + lh.humanizeStructuredLogs(sl, msgFmt, append(prefixes, "  ")...)
	}

	// Generate other keys
	for k, v := range sl {
		switch tv := v.(type) {
		case map[string]any:
			children := lh.humanizeStructuredLogs(tv, msgFmt, append(prefixes, "  ")...)
			s += lh.greyPrinter.Sprintf("%s%s:\n%s", joinedPrefixes, k, children)
		default:
			s += lh.greyPrinter.Sprintf("%s%s: %v\n", joinedPrefixes, k, tv)
		}
	}
	return s
}

func (lh Human) setInfoFromContext(ctx context.Context, msg string, info ...extensions.LogInfo) []extensions.LogInfo {
	// Add additional keys from context
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsChannel, func(value any) {
		info = append(info, extensions.LogInfo{Key: "Channel", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsCorrelationID, func(value any) {
		info = append(info, extensions.LogInfo{Key: "CorrelationID", Value: value})
	})
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsMessage, func(value any) {
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

func (lh Human) formatLog(ctx context.Context, msgFmt *color.Color, msg string, info ...extensions.LogInfo) string {
	// Set additional fields
	info = lh.setInfoFromContext(ctx, msg, info...)

	// Structure log
	sl := structureLogs(info)

	// Humanize structured logs
	return lh.humanizeStructuredLogs(sl, msgFmt)
}

// Info logs a message at info level with context and additional info.
func (lh Human) Info(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(lh.formatLog(ctx, lh.boldWhitePrinter, msg, info...))
}

// Warning logs a message at warning level with context and additional info.
func (lh Human) Warning(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(lh.formatLog(ctx, lh.boldOrangePrinter, msg, info...))
}

// Error logs a message at error level with context and additional info.
func (lh Human) Error(ctx context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(lh.formatLog(ctx, lh.boldRedPrinter, msg, info...))
}
