package tslogs

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/hpcloud/tail"
)

func getFiles(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func tailFile(filePath string, group *Group, wg *sync.WaitGroup) error {
	defer wg.Done()
	t, err := tail.TailFile(filePath, tail.Config{
		Poll:     true,
		Follow:   true,
		ReOpen:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
	if err != nil {
		Log.Printf("[WARN] can't tail file %q err: %v", filePath, err)
		return err
	}
	Log.Printf("[INFO] start watching file %q", filePath)
	for line := range t.Lines {
		for _, rule := range group.Rules {
			if len(rule.Match) > 0 && !strings.Contains(line.Text, rule.Match) {
				continue
			}
			if rule.regexp != nil {
				matches := rule.regexp.FindStringSubmatch(line.Text)
				if len(matches) > 1 {
					tags := make(map[string]string)
					vals := make(map[string]*Value)
					for i, value := range matches[1:] {
						val := &Value{val: value}
						switch groupName := rule.subexpNames[i+1]; {
						default:
							tags[rule.subexpNames[i+1]] = value
						case len(groupName) >= 3 && groupName[:3] == "val":
							vals[groupName] = val
						}
					}
					go addMetric(group, rule, vals, tags)
				}
			}
		}
	}
	return nil
}

func addMetric(group *Group, rule *Rule, vals map[string]*Value, tags map[string]string) error {
	for _, dst := range group.destinations {
		for groupName,val := range vals {
			name := rule.Name
			if groupName != "val" {
				name = name + "." + strings.Split(groupName, "_")[1]
			}
			dst.Add(name, val, tags, rule.Aggs[dst.Name()])
		}
	}
	return nil
}

func Watch(config *Config) error {
	wg := &sync.WaitGroup{}
	for _, group := range config.Groups {
		filePaths, err := getFiles(group.Mask)
		if err != nil {
			return err
		}
		for _, fPath := range filePaths {
			wg.Add(1)
			go tailFile(fPath, group, wg)
		}
	}
	wg.Wait()
	return nil
}
