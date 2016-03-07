package tslogs

import (
	logger "log"
	"os"
	"time"

	"github.com/hashicorp/logutils"
)

const (
	DEFAULT_LOG_LEVEL = "INFO"
)

var (
	Log *logger.Logger
)

func init() {
	SetLogger(DEFAULT_LOG_LEVEL)
}

func execTime(fName string, logLevel string, startTime time.Time) {
	Log.Printf("[%s] %s, time: %v", logLevel, fName, time.Since(startTime))
}

func SetLogger(logLevel string) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	Log = &logger.Logger{}
	Log.SetOutput(filter)
}
