package tslogs

import (
	"path/filepath"
	"strconv"
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
					var val float64
					for i, value := range matches[1:] {
						switch rule.subexpNames[i+1] {
						default:
							tags[rule.subexpNames[i+1]] = value
						case "val":
							if val, err = strconv.ParseFloat(value, 64); err != nil {
								Log.Printf("[WARN] can't parse value %q to float64, err: %v", value, err)
								break
							}
						}
					}
					go addMetric(group, rule, val, tags)
				}
			}
		}
	}
	return nil
}

func addMetric(group *Group, rule *Rule, val float64, tags map[string]string) error {
	for _, dst := range group.destinations {
		dst.Add(rule.Name, val, tags, rule.Aggs[dst.Name()])
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
