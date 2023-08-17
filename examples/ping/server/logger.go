package main

import (
	"fmt"
	"log"
)

type SimpleLogger struct{}

func (logger SimpleLogger) formatlogInfo(keyvals ...interface{}) string {
	var formattedLogInfo string
	for i := 0; i < len(keyvals)-1; i += 2 {
		formattedLogInfo = fmt.Sprintf("%s, %s: %+v", formattedLogInfo, keyvals[i], keyvals[i+1])
	}
	return formattedLogInfo
}

func (logger SimpleLogger) Info(msg string, keyvals ...interface{}) {
	log.Printf("INFO: %s%s", msg, logger.formatlogInfo(keyvals...))
}

func (logger SimpleLogger) Error(msg string, keyvals ...interface{}) {
	log.Printf("ERROR: %s%s", msg, logger.formatlogInfo(keyvals...))
}
