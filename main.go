package main

import (
	"flag"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ezotrank/tslogs/tslogs"
)

const (
	VERSION = "0.1"
)

var (
	configFile = flag.String("config", "./config.toml", "config file")
	version    = flag.Bool("version", false, "version number")
	dryRun     = flag.Bool("dry-run", false, "dry run")
)

func main() {
	flag.Parse()

	if *version {
		log.Infof("tslogs version %s\n", VERSION)
		os.Exit(0)
	}
	rawConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("can't read file, err: %v", err)
	}
	config, err := tslogs.LoadConfig(rawConfig)
	if err != nil {
		log.Fatalf("can't load config, err: %v", err)
	}
	config.DryRun = *dryRun
	err = tslogs.Watch(config)
	if err != nil {
		log.Fatalf("can't run Watch, err: %v", err)
	}
}
