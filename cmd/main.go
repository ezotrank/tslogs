package main

import (
	"flag"
	"os"

	"github.com/ezotrank/tslogs"
)

var version string

var (
	configFile = flag.String("config", "./config.conf", "config file")
	logLevel   = flag.String("logging", "INFO", "log level DEBUG, INFO, WARN, ERROR")
	logFile    = flag.String("log", "", "log file")
)

func main() {
	flag.Parse()
	tslogs.SetLogger(*logLevel, *logFile)
	log := tslogs.GetLogger()
	config, err := tslogs.LoadConfigFile(*configFile)
	if err != nil {
		log.Printf("[ERROR] can't load config, err: %v", err)
		os.Exit(1)
	}
	if err := tslogs.Watch(config); err != nil {
		log.Printf("[ERROR] can't run Watch, err: %v", err)
		os.Exit(1)
	}
}
