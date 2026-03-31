package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadBytesReadsFileContents(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.txt")

	const want = "hello-from-file-provider"
	if err := os.WriteFile(path, []byte(want), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	p := Provider(path)
	b, err := p.ReadBytes()
	if err != nil {
		t.Fatalf("ReadBytes() returned error: %v", err)
	}

	if string(b) != want {
		t.Fatalf("unexpected file contents: got %q want %q", string(b), want)
	}
}

func TestReadBytesReturnsErrorForMissingFile(t *testing.T) {
	p := Provider(filepath.Join(t.TempDir(), "does-not-exist.txt"))

	b, err := p.ReadBytes()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if b != nil {
		t.Fatalf("expected nil bytes, got %v", b)
	}
}

func TestReadReturnsUnsupportedError(t *testing.T) {
	p := Provider("unused-path")

	mp, err := p.Read()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if mp != nil {
		t.Fatalf("expected nil map, got %v", mp)
	}

	const wantErr = "file provider does not support this method"
	if err.Error() != wantErr {
		t.Fatalf("unexpected error: got %q want %q", err.Error(), wantErr)
	}
}
