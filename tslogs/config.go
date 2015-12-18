package tslogs

import (
	"regexp"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Host   string
	Port   int
	Groups map[string]*Group
}

func (self *Config) load() error {
	for _, group := range self.Groups {
		err := group.prepareRegexp()
		if err != nil {
			return err
		}
	}
	return nil
}

type Rule struct {
	Name   string
	Regexp *regexp.Regexp
}

type Group struct {
	Mask  string
	Rules [][]string
	rules []*Rule
}

func (self *Group) prepareRegexp() error {
	self.rules = make([]*Rule, 0)
	for _, rule := range self.Rules {
		r, err := regexp.Compile(rule[1])
		if err != nil {
			return err
		}
		self.rules = append(self.rules, &Rule{rule[0], r})
	}
	return nil
}

func LoadConfig(raw []byte) (*Config, error) {
	config := &Config{}
	_, err := toml.Decode(string(raw), config)
	if err != nil {
		return config, err
	}
	err = config.load()
	return config, err
}
