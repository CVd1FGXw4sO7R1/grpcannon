package report

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func makeBucketResults(durations ...time.Duration) []Result {
	out := make([]Result, len(durations))
	for i, d := range durations {
		out[i] = Result{Duration: d}
	}
	return out
}

func TestBucketize_Empty(t *testing.T) {
	if got := Bucketize(nil, 5); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBucketize_ZeroBuckets(t *testing.T) {
	res := makeBucketResults(ms(10), ms(20))
	if got := Bucketize(res, 0); got != nil {
		t.Fatalf("expected nil for n=0")
	}
}

func TestBucketize_AllErrors(t *testing.T) {
	res := []Result{
		{Duration: ms(10), Err: errors.New("fail")},
	}
	if got := Bucketize(res, 4); got != nil {
		t.Fatalf("expected nil when all errors")
	}
}

func TestBucketize_BucketCount(t *testing.T) {
	res := makeBucketResults(ms(10), ms(20), ms(30), ms(40), ms(50))
	buckets := Bucketize(res, 5)
	if len(buckets) != 5 {
		t.Fatalf("expected 5 buckets, got %d", len(buckets))
	}
}

func TestBucketize_CountsSum(t *testing.T) {
	res := makeBucketResults(ms(10), ms(20), ms(30), ms(40), ms(50))
	buckets := Bucketize(res, 4)
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != 5 {
		t.Fatalf("expected total count 5, got %d", total)
	}
}

func TestBucketize_UniformRange(t *testing.T) {
	res := makeBucketResults(ms(5), ms(5), ms(5))
	buckets := Bucketize(res, 3)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket for uniform input, got %d", len(buckets))
	}
	if buckets[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", buckets[0].Count)
	}
}

func TestWriteBucketize_Output(t *testing.T) {
	res := makeBucketResults(ms(10), ms(50), ms(100))
	var buf bytes.Buffer
	WriteBucketize(&buf, res, 3)
	out := buf.String()
	if out == "" || out == "no data\n" {
		t.Fatal("expected non-empty output")
	}
}

func TestWriteBucketize_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteBucketize(&buf, nil, 5)
	if buf.String() != "no data\n" {
		t.Fatalf("expected 'no data', got %q", buf.String())
	}
}
