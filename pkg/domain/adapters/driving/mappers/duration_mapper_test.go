package mappers

import (
	"testing"
	"time"
)

func TestUnit_ToIso8601Duration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{duration: 0, expected: "PT0S"},
		{duration: time.Second, expected: "PT1S"},
		{duration: time.Minute, expected: "PT1M"},
		{duration: time.Hour, expected: "PT1H"},
		{duration: 24 * time.Hour, expected: "P1D"},
		{duration: 365 * 24 * time.Hour, expected: "P365D"},
		{duration: 2 * 365 * 24 * time.Hour, expected: "P730D"},
		{duration: 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second, expected: "P2DT3H4M5S"},
		{duration: time.Second + 250*time.Millisecond, expected: "PT1.25S"},
		{duration: -(time.Hour + 30*time.Minute), expected: "-PT1H30M"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if actual := toIso8601Duration(test.duration); actual != test.expected {
				t.Fatalf("toIso8601Duration(%v) = %q, expected %q", test.duration, actual, test.expected)
			}
		})
	}
}
