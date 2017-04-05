package generation

import yaml "gopkg.in/yaml.v2"

type DefaultTemplateLoader struct {
	fileHelper FileHelper
}

// Load loads data from path into dictionary
func (tl DefaultTemplateLoader) Load(path string) (map[string]interface{}, error) {
	text, err := tl.fileHelper.Read(path)
	if err != nil {
		return nil, err
	}

	tmplYaml := make(map[string]interface{})
	err = yaml.Unmarshal(text, &tmplYaml)

	if err != nil {
		return nil, err
	}

	return tmplYaml, nil
}
