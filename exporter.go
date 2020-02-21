package main

import (
	"encoding/json"
)

type Exporter struct {
	collector *Collector
}

func NewExporter(collector *Collector) Exporter {
	return Exporter{
		collector: collector,
	}
}

func (e *Exporter) Export() (string, error) {
	res, err := json.Marshal(e.collector.getExceptionAggregate())
	if err != nil {
		return "", err
	}
	return string(res), nil
}
