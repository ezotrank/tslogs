package main

import (
	"flag"
	"io/ioutil"
	"strings"
  "fmt"
	"github.com/ezotrank/tslogs"
)

var version string

var (
	configFile = flag.String("config", "./config.conf", "config file")
	logLevel = flag.String("logging", "INFO", "log level DEBUG, INFO, WARN, ERROR")
	logFile = flag.String("log", "", "log file")
	showVersion = flag.Bool("version", false, "show version of build and exit")
	tags = &tslogs.Tags{}
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("Version: %q\n", version)
		return
	}
	tslogs.SetLogger(*logLevel, *logFile)
	for _,arg := range flag.Args() {
		tags.Add(strings.Split(arg,"=")[0], strings.Split(arg,"=")[1])
	}
	rawConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't read file, err: %v", err)
		panic(err)
	}
	config, err := tslogs.LoadConfig(rawConfig, tags)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't load config, err: %v", err)
		panic(err)
	}
	err = tslogs.Watch(config)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't run Watch, err: %v", err)
		panic(err)
	}
}
