package helper

// ContainsInt checks if the given slice of integers contains the specified integer.
// It returns true if the integer is within the slice, otherwise false.
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
