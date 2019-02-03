package main

// interface inspired by regexp.Regex.FindSubmatchIndex but designed for regular expression itself than it's match instead
// expression should compile correctly first, otherwise behavior is undefined and may panic as well
func ExprFindSubmatchIndex(expr string) (returnIndexes []int) {
	returnIndexes = append(returnIndexes, 0, len(expr)) // depth 0

	var indexPairs [][2]int // keeps proper order of founded submatches

	var depth int                   // depth of submatch recursion level
	var counter int                 // global submatch occurence
	var tracker = make(map[int]int) // depth -> submatch occurrence counter value

	var previous, current rune

	for i := 0; i < len(expr); i++ {
		// todo: panic when depth < 0

		current = rune(expr[i])

		switch {
		case current == '(' && previous != '\\':
			pair := [2]int{i, -1}
			indexPairs = append(indexPairs, pair)
			tracker[depth] = counter
			counter += 1
			depth += 1
		case current == ')' && previous != '\\':
			depth -= 1
			indexPairs[tracker[depth]][1] = i + 1
		}
		previous = current
	}

	// todo: panic when depth != 0

	for _, indexes := range indexPairs {
		returnIndexes = append(returnIndexes, indexes[0], indexes[1])
	}

	return
}

func ExprSeparateSubmatchName(expr string) ([]int, error) {
	return []int{}, nil
}
