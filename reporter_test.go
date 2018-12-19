package datadog

import "testing"

func TestGetPercentileNamesOutOfRange(t *testing.T) {
	_, err := getPercentileNames([]float64{0.23, 0})
	if err == nil {
		t.Fatal("Expected error")
	}
	_, err = getPercentileNames([]float64{0.23, 1})
	if err == nil {
		t.Fatal("Expected error")
	}
	_, err = getPercentileNames([]float64{0.23, -0.1})
	if err == nil {
		t.Fatal("Expected error")
	}
	_, err = getPercentileNames([]float64{0.23, 2})
	if err == nil {
		t.Fatal("Expected error")
	}
}
func TestGetPercentileNames(t *testing.T) {
	percentiles := []float64{0.23, 0.4, 0.99999, 0.45346356}
	names, err := getPercentileNames(percentiles)
	if err != nil {
		t.Fatal(err)
	}
	expectedNames := []string{".p23", ".p4", ".p99999", ".p45346356"}
	if len(expectedNames) != len(names) {
		t.Fatalf("Expected names: %v, got: %v", expectedNames, names)
	}
	for i, expectedName := range expectedNames {
		if names[i] != expectedName {
			t.Fatalf("Expected names: %v, got: %v", expectedNames, names)
		}
	}
}
