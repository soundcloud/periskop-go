package periskop

import (
	"errors"
	"testing"
)

var aggregationKeyCases = []struct {
	expectedAggregationKey string
	stacktrace             []string
}{
	{"testingError@811c9dc5", []string{""}},
	{"testingError@8de5e669", []string{"line 0:", "division by zero"}},
	{"testingError@9d610c3f", []string{"line 0:", "division by zero", "line 1:", "test()", "line 4:", "checkTest()"}},
	{"testingError@9b5eca82", []string{"line 0:", "division by zero", "line 1:", "test()", "line 5:", "checkTest()"}},
}

func newMockErrorWithContext(stacktrace []string) ErrorWithContext {
	errorInstance := NewErrorInstance(errors.New("testingError"), stacktrace)
	return NewErrorWithContext(errorInstance, SeverityError, HTTPContext{})
}

func TestTypes_aggregationKey(t *testing.T) {
	for _, tt := range aggregationKeyCases {
		t.Run(tt.expectedAggregationKey, func(t *testing.T) {
			errorWithContext := newMockErrorWithContext(tt.stacktrace)
			resultAggregationKey := errorWithContext.aggregationKey()
			if resultAggregationKey != tt.expectedAggregationKey {
				t.Errorf("error in aggregationKey, expected: %s, got %s", tt.expectedAggregationKey, resultAggregationKey)
			}
		})
	}
}

func TestTypes_addError(t *testing.T) {
	errorWithContext := newMockErrorWithContext([]string{""})
	errorAggregate := NewErrorAggregate("error@hash", SeverityWarning)
	errorAggregate.addError(errorWithContext)
	if errorAggregate.TotalCount != 1 {
		t.Errorf("expected one error")
	}
	for i := 0; i < MaxErrors; i++ {
		errorAggregate.addError(errorWithContext)
	}
	if errorAggregate.TotalCount != MaxErrors+1 {
		t.Errorf("expected %v total errors", MaxErrors+1)
	}
	if len(errorAggregate.LatestErrors) != MaxErrors {
		t.Errorf("expected %v latest errors", MaxErrors)
	}
}
