package tslogs

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/DataDog/datadog-go/statsd"
)

func NewDataDogS(conf toml.Primitive) (*DataDogS, error) {
	dataDogS := &DataDogS{}
	if err := toml.PrimitiveDecode(conf, dataDogS); err != nil {
		return nil, err
	}	
	return dataDogS, nil
}

type DataDogS struct {
	Host   string
	once sync.Once
	statsd *statsd.Client
}

func (self *DataDogS) connect() (err error) {
	self.statsd, err = statsd.New(self.Host)
	return err
}

func (self *DataDogS) formatTags(tags map[string]string) []string {
	out := make([]string, len(tags))
	for k, v := range tags {
		out = append(out, k+":"+v)
	}
	return out
}

func (self *DataDogS) Add(mName string, val *Value, tags map[string]string, aggs []string) error {
	Log.Printf("[DEBUG] add metrics name: %s, val: %+v, tags: %+v, aggs: %+v", mName, val, tags, aggs)
	self.once.Do(func(){
		if err := self.connect(); err != nil {
			Log.Printf("[ERROR] can't connect to datadog statsd server, err: %v", err)
			panic(err)
		}
	})
	statsdTags := self.formatTags(tags)
	for _, agg := range aggs {
		var err error
		switch agg {
		default:
			err = fmt.Errorf("agg %q not found")
		case "timing":
			var fVal float64
			if fVal, err = val.Float64(); err == nil {
				err = self.statsd.TimeInMilliseconds(mName+".timing", fVal, statsdTags, 1)
			}
		case "incr":
			err = self.statsd.Incr(mName+".incr", statsdTags, 1)
		case "decr":
			err = self.statsd.Decr(mName+".decr", statsdTags, 1)
		case "gauge":
			var fVal float64
			if fVal, err = val.Float64(); err == nil {
				err = self.statsd.Gauge(mName+".gauge", fVal, statsdTags, 1)
			}
		case "count":
			var fVal float64
			if fVal, err = val.Float64(); err == nil {
				err = self.statsd.Count(mName+".count", int64(fVal), statsdTags, 1)
			}
		case "set":
			var sVal string
			if sVal, err = val.String(); err == nil {
				err = self.statsd.Set(mName+".set", sVal, statsdTags, 1)
			}
		case "hist":
			var fVal float64
			if fVal, err = val.Float64(); err == nil {
				err = self.statsd.Histogram(mName+".hist", fVal, statsdTags, 1)
			}
		}
		if err != nil {
			Log.Printf("[WARN] can't send metric to datadogs, err: %v", err)
		}
		Log.Printf("[DEBUG] metrics sent name: %s, val: %+v, tags: %+v, aggs: %+v", mName, val, tags, aggs)
	}
	return nil
}

func (self *DataDogS) Name() string {
	return "datadogs"
}
