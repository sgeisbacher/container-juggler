package generation

import "io/ioutil"

type DefaultFileHelper struct{}

func (fh DefaultFileHelper) Exists(path string) bool {
	return true
}

func (fh DefaultFileHelper) Write(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}
