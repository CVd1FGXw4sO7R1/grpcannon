package report

import (
	"errors"
	"testing"
)

func TestResult_IsSuccess_NoError(t *testing.T) {
	r := Result{}
	if !r.IsSuccess() {
		t.Fatal("expected IsSuccess true when no error")
	}
}

func TestResult_IsSuccess_WithError(t *testing.T) {
	r := Result{Error: errors.New("fail")}
	if r.IsSuccess() {
		t.Fatal("expected IsSuccess false when error set")
	}
}
