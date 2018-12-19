package datadog

import "testing"

func TestGetPercentileNamesOutOfRange(t *testing.T) {
	reporter := &Reporter{}
	if UsePercentiles([]float64{0.23, 0})(reporter) == nil {
		t.Fatal("Expected error")
	}
	if UsePercentiles([]float64{0.23, 1})(reporter) == nil {
		t.Fatal("Expected error")
	}
	if UsePercentiles([]float64{0.23, -0.1})(reporter) == nil {
		t.Fatal("Expected error")
	}
	if UsePercentiles([]float64{0.23, 2})(reporter) == nil {
		t.Fatal("Expected error")
	}
}
func TestGetPercentileNames(t *testing.T) {
	reporter := &Reporter{}
	percentiles := []float64{0.23, 0.4, 0.99999, 0.45346356}
	err := UsePercentiles(percentiles)(reporter)
	if err != nil {
		t.Fatal(err)
	}
	expectedNames := []string{".p23", ".p4", ".p99999", ".p45346356"}
	if len(expectedNames) != len(reporter.p) {
		t.Fatalf("Expected names: %v, got: %v", expectedNames, reporter.p)
	}
	for i, expectedName := range expectedNames {
		if reporter.p[i] != expectedName {
			t.Fatalf("Expected names: %v, got: %v", expectedNames, reporter.p)
		}
	}
}
