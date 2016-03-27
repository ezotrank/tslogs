package tslogs

import (
	"fmt"
	"sort"
	"strings"
)

type Tags struct {
	tags map[string]interface{}
}

func (self *Tags) init() {
	if len(self.tags) < 1 {
		self.tags = make(map[string]interface{})
	}
}

func (self *Tags) All() map[string]interface{} {
	return self.tags
}

func (self *Tags) Add(k string, v interface{}) {
	self.init()
	self.tags[k] = v
}

func (self *Tags) Get(k string) (interface{}, error) {
	if val, ok := self.tags[k]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("tag %q not found")
}

func (self *Tags) Key() string {
	names := make([]string, 0)
	for k, v := range self.tags {
		names = append(names, k+v.(string))
	}
	sort.Strings(names)
	return strings.Join(names, "")
}

func (self *Tags) Update(newTags *Tags) {
	self.init()
	for k, v := range newTags.tags {
		self.tags[k] = v
	}
}
