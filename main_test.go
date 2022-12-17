package main

import "testing"

func TestParseResults(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []*GitCommit
	}{
		{
			name: "Single commit",
			input: `commit abc123
Author: John Doe <john@example.com>
Date:   Mon Jan 2 15:04:05 2006 -0700

Initial commit
`,
			expected: []*GitCommit{
				{
					Sha: "abc123",
					Headers: map[string]string{
						"Author": "John Doe <john@example.com>",
						"Date":   "Mon Jan 2 15:04:05 2006 -0700",
					},
					Message: "Initial commit\n",
				},
			},
		},
		{
			name: "Multiple commits",
			input: `commit abc123
Author: John Doe <john@example.com>
Date:   Mon Jan 2 15:04:05 2006 -0700

Initial commit

commit def456
Author: Jane Doe <jane@example.com>
Date:   Mon Jan 3 16:05:06 2007 -0800

Second commit
`,
			expected: []*GitCommit{
				{
					Sha: "abc123",
					Headers: map[string]string{
						"Author": "John Doe <john@example.com>",
						"Date":   "Mon Jan 2 15:04:05 2006 -0700",
					},
					Message: "Initial commit\n",
				},
				{
					Sha: "def456",
					Headers: map[string]string{
						"Author": "Jane Doe <jane@example.com>",
						"Date":   "Mon Jan 3 16:05:06 2007 -0800",
					},
					Message: "Second commit\n",
				},
			},
		},
		{
			name: "Empty input",
			input: `
`,
			expected: []*GitCommit{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := ParseResults(test.input)
			if len(output) != len(test.expected) {
				t.Errorf("Expected %d commits, got %d", len(test.expected), len(output))
			}
			for i, commit := range output {
				expectedCommit := test.expected[i]
				if commit.Sha != expectedCommit.Sha {
					t.Errorf("Expected SHA %q, got %q", expectedCommit.Sha, commit.Sha)
				}
				if commit.Message != expectedCommit.Message {
					t.Errorf("Expected message %q, got %q", expectedCommit.Message, commit.Message)
				}
				if len(commit.Headers) != len(expectedCommit.Headers) {
					t.Errorf("Expected %d headers, got %d", len(expectedCommit.Headers), len(commit.Headers))
				}
				for key, value := range commit.Headers {
					expectedValue, ok := expectedCommit.Headers[key]
					if !ok {
						t.Errorf("Unexpected header %q", key)
					}
					if value != expectedValue {
						t.Errorf("Expected header value %q for key %q, got %q", expectedValue, key, value)
					}
				}
			}
		})
	}
}
