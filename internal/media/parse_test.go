package media

import (
	"testing"
)

func TestParseCliOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "String with only whitespace",
			input:    "   \n  \t  ",
			expected: []string{},
		},
		{
			name:     "Single line",
			input:    "line1",
			expected: []string{"line1"},
		},
		{
			name:     "Multiple lines",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Lines with empty lines in between",
			input:    "line1\n\nline3",
			expected: []string{"line1", "line3"},
		},
		{
			name:     "Lines with whitespace lines in between",
			input:    "line1\n  \t  \nline3",
			expected: []string{"line1", "line3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCliOutput(tt.input)
			
			// Check if the length of the result matches the expected length
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCliOutput() returned %d lines, expected %d", len(result), len(tt.expected))
				return
			}
			
			// Check if each line in the result matches the expected line
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("ParseCliOutput() line %d = %q, expected %q", i, line, tt.expected[i])
				}
			}
		})
	}
}