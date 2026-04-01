package env

import (
	"errors"
	"os"
	"strings"
)

type Env struct {
	prefix string
}

type Opt struct {
	Prefix string
}

func Provider(o Opt) *Env {
	return &Env{
		prefix: o.Prefix,
	}
}

// ReadBytes is not supported by the env provider.
func (e *Env) ReadBytes() ([]byte, error) {
	return nil, errors.New("env provider does not support this method")
}

// Read reads all available environment variables into a key:value map
// and returns it.
func (e *Env) Read() (map[string]any, error) {
	// Collect the environment variable keys.
	var keys []string
	for _, k := range os.Environ() {
		if e.prefix != "" {
			if strings.HasPrefix(k, e.prefix) {
				keys = append(keys, k)
			}
		} else {
			keys = append(keys, k)
		}
	}

	mp := make(map[string]any)
	for _, k := range keys {
		parts := strings.SplitN(k, "=", 2)
		mp[parts[0]] = parts[1]
	}

	return mp, nil
}
