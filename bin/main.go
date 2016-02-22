package main

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/ezotrank/tslogs"
)

var (
	configFile = flag.String("config", "./config.toml", "config file")
	dryRun     = flag.Bool("dry-run", false, "dry run")
	tags = make(map[string]interface{},0)
)

func main() {
	flag.Parse()
	for _,arg := range flag.Args() {
		tags[strings.Split(arg,"=")[0]] = strings.Split(arg,"=")[1]
	}
	rawConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't read file, err: %v", err)
		panic(err)
	}
	config, err := tslogs.LoadConfig(rawConfig)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't load config, err: %v", err)
		panic(err)
	}
	config.DryRun = *dryRun
	tslogs.NodeTags = tags
	err = tslogs.Watch(config)
	if err != nil {
		tslogs.Log.Printf("[ERROR] can't run Watch, err: %v", err)
		panic(err)
	}
}
