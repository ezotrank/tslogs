package tslogs

import (
	"sort"

	"github.com/gonum/stat"
)

var (
	Aggregators = make(map[string]Aggregator, 0)
)

func init() {
	Aggregators = map[string]Aggregator{
		"min":   Min,
		"max":   Max,
		"count": Count,
		"mean":  Mean,
		"p75":   Quantile75,
		"p90":   Quantile90,
		"p95":   Quantile95,
		"p99":   Quantile99,
	}
}

func Min(data []float64) (float64, error) {
	var result float64
	sort.Float64s(data)
	if len(data) > 0 {
		result = data[0]
	}
	return result, nil
}

func Max(data []float64) (float64, error) {
	var result float64
	sort.Float64s(data)
	if len(data) > 0 {
		result = data[len(data)-1]
	}
	return result, nil
}

func Count(data []float64) (float64, error) {
	return float64(len(data)), nil
}

func Mean(data []float64) (float64, error) {
	return stat.Mean(data, nil), nil
}

func Quantile75(data []float64) (float64, error) {
	sort.Float64s(data)
	return stat.Quantile(0.75, 1, data, nil), nil
}

func Quantile90(data []float64) (float64, error) {
	sort.Float64s(data)
	return stat.Quantile(0.90, 1, data, nil), nil
}

func Quantile95(data []float64) (float64, error) {
	sort.Float64s(data)
	return stat.Quantile(0.95, 1, data, nil), nil
}

func Quantile99(data []float64) (float64, error) {
	sort.Float64s(data)
	return stat.Quantile(0.99, 1, data, nil), nil
}
