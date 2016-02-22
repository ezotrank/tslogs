package tslogs

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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
	tstd := &OpenTSTD{Host: host, Port: port}
	return tstd
}

type OpenTSTD struct {
	Host  string
	Port  int
	conn  net.Conn
	mutex sync.Mutex
}

func (self *OpenTSTD) Connect() error {
	addr := strings.Join([]string{self.Host, strconv.Itoa(self.Port)}, ":")
	var err error
	self.conn, err = net.Dial("tcp", addr)
	return err
}

func (self *OpenTSTD) Send(m *Metric) error {
	msg := fmt.Sprintf("put %s %s %s %s\n", m.Name, m.StringTime(), m.Value, m.StringTags())
	_, err := self.conn.Write([]byte(msg))
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			defer self.mutex.Unlock()
			self.mutex.Lock()
			for {
				Log.Printf("[WARN] net.Error %v", err)
				if self.Connect() == nil {
					break
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}
	return err
}

type Metric struct {
	Name  string
	Value string
	Time  *time.Time
	Tags  map[string]interface{}
}

func (self *Metric) StringTime() string {
	nTime := self.Time.UTC().UnixNano()
	strTime := strconv.FormatInt(nTime, 10)
	return strTime[:10] + "." + strTime[10:13]
}

func (self *Metric) StringTags() string {
	out := []string{fmt.Sprintf("host=%s", hostname)}
	for k, v := range self.Tags {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(out, " ")
}

func getFiles(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func tailFile(filePath string, group *Group, wg *sync.WaitGroup) error {
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
			t := time.Now()
			metric := &Metric{Name: rule.Name, Value: val, Time: &t, Tags: tags}
			tstd.Send(metric)
		}
	}
	return nil
}

func Watch(config *Config) error {
	tstd = NewOpenTSTD(config.Host, config.Port)
	err := tstd.Connect()
	if err != nil {
		return err
	}
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
