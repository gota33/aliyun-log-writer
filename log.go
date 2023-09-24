package sls

import (
	"io"
	"log"
	"os"
)

var (
	debugMode = false
	logger    = log.New(io.Discard, "[LOG] ", log.Flags())
)

func SetDebug(debug bool) {
	debugMode = debug
	if debug {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(io.Discard)
	}
}
