package generation

import (
	"io/ioutil"
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
	return ioutil.WriteFile(path, []byte(data), 0644)
}
