package state

import (
	"strings"

	"github.com/rotisserie/eris"
)

var (
	ErrVersionNotFound = eris.New("version not found")
)

// VersionMap maps a service and it's version.
type VersionMap map[string]string

// NewVersionMap initializes and returns an empty map.
func NewVersionMap() VersionMap {
	return make(VersionMap)
}

// Set creates an entry to the map.
// It a version already exists with the given service, it will be overwritten.
// Service and version will be smaller cases to stay consistent.
func (m VersionMap) Set(svc, version string) {
	svc = strings.ToLower(svc)
	version = strings.ToLower(version)
	m[svc] = version
}

// Get returns a version for the given service.
func (m VersionMap) Get(svc string) (string, error) {
	v, ok := m[svc]
	if v == "" || !ok {
		return "", ErrVersionNotFound
	}
	return v, nil
}

// Remove a service and it's version from the map.
func (m VersionMap) Remove(svc string) {
	delete(m, svc)
}

// Copy returns a deep copy of the original map.
// Golang map always modifies the original memory address.
// So, to make things safe to modify, we need to use a copy of the original map before modifying anything.
func (m VersionMap) Copy() VersionMap {
	newM := make(VersionMap)
	for k, v := range m {
		newM[k] = v
	}
	return newM
}
