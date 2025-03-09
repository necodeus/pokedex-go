package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("expected %v, got %v", c.expected, actual)
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("expected %v, got %v", c.expected, actual)
			}
		}
	}
}
