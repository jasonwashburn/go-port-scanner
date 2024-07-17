package main

import (
	"slices"
	"testing"
)

func TestParsePortRange(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected []int
		errMsg   string
	}{
		{
			name:     "single, single, range",
			input:    "22,80,5-10",
			expected: []int{5, 6, 7, 8, 9, 10, 22, 80},
			errMsg:   "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: []int{},
			errMsg:   "",
		},
		{
			name:     "whitespace",
			input:    "   ",
			expected: []int{},
			errMsg:   "",
		},
		{
			name:     "empty commas",
			input:    ", , ,,",
			expected: []int{},
			errMsg:   "",
		},
		{
			name:     "letters",
			input:    "abc, 1b2",
			expected: []int{},
			errMsg:   "unable to process port string, cannot convert to int",
		},
		{
			name:     "invalid range",
			input:    "1,2,3-4-5",
			expected: []int{},
			errMsg:   "unable to process port string, invalid range",
		},
		{
			name:     "invalid start of range",
			input:    "1,2,a-5",
			expected: []int{},
			errMsg:   "unable to convert start of range to integer",
		},
		{
			name:     "invalid stop of range",
			input:    "1,2,5-a",
			expected: []int{},
			errMsg:   "unable to convert stop of range to integer",
		},
		{
			name:     "decreasing order range",
			input:    "1,2,10-5",
			expected: []int{},
			errMsg:   "start of range must be greater than stop",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParsePortRange(tc.input)
			if !slices.Equal(result, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != tc.errMsg {
				t.Errorf("expected error: %s, got %s", tc.errMsg, errMsg)
			}
		})
	}
}
