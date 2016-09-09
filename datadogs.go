package tslogs

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/DataDog/datadog-go/statsd"
)

func NewDataDogS(conf toml.Primitive) (*DataDogS, error) {
	dataDogS := &DataDogS{}
	if err := toml.PrimitiveDecode(conf, dataDogS); err != nil {
		return nil, err
	}
	var err error
	dataDogS.statsd, err = statsd.New(dataDogS.Host)
	if err != nil {
		return nil, err
	}
	return dataDogS, nil
}

type DataDogS struct {
	Host   string
	statsd *statsd.Client
}

func (self *DataDogS) formatTags(tags map[string]string) []string {
	out := make([]string, len(tags))
	for k, v := range tags {
		out = append(out, k+":"+v)
	}
	return out
}

func (self *DataDogS) Add(mName string, val float64, tags map[string]string, aggs []string) error {
	statsdTags := self.formatTags(tags)
	for _, agg := range aggs {
		var err error
		switch agg {
		default:
			err = fmt.Errorf("agg %q not found")
		case "timing":
			err = self.statsd.TimeInMilliseconds(mName+".timing", val, statsdTags, 1)
		case "incr":
			err = self.statsd.Incr(mName+".incr", statsdTags, 1)
		case "decr":
			err = self.statsd.Decr(mName+".decr", statsdTags, 1)
		case "gauge":
			err = self.statsd.Gauge(mName+".gauge", val, statsdTags, 1)
		case "count":
			err = self.statsd.Count(mName+".count", int64(val), statsdTags, 1)
		// case "set":
		// 	err = self.statsd.Set(mName+".set", int64(val), statsdTags, 1)
		case "hist":
			err = self.statsd.Histogram(mName+".hist", val, statsdTags, 1)
		}
		if err != nil {
			Log.Printf("[WARN] can't send metric to datadogs, err: %v", err)
		}
	}
	return nil
}

func (self *DataDogS) Name() string {
	return "datadogs"
}
