package metric

// Metric is internal struct for operating metric values within the exporter
type Metric struct {
	Name   string
	Value  interface{}
	Labels map[string]string
}

// NewMetric creates new internal metric struct
func NewMetric(name string, value interface{}, tags map[string]string) Metric {
	return Metric{Name: name, Value: value, Labels: tags}
}
