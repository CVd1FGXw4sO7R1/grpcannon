package report

import (
	"errors"
	"testing"
	"time"
)

func TestResult_IsSuccess_NoError(t *testing.T) {
	r := Result{Start: time.Now(), Duration: time.Millisecond, Err: nil}
	if !r.IsSuccess() {
		t.Error("expected success")
	}
}

func TestResult_IsSuccess_WithError(t *testing.T) {
	r := Result{Start: time.Now(), Duration: time.Millisecond, Err: errors.New("fail")}
	if r.IsSuccess() {
		t.Error("expected failure")
	}
}
