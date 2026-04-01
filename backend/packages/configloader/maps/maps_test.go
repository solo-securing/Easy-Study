package maps

import (
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name         string
		input        map[string]any
		delim        string
		expectedKeys map[string][]string
		expectedLen  int
	}{
		{
			name: "simple nested map",
			input: map[string]any{
				"parent": map[string]any{
					"child": 123,
				},
			},
			delim: ".",
			expectedKeys: map[string][]string{
				"parent.child": {"parent", "child"},
			},
			expectedLen: 1,
		},
		{
			name: "deeply nested map",
			input: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": "value",
					},
				},
			},
			delim: ".",
			expectedKeys: map[string][]string{
				"level1.level2.level3": {"level1", "level2", "level3"},
			},
			expectedLen: 1,
		},
		{
			name: "multiple keys at same level",
			input: map[string]any{
				"parent": map[string]any{
					"child1": 1,
					"child2": 2,
				},
			},
			delim:       ".",
			expectedLen: 2,
		},
		{
			name: "empty nested map",
			input: map[string]any{
				"parent": map[string]any{},
			},
			delim:       ".",
			expectedLen: 1,
		},
		{
			name: "mixed nested and flat",
			input: map[string]any{
				"flat": "value",
				"nested": map[string]any{
					"key": "val",
				},
			},
			delim:       ".",
			expectedLen: 2,
		},
		{
			name: "custom delimiter",
			input: map[string]any{
				"parent": map[string]any{
					"child": 123,
				},
			},
			delim: "-",
			expectedKeys: map[string][]string{
				"parent-child": {"parent", "child"},
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, keyMap := Flatten(tt.input, []string{}, tt.delim)

			if len(result) != tt.expectedLen {
				t.Errorf("expected %d keys, got %d", tt.expectedLen, len(result))
			}

			for expectedKey, expectedParts := range tt.expectedKeys {
				actualParts, ok := keyMap[expectedKey]
				if !ok {
					t.Errorf("expected key %q not found in result", expectedKey)
					continue
				}

				if len(actualParts) != len(expectedParts) {
					t.Errorf("key %q: expected %d parts, got %d", expectedKey, len(expectedParts), len(actualParts))
					continue
				}

				for i, part := range expectedParts {
					if actualParts[i] != part {
						t.Errorf("key %q part %d: expected %q, got %q", expectedKey, i, part, actualParts[i])
					}
				}
			}
		})
	}
}

func TestUnflatten(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		delim    string
		validate func(result map[string]any) bool
	}{
		{
			name: "simple unflatten",
			input: map[string]any{
				"parent.child": 123,
			},
			delim: ".",
			validate: func(result map[string]any) bool {
				parent, ok := result["parent"]
				if !ok {
					return false
				}
				parentMap, ok := parent.(map[string]any)
				if !ok {
					return false
				}
				child, ok := parentMap["child"]
				if !ok {
					return false
				}
				return child == 123
			},
		},
		{
			name: "deep unflatten",
			input: map[string]any{
				"a.b.c.d": "value",
			},
			delim: ".",
			validate: func(result map[string]any) bool {
				a, ok := result["a"].(map[string]any)
				if !ok {
					return false
				}
				b, ok := a["b"].(map[string]any)
				if !ok {
					return false
				}
				c, ok := b["c"].(map[string]any)
				if !ok {
					return false
				}
				return c["d"] == "value"
			},
		},
		{
			name: "multiple keys",
			input: map[string]any{
				"parent.child1": 1,
				"parent.child2": 2,
			},
			delim: ".",
			validate: func(result map[string]any) bool {
				parent, ok := result["parent"].(map[string]any)
				if !ok {
					return false
				}
				return parent["child1"] == 1 && parent["child2"] == 2
			},
		},
		{
			name: "no delimiter",
			input: map[string]any{
				"key": "value",
			},
			delim: "",
			validate: func(result map[string]any) bool {
				return result["key"] == "value"
			},
		},
		{
			name: "mixed hierarchy",
			input: map[string]any{
				"a.b":     1,
				"a.c.d":   2,
				"e":       3,
				"f.g.h.i": 4,
			},
			delim: ".",
			validate: func(result map[string]any) bool {
				// Check all keys exist and have correct structure
				a, ok := result["a"].(map[string]any)
				if !ok {
					return false
				}
				if a["b"] != 1 {
					return false
				}
				c, ok := a["c"].(map[string]any)
				if !ok {
					return false
				}
				if c["d"] != 2 {
					return false
				}
				if result["e"] != 3 {
					return false
				}
				return result["f"] != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unflatten(tt.input, tt.delim)
			if !tt.validate(result) {
				t.Errorf("validation failed for unflatten result: %v", result)
			}
		})
	}
}

