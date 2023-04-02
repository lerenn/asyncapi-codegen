package main

import (
	"fmt"
	"log"
)

type SimpleLogger struct{}

func (logger SimpleLogger) formatKeyValues(keyvals ...interface{}) string {
	var formattedKeyValues string
	for i := 0; i < len(keyvals)-1; i += 2 {
		formattedKeyValues = fmt.Sprintf("%s, %s: %+v", formattedKeyValues, keyvals[i], keyvals[i+1])
	}
	return formattedKeyValues
}

func (logger SimpleLogger) Info(msg string, keyvals ...interface{}) {
	log.Printf("INFO: %s%s", msg, logger.formatKeyValues(keyvals...))
}

func (logger SimpleLogger) Error(msg string, keyvals ...interface{}) {
	log.Printf("ERROR: %s%s", msg, logger.formatKeyValues(keyvals...))
}
