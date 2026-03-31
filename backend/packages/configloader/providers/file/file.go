package file

import (
	"errors"
	"os"
)

type File struct {
	path string
}

func Provider(path string) *File {
	return &File{
		path: path,
	}
}

// ReadBytes reads the contents of a file on disk and returns the bytes.
func (f *File) ReadBytes() ([]byte, error) {
	return os.ReadFile(f.path)
}

// Read is not supported by the file provider.
func (f *File) Read() (map[string]any, error) {
	return nil, errors.New("file provider does not support this method")
}
