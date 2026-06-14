package loggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Text is a logger that will print logs in Elastic Common Schema format.
type Text struct {
	boldRedPrinter    *color.Color
	boldOrangePrinter *color.Color
	boldWhitePrinter  *color.Color
	greyPrinter       *color.Color
	level             Level
}

func (tl *Text) setLevel(level Level) {
	tl.level = level
}

// NewText creates a new Human logger.
func NewText(options ...Option) Text {
	// Create red color
	red := color.New(color.FgHiRed)
	boldRed := red.Add(color.Bold)

	// Create orange color
	orange := color.New(color.FgHiYellow)
	boldOrange := orange.Add(color.Bold)

	// Create white color
	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)

	tl := Text{
		boldRedPrinter:    boldRed,
		boldOrangePrinter: boldOrange,
		boldWhitePrinter:  boldWhite,
		greyPrinter:       color.New(color.FgHiBlack),
	}
	for _, option := range options {
		option(&tl)
	}
	return tl
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

func (tl Text) addStandardInfo(msg string, info ...extensions.LogInfo) []extensions.LogInfo {
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

func (tl Text) formatLog(msgFmt *color.Color, msg string, info ...extensions.LogInfo) string {
	// Set additional fields
	info = tl.addStandardInfo(msg, info...)

	// Structure log
	sl := structureLogs(info)

	// Humanize structured logs
	return tl.humanizeStructuredLogs(sl, msgFmt)
}

// Info logs a message at info level with context and additional info.
func (tl Text) Info(_ context.Context, msg string, info ...extensions.LogInfo) {
	if tl.level > LevelInfo {
		return
	}
	fmt.Println(tl.formatLog(tl.boldWhitePrinter, msg, info...))
}

// Warning logs a message at warning level with context and additional info.
func (tl Text) Warning(_ context.Context, msg string, info ...extensions.LogInfo) {
	if tl.level > LevelWarning {
		return
	}
	fmt.Println(tl.formatLog(tl.boldOrangePrinter, msg, info...))
}

// Error logs a message at error level with context and additional info.
func (tl Text) Error(_ context.Context, msg string, info ...extensions.LogInfo) {
	fmt.Println(tl.formatLog(tl.boldRedPrinter, msg, info...))
}
