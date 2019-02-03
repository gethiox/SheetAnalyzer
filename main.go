package main

import "regexp"

// interface inspired by regexp.Regex.FindSubmatchIndex but designed for regular expression itself than it's match instead
// expression should compile correctly first, otherwise behavior is undefined and may panic as well
func ExprFindSubmatchIndex(expr string) (returnIndexes []int) {
	returnIndexes = append(returnIndexes, 0, len(expr)) // depth 0

	var indexPairs [][2]int // keeps proper order of founded submatches

	var depth int                    // depth of submatch recursion level
	var counter int                  // global submatch occurence
	var toUpdate = make(map[int]int) // depth -> submatch occurrence

	var previous, current rune

	for i := 0; i < len(expr); i++ {
		// todo: panic when depth < 0

		current = rune(expr[i])

		switch {
		case current == '(' && previous != '\\':
			pair := [2]int{i, -1} // saving start index for new submatch
			indexPairs = append(indexPairs, pair)

			toUpdate[depth] = counter
			counter += 1
			depth += 1
		case current == ')' && previous != '\\':
			depth -= 1

			pair := &indexPairs[toUpdate[depth]] // read pair on current depth that is claimed to be closed
			pair[1] = i + 1                      // updating end index
			// "delete(toUpdate, depth)" could be used here to keep logic clean but saving CPU cycles sounds better
		}
		previous = current
	}

	// todo: panic when depth != 0

	for _, indexes := range indexPairs {
		returnIndexes = append(returnIndexes, indexes[0], indexes[1])
	}

	return
}

var regexpName = regexp.MustCompile(`\(\?P<([a-zA-Z_][a-zA-Z0-9_]*)>(.*)\)`)

func ExprSeparateSubmatchName(expr string) []int {
	match := regexpName.MatchString(expr)
	if !match {
		return []int{1, len(expr) - 1} // assume submatch without a name was passed, returns everything excluding parenthesizes
	}

	return regexpName.FindStringSubmatchIndex(expr)[2:]
}
