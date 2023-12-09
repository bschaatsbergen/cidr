package helper

// ContainsInt checks if the given slice of integers contains the specified integer.
// It returns true if the integer is within the slice, otherwise false.
func ContainsInt(ints []int, specifiedInt int) bool {
	for _, i := range ints {
		if i == specifiedInt {
			return true
		}
	}
	return false
}
