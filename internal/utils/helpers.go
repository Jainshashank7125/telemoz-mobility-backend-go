package utils

// Float64ToIntPointer converts a float64 to *int (rounded)
func Float64ToIntPointer(f float64) *int {
	i := int(f)
	return &i
}
