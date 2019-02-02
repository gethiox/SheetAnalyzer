package main

import (
	"regexp"
	"testing"
)

// helper, designed without TTD, implementer promised this function is free from bugs™
func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIndexGather(t *testing.T) {
	type testCase struct {
		regexp   string
		expected []int
	}

	var testCases = []testCase{
		{ // depth 0 range (basically whole string range)
			"",
			[]int{0, 0},
		},
		{ // same as above
			"abc",
			[]int{0, 3},
		},
		{ // "depth 0" starts at 0 and ends at 2 but submatch also starts at 0 and ends at 2
			"()",
			[]int{0, 2, 0, 2},
		},
		{ // similar as above but additional submatch (at depth 1) in other submatch (at depth 0)
			"(())",
			[]int{0, 4, 0, 4, 1, 3},
		},
		{ // same as above but with wider "depth 0" range (space padding)
			" (()) ",
			[]int{0, 6, 1, 5, 2, 4},
		},
		{ // there is no additional submatch group (because of escape characters)
			`\(nope\)`,
			[]int{0, 8},
		},
		{ // the most ultimate test-case I've ever seen
			//        0/     7/   13/16/ 20/ 24/     32/         44/      53/                   75/
			regexp:   `prefix (first), (?P<name>second (third inner) padding) \(escaped non-group\)`,
			expected: []int{0, 76, 7, 14, 16, 54, 32, 45}, // counted manually, I am surprised too
		},
	}

	for _, c := range testCases {
		regexp.MustCompile(c.regexp)

		indexes := ExprFindSubmatchIndex(c.regexp)

		equal := intSliceEqual(indexes, c.expected)
		if !equal {
			t.Errorf("\n    expr: \"%s\"\nexpected: %#v\n     got: %#v", c.regexp, c.expected, indexes)
		}

	}
}
