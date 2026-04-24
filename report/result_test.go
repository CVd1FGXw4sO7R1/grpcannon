package report

import (
	"errors"
	"testing"
)

func TestResult_IsSuccess_NoError(t *testing.T) {
	r := Result{}
	if !r.IsSuccess() {
		t.Error("expected IsSuccess true when Err is nil")
	}
}

func TestResult_IsSuccess_WithError(t *testing.T) {
	r := Result{Err: errors.New("boom")}
	if r.IsSuccess() {
		t.Error("expected IsSuccess false when Err is set")
	}
}
