package match

import (
	"testing"

	"github.com/agrski/greg/pkg/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	type test struct {
		name       string
		isBinary   bool
		text       string
		pattern    string
		expected   *Match
		expectedOk bool
	}

	tests := []test{
		{
			name:       "should ignore binary files",
			isBinary:   true,
			text:       "asdf",
			pattern:    "as",
			expected:   nil,
			expectedOk: false,
		},
		{
			name:       "should reject non-matching text file",
			isBinary:   false,
			text:       "asdf",
			pattern:    "foo",
			expected:   nil,
			expectedOk: false,
		},
		{
			name:       "should reject empty text file with non-empty pattern",
			isBinary:   false,
			text:       "",
			pattern:    "foo",
			expected:   nil,
			expectedOk: false,
		},
		{
			name:       "should reject empty text file with empty pattern",
			isBinary:   false,
			text:       "",
			pattern:    "",
			expected:   nil,
			expectedOk: false,
		},
		{
			name:     "should accept matching text file",
			isBinary: false,
			text:     "foo bar baz",
			pattern:  "bar",
			expected: &Match{
				Positions: []*FilePosition{
					{
						Line:        0,
						ColumnStart: 4,
					},
				},
			},
			expectedOk: true,
		},
		{
			name:     "should accept matching multi-line text file",
			isBinary: false,
			text: `first
second

fourth
foo
			`,
			pattern: "foo",
			expected: &Match{
				Positions: []*FilePosition{
					{
						Line:        4,
						ColumnStart: 0,
					},
				},
			},
			expectedOk: true,
		},
		{
			name:     "should accept multiple matches in multi-line text file",
			isBinary: false,
			text: `first
second foo

fourth
foo fifth
			`,
			pattern: "foo",
			expected: &Match{
				Positions: []*FilePosition{
					{
						Line:        1,
						ColumnStart: 7,
					},
					{
						Line:        4,
						ColumnStart: 0,
					},
				},
			},
			expectedOk: true,
		},
		{
			name:     "should accept multiple matches on same line",
			isBinary: false,
			text:     "foo bar foo",
			pattern:  "foo",
			expected: &Match{
				Positions: []*FilePosition{
					{
						Line:        0,
						ColumnStart: 0,
					},
					{
						Line:        0,
						ColumnStart: 8,
					},
				},
			},
			expectedOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo := &types.FileInfo{}
			fileInfo.IsBinary = tt.isBinary
			fileInfo.Text = tt.text

			matcher := newExactMatcher(zerolog.Nop())

			actual, ok := matcher.Match(tt.pattern, fileInfo)

			require.Equal(t, tt.expectedOk, ok)
			require.Equal(t, tt.expected, actual)
		})
	}
}
