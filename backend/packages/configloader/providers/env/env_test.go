package env

import (
	"strings"
	"testing"
)

func TestReadBytesReturnsUnsupportedError(t *testing.T) {
	p := Provider(Opt{})

	b, err := p.ReadBytes()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if b != nil {
		t.Fatalf("expected nil bytes, got %v", b)
	}

	const wantErr = "env provider does not support this method"
	if err.Error() != wantErr {
		t.Fatalf("unexpected error: got %q want %q", err.Error(), wantErr)
	}
}

func TestReadWithoutPrefixIncludesSetEnv(t *testing.T) {
	key := "CONFIGLOADER_TEST_NO_PREFIX_KEY"
	val := "configloader-test-no-prefix-value"
	t.Setenv(key, val)

	p := Provider(Opt{})
	mp, err := p.Read()
	if err != nil {
		t.Fatalf("Read() returned error: %v", err)
	}

	got, ok := mp[key]
	if !ok {
		t.Fatalf("expected key %q in map", key)
	}
	if got != val {
		t.Fatalf("unexpected value for %q: got %v want %q", key, got, val)
	}
}

func TestReadWithPrefixFiltersKeys(t *testing.T) {
	prefix := "CONFIGLOADER_TEST_PREFIX_"
	includedKey := prefix + "INCLUDED"
	anotherIncludedKey := prefix + "ANOTHER"
	excludedKey := "CONFIGLOADER_TEST_OTHER_EXCLUDED"

	t.Setenv(includedKey, "value-1")
	t.Setenv(anotherIncludedKey, "value-2")
	t.Setenv(excludedKey, "should-not-be-included")

	p := Provider(Opt{Prefix: prefix})
	mp, err := p.Read()
	if err != nil {
		t.Fatalf("Read() returned error: %v", err)
	}

	if len(mp) == 0 {
		t.Fatal("expected non-empty map for matching prefix")
	}

	for k := range mp {
		if !strings.HasPrefix(k, prefix) {
			t.Fatalf("found non-prefixed key %q in result", k)
		}
	}

	if got := mp[includedKey]; got != "value-1" {
		t.Fatalf("unexpected value for %q: got %v want %q", includedKey, got, "value-1")
	}
	if got := mp[anotherIncludedKey]; got != "value-2" {
		t.Fatalf("unexpected value for %q: got %v want %q", anotherIncludedKey, got, "value-2")
	}
	if _, ok := mp[excludedKey]; ok {
		t.Fatalf("did not expect key %q in result", excludedKey)
	}
}
