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

type StringSet struct {
	values map[string]bool
}

func NewStringSet() StringSet {
	return StringSet{make(map[string]bool)}
}

func (c StringSet) Add(value string) {
	c.values[value] = true
}

func (c StringSet) AddMany(values []string) {
	for _, value := range values {
		c.values[value] = true
	}
}

func (c StringSet) Remove(value string) {
	delete(c.values, value)
}

func (c StringSet) Has(value string) bool {
	_, ok := c.values[value]
	return ok
}

func (c StringSet) GetAll() []string {
	all := make([]string, len(c.values))
	i := 0
	for k := range c.values {
		all[i] = k
		i = i + 1
	}
	return all
}
