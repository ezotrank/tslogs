package tslogs

type Destination interface {
	Add(metricName string, val *Value, tags map[string]string, aggs []string) error
	Name() string
}