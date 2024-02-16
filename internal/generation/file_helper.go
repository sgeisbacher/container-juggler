package generation

import (
	"os"
)

// DefaultFileHelper using real filesystem
type DefaultFileHelper struct{}

// Exists checks if file at path exists
func (fh DefaultFileHelper) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fh DefaultFileHelper) Write(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func (fh DefaultFileHelper) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
