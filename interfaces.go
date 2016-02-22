package tslogs

import (
	logger "log"
	"os"

	"github.com/hashicorp/logutils"
)

var (
	Log *logger.Logger
)

func init() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
		Writer:   os.Stderr,
	}
	Log = &logger.Logger{}
	Log.SetOutput(filter)
}
