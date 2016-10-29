package generation

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type TemplateLoader interface {
	Load(path string) (map[string]interface{}, error)
}

type FileHelper interface {
	Exists(path string) bool
}

type Generator struct {
	tmplLoader TemplateLoader
	fileHelper FileHelper
}

func (g Generator) checkPrerequisites(env string) error {
	if !viper.IsSet("scenarios.all") {
		return errors.New("'scenarios.all' not configured'")
	}

	templateFolderPath := getTemplateFolderPath()

	allScenarioServices := viper.GetStringSlice("scenarios.all")
	for _, templateName := range allScenarioServices {
		tmplPath := templateFolderPath + templateName + ".yml"
		if !g.fileHelper.Exists(tmplPath) {
			return errors.New("template '" + tmplPath + "' not found")
		}
	}

	return nil
}

func (g Generator) Generate(env string) error {
	templateFolderPath := getTemplateFolderPath()
	allScenarioServices := viper.GetStringSlice("scenarios.all")
	for _, templateName := range allScenarioServices {
		g.tmplLoader.Load(templateFolderPath + templateName + ".yml")
	}
	return g.checkPrerequisites(env)
}

func getTemplateFolderPath() string {
	templateFolderPath := viper.GetString("templateFolderPath")
	return convertToFolderPath(templateFolderPath)
}

func convertToFolderPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")
	return path + "/"
}
