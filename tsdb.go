package tslogs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	MAX_METRICS_PER_CHUNK = 64
)

func NewOpenTSDB(host string, port int) *OpenTSDB {
	tsdb := &OpenTSDB{Host: host, Port: port, client: &http.Client{}}
	return tsdb
}

type OpenTSDB struct {
	Host   string
	Port   int
	client *http.Client
}

func (self *OpenTSDB) apiSendUrl() string {
	return "http://" + strings.Join([]string{self.Host, strconv.Itoa(self.Port)}, ":") + "/api/put"
}

func (self *OpenTSDB) Send(metrics []*Metric) error {
	for _, m := range metrics {
		m.Timestamp = m.IntTime()
		m.Tags = m.PrepareTags()
	}
	offset := 0
	for offset < len(metrics) {
		var chunk []*Metric
		if len(metrics)-offset > MAX_METRICS_PER_CHUNK {
			chunk = metrics[offset : offset+MAX_METRICS_PER_CHUNK]
		} else {
			chunk = metrics
		}
		offset += len(chunk)
		data, err := json.Marshal(chunk)
		if err != nil {
			Log.Printf("[ERROR] can't marshal metrics, err: %v", err)
			return err
		}
		Log.Printf("[DEBUG] send %d metrics %v", len(chunk), string(data))
		req, err := http.NewRequest("POST", self.apiSendUrl(), bytes.NewBuffer(data))
		if err != nil {
			Log.Printf("[ERROR] can make request, err: %v", err)
			return err
		}
		resp, err := self.client.Do(req)
		defer resp.Body.Close()
		if err != nil || resp.StatusCode >= 300 {
			Log.Printf("[ERROR] can't send request, code: %d, resp: %v, err: %v", resp.StatusCode, resp, err)
			return err
		}
		Log.Print("[DEBUG] OpenTSDB chunk sended")
	}
	return nil
}
