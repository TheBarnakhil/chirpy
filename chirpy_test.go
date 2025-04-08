package main

import "testing"

func TestReplaceProfanity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			input:    "testing to check if this damn kerfuffle in the sharbert is a valid fornax",
			expected: "testing to check if this damn **** in the **** is a valid ****",
		},
	}

	for i, testCase := range cases {
		result, err := filterProfaneWords(testCase.input)
		if err != nil {
			t.Errorf("Case %d: Expected %q, got %q", i, testCase.expected, err)
		}
		if result != testCase.expected {
			t.Errorf("Case %d: Expected %q, got %q", i, testCase.expected, result)
		}
	}
}
