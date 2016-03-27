package tslogs

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

func getFiles(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func tailFile(filePath string, group *Group, wg *sync.WaitGroup, tsdb *OpenTSDB) error {
	defer wg.Done()
	t, err := tail.TailFile(filePath, tail.Config{Poll: true, Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
	if err != nil {
		Log.Printf("[WARN] can't tail file %q err: %v", filePath, err)
		return err
	}
	Log.Printf("[INFO] start watching file %q", filePath)
	buff := make(map[*Rule][]*Metric, 0)
	mutex := &sync.Mutex{}
	go func() {
		c := time.Tick(group.tick)
		for now := range c {
			Log.Printf("[DEBUG] tick %v", now)
			mutex.Lock()
			allMetrics := make([]*Metric, 0)
			for rule, metrics := range buff {
				if len(rule.aggs) > 0 {
					allMetrics = append(allMetrics, aggregateMetrics(rule, metrics)...)
				} else {
					allMetrics = append(allMetrics, metrics...)
				}
			}
			buff = make(map[*Rule][]*Metric, 0)
			mutex.Unlock()
			if len(allMetrics) > 0 {
				go tsdb.Send(allMetrics)
			} else {
				Log.Printf("[DEBUG] nothing to send")
			}
		}
	}()
	for line := range t.Lines {
		for _, rule := range group.Rules {
			var val float64
			tags := &Tags{}
			if len(rule.Match) > 0 {
				if !strings.Contains(line.Text, rule.Match) {
					continue
				}
				val = float64(1)
			}
			if len(rule.Regexp) > 0 {
				matches := rule.regexp.FindStringSubmatch(line.Text)
				if len(matches) == 0 {
					Log.Printf("[DEBUG] regexp %q doesn't match string %q", rule.Regexp, line.Text)
					continue
				}
				Log.Printf("[DEBUG] regexp %q match string %q", rule.Regexp, line.Text)
				for i, value := range matches[1:] {
					switch rule.subexpNames[i+1] {
					default:
						tags.Add(rule.subexpNames[i+1], value)
					case "val":
						if s, err := strconv.ParseFloat(value, 64); err == nil {
							val = s
						} else {
							Log.Printf("[WARN] can't parse value %q to float64, err: %v", value, err)
							break
						}
					}
				}
			}
			if group.presetTags != nil {
				tags.Update(group.presetTags)
			}
			metric := &Metric{Metric: rule.Name, Value: val, time: time.Now(), tags: tags}
			Log.Printf("[DEBUG] create metric %v", metric)
			if _, ok := buff[rule]; !ok {
				buff[rule] = make([]*Metric, 0)
			}
			buff[rule] = append(buff[rule], metric)
		}
	}
	return nil
}

func Watch(config *Config) error {
	tsdb := NewOpenTSDB(config.Host, config.Port)

	wg := &sync.WaitGroup{}
	for _, group := range config.Groups {
		filePaths, err := getFiles(group.Mask)
		if err != nil {
			return err
		}
		for _, fPath := range filePaths {
			wg.Add(1)
			go tailFile(fPath, group, wg, tsdb)
		}
	}
	wg.Wait()
	return nil
}
