package tslogs

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
)

func LoadConfigFile(fpath string) (*Config, error) {
	configFile := &ConfigFile{}
	_, err := toml.DecodeFile(fpath, configFile)
	if err != nil {
		return nil, err
	}
	return ConfigFile2Config(configFile)
}

func loadDest(name string, conf toml.Primitive) (Destination, error) {
	switch name {
	case "tsdb":
		return NewTSDB(conf)
	case "datadogs":
		return NewDataDogS(conf)
	default:
		return nil, fmt.Errorf("destination with name %q doesn't exists", name)
	}
}

func ConfigFile2Config(cfile *ConfigFile) (*Config, error) {
	config := &Config{
		Destinations: make(map[string]Destination),
		Tags:         make([]string, 0),
		Groups:       make(map[string]*Group),
	}
	for name, dstConfig := range cfile.Destinations {
		if dst, err := loadDest(name, dstConfig); err != nil {
			return nil, fmt.Errorf("can't load dst %q, err: %v", name, err)
		} else {
			config.Destinations[name] = dst
		}
	}
	config.Tags = cfile.Tags
	for name, prm := range cfile.Groups {
		group := &Group{destinations: make([]Destination, 0)}
		if err := toml.PrimitiveDecode(prm, group); err != nil {
			return nil, err
		}
		for _, dst := range group.Destinations {
			if val, ok := config.Destinations[dst]; ok {
				group.destinations = append(group.destinations, val)
			} else {
				return nil, fmt.Errorf("can't find destination with name %q", dst)
			}
		}
		if err := group.LoadRules(); err != nil {
			return nil, err
		}
		config.Groups[name] = group
	}
	return config, nil
}

type ConfigFile struct {
	Tags         []string
	Destinations map[string]toml.Primitive
	Groups       map[string]toml.Primitive
}

type Config struct {
	Destinations map[string]Destination
	Tags         []string
	Groups       map[string]*Group
}

type Rule struct {
	Name        string
	Regexp      string
	Match       string
	Aggs        map[string][]string
	subexpNames []string
	regexp      *regexp.Regexp
}

func (self *Rule) Load() error {
	var err error
	if len(self.Regexp) > 0 {
		if self.regexp, err = regexp.Compile(self.Regexp); err != nil {
			return err
		}
		self.subexpNames = self.regexp.SubexpNames()
	}
	return nil
}

type Group struct {
	Mask         string
	Rules        []*Rule
	Destinations []string
	destinations []Destination
}

func (self *Group) LoadRules() error {
	for _, rule := range self.Rules {
		if err := rule.Load(); err != nil {
			return err
		}
	}
	return nil
}
