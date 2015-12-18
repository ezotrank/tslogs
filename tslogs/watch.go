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

	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
)

var (
	hostname string
	tstd     *OpenTSTD
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalf("can't get hostname, err: %v", err)
	}
}

func NewOpenTSTD(host string, port int) (*OpenTSTD, error) {
	tstd := &OpenTSTD{Host: host, Port: port}
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	var err error
	tstd.conn, err = net.Dial("tcp", addr)
	return tstd, err
}

type OpenTSTD struct {
	Host string
	Port int
	conn net.Conn
}

func (self *OpenTSTD) Send(m *Metric) error {
	msg := fmt.Sprintf("put %s %s %s %s\n", m.Name, m.StringTime(), m.Value, m.StringTags())
	_, err := self.conn.Write([]byte(msg))
	if err != nil {
		log.Errorf("can't send data to server, err: %v", err)
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

func tailFile(filePath string, rules []*Rule, wg *sync.WaitGroup) error {
	defer wg.Done()
	t, err := tail.TailFile(filePath, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
	if err != nil {
		log.Errorf("can't tail file %q err: %v", filePath, err)
		return err
	}
	log.Infof("start watching file %q", filePath)
	for line := range t.Lines {
		for _, rule := range rules {
			if rule.Regexp.Match([]byte(line.Text)) {
				matches := rule.Regexp.FindStringSubmatch(line.Text)
				if len(matches) == 0 {
					return nil
				}
				tags := make(map[string]interface{})
				var val string
				for i, value := range matches[1:] {
					if rule.Regexp.SubexpNames()[i+1] == "val" {
						val = value
					} else {
						tags[rule.Regexp.SubexpNames()[i+1]] = value
					}
				}
				t := time.Now()
				metric := &Metric{Name: rule.Name, Value: val, Time: &t, Tags: tags}
				tstd.Send(metric)
			}
		}
	}
	return nil
}

func Watch(config *Config) error {
	var err error
	tstd, err = NewOpenTSTD(config.Host, config.Port)
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
			go tailFile(fPath, group.rules, wg)
		}
	}
	log.Info("watching...")
	wg.Wait()
	return nil
}
