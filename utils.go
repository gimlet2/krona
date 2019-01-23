package main

// Set dummy set
type Set interface {
	Has(name string) bool
}

// SetImpl implements Set
type SetImpl struct {
	store *map[string]bool
}

// NewSet constructor for Set
func NewSet(a []string) Set {
	m := map[string]bool{}
	for _, s := range a {
		m[s] = true
	}
	return SetImpl{&m}
}

// Has checks that set has value
func (s SetImpl) Has(name string) bool {
	return (*s.store)[name]
}
