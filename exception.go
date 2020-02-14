package main

import (
	"github.com/google/uuid"
	"time"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type ErrorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	Severity       string             `json:"severity"`
	LatestErrors   []ErrorWithContext `json:"latest_errors"`
}

type ErrorWithContext struct {
	Error       ErrorInstance `json:"error"`
	UUID        uuid.UUID        `json:"uuid"`
	Timestamp   int64         `json:"timestamp"`
	Severity    Severity      `json:"severity"`
	HTTPContext HTTPContext   `json:"http_context"`
}

type ErrorInstance struct {
	Class      string   `json:"class"`
	Message    string   `json:"message"`
	Stacktrace []string `json:"stacktrace"`
	Cause      *error   `json:"cause"`
}

type HTTPContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
}

func NewErrorInstance(err *error, stacktrace []string) ErrorInstance {
	return ErrorInstance{
		Cause: err,
		Stacktrace: stacktrace,
	}
}

func NewErrorWithContext(errorInstance ErrorInstance, severity Severity) ErrorWithContext {
	return ErrorWithContext{
		Error: errorInstance,
		UUID: uuid.New(),
		Timestamp: time.Now().Unix(),
		Severity: severity,
		HTTPContext: HTTPContext{},
	}
}