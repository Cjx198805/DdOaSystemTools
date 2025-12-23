package util

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// IntValue returns the value of an int pointer or a default value
func IntValue(i *int, defaultVal int) int {
	if i == nil {
		return defaultVal
	}
	return *i
}
