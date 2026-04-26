package env

import (
	"errors"
	"fmt"
	"strings"
)

// EnvSet represents a named collection of environment variables.
type EnvSet struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

// New creates a new empty EnvSet with the given name.
func New(name string) (*EnvSet, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("env set name must not be empty")
	}
	return &EnvSet{
		Name:      name,
		Variables: make(map[string]string),
	}, nil
}

// Set adds or updates a key-value pair in the EnvSet.
func (e *EnvSet) Set(key, value string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("key must not be empty")
	}
	if strings.ContainsAny(key, " =") {
		return fmt.Errorf("invalid key %q: must not contain spaces or '='" , key)
	}
	e.Variables[key] = value
	return nil
}

// Get retrieves the value for a given key.
func (e *EnvSet) Get(key string) (string, bool) {
	v, ok := e.Variables[key]
	return v, ok
}

// Delete removes a key from the EnvSet.
func (e *EnvSet) Delete(key string) {
	delete(e.Variables, key)
}

// Keys returns all keys in the EnvSet sorted alphabetically.
func (e *EnvSet) Keys() []string {
	keys := make([]string, 0, len(e.Variables))
	for k := range e.Variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ToExportLines returns the env set as a slice of "export KEY=VALUE" strings.
func (e *EnvSet) ToExportLines() []string {
	lines := make([]string, 0, len(e.Variables))
	for _, k := range e.Keys() {
		lines = append(lines, fmt.Sprintf("export %s=%s", k, e.Variables[k]))
	}
	return lines
}

// Merge merges another EnvSet into this one. Existing keys are overwritten.
func (e *EnvSet) Merge(other *EnvSet) {
	for k, v := range other.Variables {
		e.Variables[k] = v
	}
}
