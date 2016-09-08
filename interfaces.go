package tslogs

type Destination interface {
	Add(metricName string, val float64, tags map[string]string, aggs []string) error
	Name() string
}
