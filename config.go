package tslogs

import (
	"regexp"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_TICK_TIME = 1000
)

type Config struct {
	DryRun bool
	Host   string
	Port   int
	Tick  uint
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
	Name           string
	Regexp         *regexp.Regexp
	SubexpNames    []string
	stringContains string
}

type Group struct {
	Mask       string
	Rules      [][]string
	PresetTags map[string]interface{} `json:"preset_tags"`
	rules      []*Rule
}

func (self *Group) prepareRegexp() error {
	self.rules = make([]*Rule, 0)
	for _, rule := range self.Rules {
		r, err := regexp.Compile(rule[1])
		if err != nil {
			return err
		}
		stringContains := ""
		if len(rule) == 3 {
			stringContains = rule[2]
		}
		self.rules = append(self.rules, &Rule{rule[0], r, r.SubexpNames(), stringContains})
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
	if err != nil {
		return config, err
	}
	if config.Tick < 1 {
		config.Tick = DEFAULT_TICK_TIME
	}
	return config, err
}
