package main

import (
	"regexp"
	"testing"
)

// helper, designed without TDD, implementer promised this function is free from bugs™
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

func intSliceError(t *testing.T, expr string, expected, got []int) {
	t.Errorf("\n    expr: \"%s\"\nexpected: %#v\n     got: %#v", expr, expected, got)
}

func TestExprFindSubmatchIndex(t *testing.T) {
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
		{ // test for proper submatch order, sequence should be odered by "beginning of submatch", not by depth
			//        0/      8/     16/     24/27/     35/
			regexp:   `depth_0 (group_1(group_2)) (group_3)`,
			expected: []int{0, 36, 8, 26, 16, 25, 27, 36}, // expected order: depth_0, group_1, group_2, group_3
		},
		{ // continues above but in more complex way, many groups on same depth
			//        0/       10/       20/       30/       40/       50/ 54/
			regexp:   `main (g1 (g2 (g3) (g4)) (g5 (g6) (g7))) (g8 (g9) (g10))`, // expected order: as numbers goes
			expected: []int{0, 55, 5, 39, 9, 23, 13, 17, 18, 22, 24, 38, 28, 32, 33, 37, 40, 55, 44, 48, 49, 54},
			//                -- , g1   , g2   , g3    , g4    , g5    , g6    , g7    , g8    , g9    , g10
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
			intSliceError(t, c.regexp, c.expected, indexes)
		}

	}
}

func TestExprSeparateSubmatchName(t *testing.T) {
	type testCase struct {
		regexp   string
		expected []int
	}

	// happy path only, I assume regexp was compiled correctly before
	var testCases = []testCase{
		{
			`()`,
			[]int{1, 1},
		},
		{
			`(content)`,
			[]int{1, 8},
		},
		{
			`(?P<name>content)`,
			[]int{4, 8, 9, 16}, // name should be recognized and ranges returned separately (name and content)
		},
	}

	for _, c := range testCases {
		indexes := ExprSeparateSubmatchName(c.regexp)

		equal := intSliceEqual(indexes, c.expected)
		if !equal {
			intSliceError(t, c.regexp, c.expected, indexes)
		}
	}
}