func TestFlattenUnflattenRoundTrip(t *testing.T) {
	original := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]any{
				"user":     "admin",
				"password": "secret",
			},
		},
		"app": map[string]any{
			"name": "myapp",
		},
	}

	// Flatten and unflatten
	flattened, _ := Flatten(Copy(original), []string{}, ".")
	unflattened := Unflatten(flattened, ".")

	// Validate structure
	db, ok := unflattened["database"].(map[string]any)
	if !ok {
		t.Fatal("database key not found or not a map")
	}

	if db["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", db["host"])
	}

	if db["port"] != 5432 {
		t.Errorf("expected port=5432, got %v", db["port"])
	}

	creds, ok := db["credentials"].(map[string]any)
	if !ok {
		t.Fatal("credentials key not found or not a map")
	}

	if creds["user"] != "admin" {
		t.Errorf("expected user=admin, got %v", creds["user"])
	}
}

func TestMergeStrict(t *testing.T) {
	tests := []struct {
		name      string
		a         map[string]any
		b         map[string]any
		shouldErr bool
		validate  func(result map[string]any) bool
	}{
		{
			name: "simple merge",
			a: map[string]any{
				"key1": "value1",
			},
			b: map[string]any{
				"key2": "value2",
			},
			shouldErr: false,
			validate: func(result map[string]any) bool {
				return result["key1"] == "value1" && result["key2"] == "value2"
			},
		},
		{
			name: "overwrite value",
			a: map[string]any{
				"key": "new",
			},
			b: map[string]any{
				"key": "old",
			},
			shouldErr: false,
			validate: func(result map[string]any) bool {
				return result["key"] == "new"
			},
		},
		{
			name: "merge nested maps",
			a: map[string]any{
				"parent": map[string]any{
					"child1": 1,
				},
			},
			b: map[string]any{
				"parent": map[string]any{
					"child2": 2,
				},
			},
			shouldErr: false,
			validate: func(result map[string]any) bool {
				parent, ok := result["parent"].(map[string]any)
				if !ok {
					return false
				}
				return parent["child1"] == 1 && parent["child2"] == 2
			},
		},
		{
			name: "type mismatch error",
			a: map[string]any{
				"key": "string",
			},
			b: map[string]any{
				"key": 123,
			},
			shouldErr: true,
		},
		{
			name: "type mismatch in nested map",
			a: map[string]any{
				"parent": map[string]any{
					"child": "string",
				},
			},
			b: map[string]any{
				"parent": map[string]any{
					"child": 123,
				},
			},
			shouldErr: true,
		},
		{
			name: "replace non-map with map",
			a: map[string]any{
				"key": map[string]any{
					"nested": "value",
				},
			},
			b: map[string]any{
				"key": "scalar",
			},
			shouldErr: false,
			validate: func(result map[string]any) bool {
				_, ok := result["key"].(map[string]any)
				return ok
			},
		},
		{
			name: "empty maps",
			a:    map[string]any{},
			b:    map[string]any{},
			validate: func(result map[string]any) bool {
				return len(result) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of b since MergeStrict mutates it
			bCopy := Copy(tt.b)
			err := MergeStrict(tt.a, bCopy)

			if (err != nil) != tt.shouldErr {
				t.Errorf("expected error=%v, got error=%v", tt.shouldErr, err != nil)
			}

			if !tt.shouldErr && tt.validate != nil {
				if !tt.validate(bCopy) {
					t.Errorf("validation failed for merge result: %v", bCopy)
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]any
		verify func(original, copied map[string]any) bool
	}{
		{
			name: "simple copy",
			input: map[string]any{
				"key": "value",
			},
			verify: func(original, copied map[string]any) bool {
				return copied["key"] == original["key"]
			},
		},
		{
			name: "nested map copy",
			input: map[string]any{
				"parent": map[string]any{
					"child": "value",
				},
			},
			verify: func(original, copied map[string]any) bool {
				origParent := original["parent"].(map[string]any)
				copiedParent := copied["parent"].(map[string]any)
				return copiedParent["child"] == origParent["child"]
			},
		},
		{
			name: "deep copy independence",
			input: map[string]any{
				"parent": map[string]any{
					"child": "original",
				},
			},
			verify: func(original, copied map[string]any) bool {
				// Modify the copied map's nested value
				copiedParent := copied["parent"].(map[string]any)
				copiedParent["child"] = "modified"

				// Check that original is unchanged
				origParent := original["parent"].(map[string]any)
				return origParent["child"] == "original"
			},
		},
		{
			name: "copy with slice values",
			input: map[string]any{
				"list": []any{1, 2, 3},
			},
			verify: func(original, copied map[string]any) bool {
				return len(copied["list"].([]any)) == 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := Copy(tt.input)
			if !tt.verify(tt.input, copied) {
				t.Errorf("copy verification failed")
			}
		})
	}
}

func TestSearch(t *testing.T) {
	testMap := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]any{
				"user":     "admin",
				"password": "secret",
			},
		},
		"app": map[string]any{
			"name": "myapp",
		},
	}

	tests := []struct {
		name      string
		path      []string
		expected  any
		shouldEnd bool
	}{
		{
			name:      "first level key",
			path:      []string{"database"},
			shouldEnd: false,
		},
		{
			name:      "second level string value",
			path:      []string{"database", "host"},
			expected:  "localhost",
			shouldEnd: true,
		},
		{
			name:      "second level numeric value",
			path:      []string{"database", "port"},
			expected:  5432,
			shouldEnd: true,
		},
		{
			name:      "third level value",
			path:      []string{"database", "credentials", "user"},
			expected:  "admin",
			shouldEnd: true,
		},
		{
			name:      "nested map path",
			path:      []string{"database", "credentials"},
			shouldEnd: false,
		},
		{
			name:      "non-existent key at first level",
			path:      []string{"nonexistent"},
			expected:  nil,
			shouldEnd: true,
		},
		{
			name:      "non-existent nested key",
			path:      []string{"database", "nonexistent"},
			expected:  nil,
			shouldEnd: true,
		},
		{
			name:      "path through non-map value",
			path:      []string{"database", "host", "nested"},
			expected:  nil,
			shouldEnd: true,
		},
		{
			name:      "app name",
			path:      []string{"app", "name"},
			expected:  "myapp",
			shouldEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Search(testMap, tt.path)

			if tt.shouldEnd && result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			if !tt.shouldEnd && result == nil {
				t.Errorf("expected non-nil result for nested search")
			}
		})
	}
}

func TestIntfaceKeysToStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		validate func(result map[string]any) bool
	}{
		{
			name: "simple map with map[any]any",
			input: map[string]any{
				"key": map[any]any{
					"nested": "value",
				},
			},
			validate: func(result map[string]any) bool {
				nested, ok := result["key"].(map[string]any)
				if !ok {
					return false
				}
				return nested["nested"] == "value"
			},
		},
		{
			name: "mixed map[any]any keys",
			input: map[string]any{
				"config": map[any]any{
					"timeout": 30,
					1:         "numeric key",
					"string":  "value",
				},
			},
			validate: func(result map[string]any) bool {
				config, ok := result["config"].(map[string]any)
				if !ok {
					return false
				}
				return config["timeout"] == 30 && config["string"] == "value"
			},
		},
		{
			name: "nested map[any]any conversions",
			input: map[string]any{
				"level1": map[any]any{
					"level2": map[any]any{
						"key": "value",
					},
				},
			},
			validate: func(result map[string]any) bool {
				l1, ok := result["level1"].(map[string]any)
				if !ok {
					return false
				}
				l2, ok := l1["level2"].(map[string]any)
				if !ok {
					return false
				}
				return l2["key"] == "value"
			},
		},
		{
			name: "slice containing map[any]any",
			input: map[string]any{
				"items": []any{
					map[any]any{
						"id":   1,
						"name": "item1",
					},
					map[any]any{
						"id":   2,
						"name": "item2",
					},
				},
			},
			validate: func(result map[string]any) bool {
				items, ok := result["items"].([]any)
				if !ok {
					return false
				}
				firstItem, ok := items[0].(map[string]any)
				if !ok {
					return false
				}
				return firstItem["id"] == 1 && firstItem["name"] == "item1"
			},
		},
		{
			name: "already converted map[string]any",
			input: map[string]any{
				"key": map[string]any{
					"nested": "value",
				},
			},
			validate: func(result map[string]any) bool {
				nested, ok := result["key"].(map[string]any)
				if !ok {
					return false
				}
				return nested["nested"] == "value"
			},
		},
		{
			name: "empty map[any]any",
			input: map[string]any{
				"empty": map[any]any{},
			},
			validate: func(result map[string]any) bool {
				empty, ok := result["empty"].(map[string]any)
				if !ok {
					return false
				}
				return len(empty) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IntfaceKeysToStrings(tt.input)
			if !tt.validate(tt.input) {
				t.Errorf("validation failed for input: %v", tt.input)
			}
		})
	}
}
