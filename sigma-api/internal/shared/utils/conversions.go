package utils

import (
	"strconv"
)

// StringToUint converts a string to uint, returns 0 if invalid
func StringToUint(s string) uint {
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(u)
}

// StringToInt converts a string to int, returns 0 if invalid
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
