package main

// interface inspired by regexp.Regex.FindSubmatchIndex but designed for regular expression itself than it's match instead
// expression should compile correctly first, otherwise behavior is undefined and may panic as well
func ExprFindSubmatchIndex(expr string) (returnIndexes []int) {
	returnIndexes = append(returnIndexes, 0, len(expr)) // depth 0

	var previous rune
	var depth int
	var depthIndexes = make(map[int][]int) // keeps

	for i := 0; i < len(expr); i++ {
		// todo: panic when depth < 0

		current := rune(expr[i])

		switch {
		case current == '(' && previous != '\\':
			depthIndexes[depth] = append(depthIndexes[depth], i)
			depth += 1
		case current == ')' && previous != '\\':
			depth -= 1
			depthIndexes[depth] = append(depthIndexes[depth], i+1)
		}
		previous = current
	}

	// todo: panic when depth != 0

	for i := 0; i <= len(depthIndexes); i++ {
		for j := 0; j < len(depthIndexes[i]); j += 2 {
			returnIndexes = append(returnIndexes, depthIndexes[i][j], depthIndexes[i][j+1])
		}
	}

	return
}
