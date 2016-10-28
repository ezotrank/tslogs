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
	checkConfig = flag.Bool("check", false, "check config and exit")
)

func main() {
	flag.Parse()
	tslogs.SetLogger(*logLevel, *logFile)
	log := tslogs.GetLogger()
	log.Printf("[INFO] Start with log level %s", *logLevel)
	config, err := tslogs.LoadConfigFile(*configFile)
	if err != nil {
		log.Printf("[ERROR] can't load config, err: %v", err)
		os.Exit(1)
	}
	if *checkConfig {
		log.Printf("[INFO] config has checked %q", *configFile)
		os.Exit(0)
	}
	if err := tslogs.Watch(config); err != nil {
		log.Printf("[ERROR] can't run Watch, err: %v", err)
		os.Exit(1)
	}
}
