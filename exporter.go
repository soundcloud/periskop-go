package periskop

import (
	"encoding/json"
)

type ErrorExporter struct {
	collector *ErrorCollector
}

func NewErrorExporter(collector *ErrorCollector) ErrorExporter {
	return ErrorExporter{
		collector: collector,
	}
}

func (e *ErrorExporter) Export() (string, error) {
	res, err := json.Marshal(e.collector.getAggregatedErrors())
	if err != nil {
		return "", err
	}
	return string(res), nil
}
