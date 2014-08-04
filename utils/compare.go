package utils

import (
	"log"
)

func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		log.Print("Slices are not of equal length")
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
