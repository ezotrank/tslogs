package tslogs

import (
	"os"
	"time"
)

var (
	hostname string
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		Log.Printf("[ERROR] can't get hostname, err: %v", err)
	}
}

type Metric struct {
	Metric    string                 `json:"metric"`
	Value     float64                `json:"value"`
	Tags      map[string]interface{} `json:"tags"`
	tags      *Tags
	Timestamp int64 `json:"timestamp"`
	time      time.Time
}

func (self *Metric) IntTime() int64 {
	return self.time.UnixNano() / int64(time.Millisecond)
}

func (self *Metric) PrepareTags() map[string]interface{} {
	if _, err := self.tags.Get("host"); err != nil {
		self.tags.Add("host", hostname)
	}
	return self.tags.All()
}

func aggregateMetrics(rule *Rule, metrics []*Metric) []*Metric {
	outMetrics := make([]*Metric, 0)
	tagMetrics := make(map[string][]*Metric)
	for _, metric := range metrics {
		key := metric.tags.Key()
		if _, ok := tagMetrics[key]; !ok {
			tagMetrics[key] = make([]*Metric, 0)
		}
		tagMetrics[key] = append(tagMetrics[key], metric)
	}
	for _, groupedMetrics := range tagMetrics {
		vals := make([]float64, 0)
		for _, metric := range groupedMetrics {
			vals = append(vals, metric.Value)
		}
		tags := groupedMetrics[0].tags
		for name, method := range rule.aggs {
			val, err := method(vals)
			if err != nil {
				Log.Printf("[ERROR] can't make aggregation %q, err: %v", name, err)
			}
			outMetrics = append(outMetrics, &Metric{Metric: rule.Name + "." + name, Value: val, time: time.Now(), tags: tags})
		}
	}
	return outMetrics
}
