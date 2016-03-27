package tslogs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_TICK = "60s"
)

type Config struct {
	Host   string
	Port   int
	Groups map[string]*Group
	Tags   *Tags
}

func (self *Config) load() error {
	var err error
	for _, group := range self.Groups {
		if len(group.Tick) < 1 {
			group.Tick = DEFAULT_TICK
		}
		group.tick, err = time.ParseDuration(group.Tick)
		if err != nil {
			return err
		}
		for _, rule := range group.Rules {
			err = rule.Prepare()
			if err != nil {
				return err
			}
		}
		group.presetTags = self.Tags
	}
	return nil
}

type Rule struct {
	Name        string
	Regexp      string
	Match       string
	Aggs        []string
	subexpNames []string
	regexp      *regexp.Regexp
	aggs        map[string]Aggregator
}

func (self *Rule) Prepare() (err error) {
	if len(self.Regexp) > 0 {
		self.regexp, err = regexp.Compile(self.Regexp)
		if err != nil {
			return
		}
		self.subexpNames = self.regexp.SubexpNames()
	}
	self.aggs = make(map[string]Aggregator, 0)
	for _, m := range self.Aggs {
		if method, ok := Aggregators[m]; ok {
			self.aggs[m] = method
		} else {
			err = fmt.Errorf("Method %q doesn't exists", method)
			return
		}
	}
	return nil
}

type Group struct {
	Mask       string
	Rules      []*Rule
	Tick       string
	tick       time.Duration
	presetTags *Tags
}

func LoadConfig(raw []byte, tags *Tags) (*Config, error) {
	config := &Config{Tags: tags}
	_, err := toml.Decode(string(raw), config)
	if err != nil {
		return config, err
	}
	err = config.load()
	if err != nil {
		return config, err
	}
	return config, err
}
