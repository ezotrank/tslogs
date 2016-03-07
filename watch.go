package tslogs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"bytes"
	"net/http"
	"sync"
	"encoding/json"

	"github.com/hpcloud/tail"
)

var (
	hostname string
	tstd     *OpenTSTD
	NodeTags map[string]interface{}
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		Log.Printf("[ERROR] can't get hostname, err: %v", err)
		panic(err)
	}
}

func NewOpenTSTD(host string, port int) *OpenTSTD {
	tstd := &OpenTSTD{Host: host, Port: port, httpClient: &http.Client{}}
	return tstd
}

type OpenTSTD struct {
	Host  string
	Port  int
	httpClient *http.Client
}

func (self *OpenTSTD) apiSendUrl() string {
	return "http://" + strings.Join([]string{self.Host, strconv.Itoa(self.Port)}, ":") + "/api/put"
}

func (self *OpenTSTD) SendMultiple(metrics []*Metric) error {
	for _,m := range metrics {
		m.Timestamp = m.IntTime()
		m.PrepareTags()
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		Log.Printf("[ERROR] can't marshal metrics, err: %v", err)
		return err
	}
	Log.Printf("[DEBUG] send %v", string(data))
	req, err := http.NewRequest("POST", self.apiSendUrl(), bytes.NewBuffer(data))
	if err != nil {
		Log.Printf("[ERROR] can make request, err: %v", err)
		return err
	}
	resp, err := self.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		Log.Printf("[ERROR] can't send request, resp: %v, err: %v", resp, err)
		return err
	}
	Log.Print("[DEBUG] OpenTSTD chunk sended")
	return nil
}

type Metric struct {
	Metric string `json:"metric"`
	Value string `json:"value"`
	Tags map[string]interface{} `json:"tags"`
	Timestamp int64 `json:"timestamp"`
	time  *time.Time
}

func (self *Metric) IntTime() int64 {
	return self.time.UTC().UnixNano() / int64(time.Millisecond)
}

func (self *Metric) PrepareTags() {
	if _,ok := self.Tags["host"]; !ok {
		self.Tags["host"] = hostname
	}
}

func getFiles(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func tailFile(filePath string, group *Group, wg *sync.WaitGroup, ch chan *Metric) error {
	defer wg.Done()
	t, err := tail.TailFile(filePath, tail.Config{Poll: true, Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
	if err != nil {
		Log.Printf("[WARN] can't tail file %q err: %v", filePath, err)
		return err
	}
	Log.Printf("[INFO] start watching file %q", filePath)
	for line := range t.Lines {
		for _, rule := range group.rules {
			if len(rule.stringContains) > 0 {
				if !strings.Contains(line.Text, rule.stringContains) {
					continue
				}
			}
			matches := rule.Regexp.FindStringSubmatch(line.Text)
			if len(matches) == 0 {
				continue
			}
			tags := make(map[string]interface{})
			for k, v := range NodeTags {
				tags[k] = v
			}
			var val string
			for i, value := range matches[1:] {
				switch rule.SubexpNames[i+1] {
				default:
					tags[rule.SubexpNames[i+1]] = value
				case "val":
					val = value
				case "val_count":
					val = "1"
				}
			}
			if len(val) < 1 {
				val = "0"
			}
			t := time.Now()
			ch <- &Metric{Metric: rule.Name, Value: val, time: &t, Tags: tags}
		}
	}
	return nil
}

func startQueue(ch chan *Metric, tick uint) error {
	buff := make([]*Metric, 0)
	sendTime := time.Now().UnixNano() + (int64(tick) * int64(time.Millisecond))
	for m := range ch {
		buff = append(buff, m)
		if m.time.UnixNano() >= sendTime {
			tstd.SendMultiple(buff)
			buff = make([]*Metric, 0)
			sendTime = m.time.UnixNano() + (int64(tick) * int64(time.Millisecond))
		}
	}
	Log.Print("[DEBUG] Queue started")
	return nil
}

func Watch(config *Config) error {
	tstd = NewOpenTSTD(config.Host, config.Port)

	wg := &sync.WaitGroup{}
	ch := make(chan *Metric)
	go startQueue(ch, config.Tick)
	for _, group := range config.Groups {
		filePaths, err := getFiles(group.Mask)
		if err != nil {
			return err
		}
		for _, fPath := range filePaths {
			wg.Add(1)
			go tailFile(fPath, group, wg, ch)
		}
	}
	wg.Wait()
	return nil
}
