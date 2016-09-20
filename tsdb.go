package tslogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gonum/stat"
)

const (
	tsdbTickTimeSeconds = 10
)

var tsdbAggregators = map[string]func([]float64) (float64, error){
	"min":   tsdbMin,
	"max":   tsdbMax,
	"count": tsdbCount,
	"mean":  tsdbMean,
	"p75":   tsdbQuantile75,
	"p90":   tsdbQuantile90,
	"p95":   tsdbQuantile95,
	"p99":   tsdbQuantile99,
}

func NewTSDB(conf toml.Primitive) (*TSDB, error) {
	tsdb := &TSDB{
		buff:   make(map[string]*tsdbMetricScope),
		client: &http.Client{},
	}
	if err := toml.PrimitiveDecode(conf, tsdb); err != nil {
		return nil, err
	}
	return tsdb, nil
}

type TSDB struct {
	Host   string
	client *http.Client
	buff   map[string]*tsdbMetricScope
	once   sync.Once
	sync.Mutex
}

func (self *TSDB) Name() string {
	return "tsdb"
}

func (self *TSDB) flush() map[string]*tsdbMetricScope {
	defer self.Unlock()
	self.Lock()
	buff := self.buff
	self.buff = make(map[string]*tsdbMetricScope)
	return buff
}

func (self *TSDB) timer() {
	go func() {
		for {
			select {
			case <-time.After(tsdbTickTimeSeconds * time.Second):
				Log.Printf("[DEBUG] tsdb tick")
				self.send()
			}
		}
	}()
}

func (self *TSDB) Add(mName string, val float64, tags map[string]string, aggs []string) error {
	self.once.Do(self.timer)
	defer self.Unlock()
	key := mName + metricScopeKey(tags, aggs)
	self.Lock()
	if scope, ok := self.buff[key]; !ok {
		self.buff[key] = &tsdbMetricScope{
			Vals: make([]float64, 0),
			Tags: tags,
			Aggs: aggs,
			Name: mName,
		}
	} else {
		scope.Vals = append(scope.Vals, val)
	}
	return nil
}

func metricScopeKey(tags map[string]string, aggs []string) string {
	key := make([]string, len(tags)+len(aggs))
	key = append(key, aggs...)
	for k, _ := range tags {
		key = append(key, k)
	}
	sort.Strings(key)
	return strings.Join(key, "")
}

type tsdbMetricScope struct {
	Name string
	Vals []float64
	Tags map[string]string
	Aggs []string
}

func (self *TSDB) send() error {
	buff := self.flush()
	metrics := make([]*TSDBMetric, 0)
	var val float64
	var err error
	for _, scope := range buff {
		for _, agg := range scope.Aggs {
			if aggm, ok := tsdbAggregators[agg]; ok {
				val, err = aggm(scope.Vals)
				if err != nil {
					Log.Printf("[DEBUG] aggs %q return error: %v", agg, err)
					continue
				}
			} else {
				Log.Printf("[ERROR] can't find agg method %q", agg)
				continue
			}
			metrics = append(metrics, &TSDBMetric{
				Metric:    scope.Name + "." + agg,
				Value:     val,
				Tags:      scope.Tags,
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			})
		}
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		Log.Printf("[ERROR] can't marshal metrics, err: %v", err)
		return err
	}
	req, err := http.NewRequest("POST", "http://"+self.Host+"/api/put", bytes.NewBuffer(data))
	if err != nil {
		Log.Printf("[ERROR] can make request, err: %v", err)
		return err
	}
	resp, err := self.client.Do(req)
	if err != nil {
		Log.Printf("[ERROR] can't send request, err: %v", err)
		return err
	}
	defer resp.Body.Close()
	Log.Print("[DEBUG] data to OpenTSDB sended, %v", string(data))
	return nil
}

type TSDBMetric struct {
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
	Timestamp int64             `json:"timestamp"`
}

func tsdbMin(data []float64) (result float64, err error) {
	sort.Float64s(data)
	if len(data) > 0 {
		result = data[0]
	} else {
		err = fmt.Errorf("slice for tsdbMin is empty")
	}
	return
}

func tsdbMax(data []float64) (result float64, err error) {
	sort.Float64s(data)
	if len(data) > 0 {
		result = data[len(data)-1]
	} else {
		err = fmt.Errorf("slice for tsdbMax is empty")
	}
	return
}

func tsdbCount(data []float64) (float64, error) {
	return float64(len(data)), nil
}

func tsdbMean(data []float64) (result float64, err error) {
	if len(data) < 1 {
		err = fmt.Errorf("slice for tsdbMean is empty")
		return
	}
	result = stat.Mean(data, nil)
	return
}

func tsdbQuantile(q float64, vals []float64) (result float64, err error) {
	if len(vals) < 1 {
		err = fmt.Errorf("slice for tsdbQuantile%v is empty", q)
		return
	}
	sort.Float64s(vals)
	result = stat.Quantile(q, 1, vals, nil)
	return
}

func tsdbQuantile75(data []float64) (float64, error) {
	return tsdbQuantile(0.75, data)
}

func tsdbQuantile90(data []float64) (float64, error) {
	return tsdbQuantile(0.90, data)
}

func tsdbQuantile95(data []float64) (result float64, err error) {
	return tsdbQuantile(0.95, data)
}

func tsdbQuantile99(data []float64) (result float64, err error) {
	return tsdbQuantile(0.99, data)
}
