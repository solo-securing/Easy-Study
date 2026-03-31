package configloader

import (
	"fmt"
	"packages/configloader/maps"
	"reflect"
	"sync"

	"github.com/mitchellh/copystructure"
)

type ConfigLoader struct {
	confMap     map[string]any
	confMapFlat map[string]any
	keyMap      KeyMap
	mu          sync.RWMutex
}

type KeyMap map[string][]string

func New() *ConfigLoader {
	return &ConfigLoader{
		confMap:     make(map[string]any),
		confMapFlat: make(map[string]any),
		keyMap:      make(KeyMap),
	}
}

// Load takes a Provider that either provides a parsed config map[string]any
// in which case pa (Parser) can be nil, or raw bytes to be parsed, where a Parser
// can be provided to parse.
func (cl *ConfigLoader) Load(p Provider, pa Parser) error {
	var (
		mp  map[string]any
		err error
	)

	if p == nil {
		return fmt.Errorf("load received a nil provider")
	}

	// No Parser is given. Call the Provider's Read() method to get the config map.
	if pa == nil {
		mp, err = p.Read()
		if err != nil {
			return err
		}
	} else {
		// There's a Parser. Get raw bytes from the Provider to parse.
		b, err := p.ReadBytes()
		if err != nil {
			return err
		}
		mp, err = pa.Unmarshal(b)
		if err != nil {
			return err
		}
	}

	return cl.merge(mp)
}

func (cl *ConfigLoader) merge(c map[string]any) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	maps.IntfaceKeysToStrings(c)
	if err := maps.MergeStrict(c, cl.confMap); err != nil {
		return err
	}

	// Maintain a flattened version as well.
	cl.confMapFlat, cl.keyMap = maps.Flatten(cl.confMap, nil, ".")
	cl.keyMap = populateKeyParts(cl.keyMap, ".")

	return nil
}

// populateKeyParts iterates a key map and generates all possible
// traversal paths. For instance, `parent.child.key` generates
// `parent`, and `parent.child`.
func populateKeyParts(m KeyMap, delim string) KeyMap {
	out := make(KeyMap, len(m)) // The size of the result is at very least same to KeyMap
	for _, parts := range m {
		// parts is a slice of [parent, child, key]
		var nk string

		for i := range parts {
			if i == 0 {
				// On first iteration only use first part
				nk = parts[i]
			} else {
				// If nk already contains a part (e.g. `parent`) append delim + `child`
				nk += delim + parts[i]
			}
			if _, ok := out[nk]; ok {
				continue
			}
			out[nk] = make([]string, i+1)
			copy(out[nk], parts[0:i+1])
		}
	}
	return out
}

// Get returns the raw, uncast any value of a given key path
// in the config map. If the key path does not exist, nil is returned.
func (cl *ConfigLoader) Get(path string) any {
	// No path. Return the whole conf map.
	if path == "" {
		return cl.Raw()
	}

	// Does the path exist?
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	p, ok := cl.keyMap[path]
	if !ok {
		return nil
	}
	res := maps.Search(cl.confMap, p)

	// Non-reference types are okay to return directly.
	// Other types are "copied" with maps.Copy or json.Marshal
	// that change the numeric types to float64.

	switch v := res.(type) {
	case int, int8, int16, int32, int64, float32, float64, string, bool:
		return v
	case map[string]any:
		return maps.Copy(v)
	case nil:
		return nil
	}

	// Skip nil pointers before copying.
	if rv := reflect.ValueOf(res); rv.Kind() == reflect.Ptr && rv.IsNil() {
		return res
	}

	out, _ := copystructure.Copy(&res)
	if ptrOut, ok := out.(*any); ok {
		return *ptrOut
	}
	return out
}

// Raw returns a copy of the full raw conf map.
// Note that it uses maps.Copy to create a copy that uses
// json.Marshal which changes the numeric types to float64.
func (cl *ConfigLoader) Raw() map[string]any {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return maps.Copy(cl.confMap)
}
