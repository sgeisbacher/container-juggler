package generation

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type DefaultTemplateLoader struct{}

func (tl DefaultTemplateLoader) Load(path string) (map[interface{}]interface{}, error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmplYaml := make(map[interface{}]interface{})
	err = yaml.Unmarshal(text, &tmplYaml)

	if err != nil {
		return nil, err
	}

	return tmplYaml, nil
}
