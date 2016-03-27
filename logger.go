package tslogs

import (
	"io"
	logger "log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	SetLogger(DEFAULT_LOG_LEVEL, "")
}

func execTime(fName string, logLevel string, startTime time.Time) {
	Log.Printf("[%s] %s, time: %v", logLevel, fName, time.Since(startTime))
}

func initLogFilter(logLevel string, logFile string) io.Writer {
	writer := os.Stderr
	if len(logFile) > 0 {
		var err error
		writer, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   writer,
	}
	return filter
}

func SetLogger(logLevel string, logFile string) {
	mutex := &sync.Mutex{}
	Log = &logger.Logger{}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for sig := range c {
			mutex.Lock()
			Log.Printf("[WARN] Go a %v Signal! Reopen logs", sig)
			Log.SetOutput(initLogFilter(logLevel, logFile))
			mutex.Unlock()
		}
	}()
	Log.SetOutput(initLogFilter(logLevel, logFile))
}
