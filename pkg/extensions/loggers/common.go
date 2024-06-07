package loggers

import (
	"fmt"
	"strings"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
)

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
			// Explicitly drop the old value
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
